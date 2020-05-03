package temp

import (
    "distributed-object-storage/api_server/locate"
    "distributed-object-storage/src/es"
    "distributed-object-storage/src/rs"
    "distributed-object-storage/src/utils"
    "io"
    "log"
    "net/http"
    "net/url"
)

func put(w http.ResponseWriter, r *http.Request) {
    token := utils.GetObjectName(r.URL.EscapedPath())
    // 反序列化token，通过token中的对象信息建立PUT流对象
    stream, err := rs.NewRSResumablePUtStreamFromToken(token)
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusForbidden)
        return
    }
    // 校验请求中的offset与实际数据节点的大小是否一致
    offset, _ := utils.GetOffsetFromHeader(r.Header)
    uploadedSize := stream.CurrentSize()
    // 首先校验数据节点上是否存在对象数据
    if uploadedSize < 0 {
        w.WriteHeader(http.StatusNotFound)
        return
    }
    // 然后校验数据节点上的数据大小是否与offset相同
    if uploadedSize != offset {
        w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
        return
    }
    // 数据上传机制：每次上传数据至数据节点时至少是一个分片的大小，故为BLOCK_SIZE
    // 这意味着如果上传的数据大小不足BLOCK_SIZE且不是最后的数据，则会将其丢弃，避免因破坏上传规则带来的数据校验问题（获取已上传数据的计算方式是获取第一个数据节点的数据大小然后乘以DATA_SHARDS）
    buff := make([]byte, rs.BLOCK_SIZE)
    for {
        // 从请求体中读取上传的数据，每次读取BLOCK_SIZE大小
        readSize, err := io.ReadFull(r.Body, buff)
        // 若是读取数据的出错原因不是已读完或者在数据中间读取到了结束标志但实际仍有数据没读完
        // 则返回内部错误
        if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
            log.Println(err)
            w.WriteHeader(http.StatusInternalServerError)
            return
        }
        uploadedSize = uploadedSize + int64(readSize)
        // 若是在读取数据的过程中读取的数据总量超过了对象总数据量的大小，则说明上传的数据有误，直接清空缓存区的临时数据
        if uploadedSize > stream.Size {
            stream.Commit(false)
            log.Println("apiServer Error: the object data to be uploaded is mismatch with the uploaded data")
            w.WriteHeader(http.StatusForbidden)
            return
        }
        // 如果请求体中剩余的数据已不足以分配给每个数据节点一个分片大小的数据，并且当前累计的数据量并不等于对象总数据量，则丢弃该部分数据
        // 丢弃的原因很简单，若上传了该部分数据，会破坏临时对象数据的计算机制（获取一个数据节点的分片大小乘以数据分片的数量即为上传的数据量）
        if readSize != rs.BLOCK_SIZE && uploadedSize < stream.Size {
            return
        }
        // 将从请求体中读取的数据写入数据分片对应的数据节点中
        stream.Write(buff)
        // 若是已上传的数据正好等于对象总数据大小，则说明整个对象的数据已全部上传，进行最后的刷新处理
        if uploadedSize == stream.Size {
            // 将不足BLOCK_SIZE的最后一部分数据写入每个数据分片所在的数据节点中
            stream.Flush()
            // 接下来，获取每个数据分片所在的数据节点的数据进行哈希校验
            getStream, err := rs.NewRSResumableGetStream(stream.Servers, stream.UUIDS, stream.Size)
            if err != nil {
                log.Println(err)
                w.WriteHeader(http.StatusInternalServerError)
                return
            }
            actualHash := url.PathEscape(utils.CalculateHash(getStream))
            if actualHash != stream.Hash {
                log.Println("apiServer Error: the actual uploaded data`s hash is mismatch")
                w.WriteHeader(http.StatusForbidden)
                return
            }
            // 定位当前对象是否存在，若存在则不用进行数据转正，否则进行数据转正
            if locate.Exist(stream.Hash) {
                stream.Commit(false)
            } else {
                stream.Commit(true)
            }
            // 在元数据服务中添加对象元数据
            err = es.AddVersion(stream.ObjectName, stream.Size, stream.Hash)
            if err != nil {
                log.Println(err)
                w.WriteHeader(http.StatusInternalServerError)
            }
            return
        }
    }
}
