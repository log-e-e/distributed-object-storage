package rs

import (
    "distributed-object-storage/src/object_stream"
    "distributed-object-storage/src/utils"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
)

// resumableToken: 封装的是对象的相关信息【对象名、总数据大小、哈希、数据分片存储的数据服务节点以及每个节点上对应的临时对象的uuid】
type resumableToken struct {
    ObjectName string
    Size int64
    Hash string
    Servers []string
    UUIDS []string
}

// RSResumablePutStream: 该结构体封装对象的临时数据写入流以及对象的信息
// 之所以将二者组合起来，是因为在将数据写入临时对象时需要校验是否与对象的相关信息【数据大小等】一致，同时也可用于转正校验
type RSResumablePutStream struct {
    *RSPutStream
    *resumableToken
}

func NewRSResumablePutStream(dataServers []string, name string, size int64, hash string) (*RSResumablePutStream, error) {
    // 获取对象临时数据的各个写入流
    putStream, err := NewRSPutStream(dataServers, hash, size)
    if err != nil {
        log.Println(err)
        return nil, err
    }
    // 获取各个数据节点上的uuid，用于创建token对象
    uuids := make([]string, ALL_SHARDS)
    for i := 0; i < ALL_SHARDS; i++ {
        uuids[i] = putStream.writers[i].(*object_stream.TempPutStream).UUID
    }
    token := &resumableToken{
        ObjectName: name,
        Size:       size,
        Hash:       hash,
        Servers:    dataServers,
        UUIDS:      uuids,
    }

    return &RSResumablePutStream{
        RSPutStream:    putStream,
        resumableToken: token,
    }, nil
}

// ToToken: 将token对象的数据序列化，以方便数据传递
func (s *RSResumablePutStream) ToToken() string {
    b, _ := json.Marshal(s)
    return base64.StdEncoding.EncodeToString(b)
}

// CurrentSize: 通过向数据服务节点发送head请求，获取已上传数据的大小
func (s *RSResumablePutStream) CurrentSize() int64 {
    // 向dataServer发送head请求，获取上传的临时对象的大小
    response, err := http.Head(fmt.Sprintf("http://%s/temp/%s", s.Servers[0], s.UUIDS[0]))
    if err != nil {
        log.Println(err)
        return NOT_FOUND
    }
    if response.StatusCode != http.StatusOK {
        log.Println("Error: http request [method: HEAD] status code:", response.StatusCode)
        return NOT_FOUND
    }
    size := utils.GetSizeFromHeader(response.Header) * DATA_SHARDS
    if size > s.Size {
        size = s.Size
    }

    return size
}

// 反序列化token，创建RSResumablePutStream对象
func NewRSResumablePUtStreamFromToken(token string) (*RSResumablePutStream, error) {
    // base64解码还原成字节序列
    b, err := base64.StdEncoding.DecodeString(token)
    if err != nil {
        log.Println(err)
        return nil, err
    }
    // 反序列化token
    var r resumableToken
    err = json.Unmarshal(b, &r)
    if err != nil {
        log.Println(err)
        return nil, err
    }

    // 创建RSResumablePutStream对象
    writers := make([]io.Writer, ALL_SHARDS)
    for i := 0; i < ALL_SHARDS; i++ {
        writers[i] = &object_stream.TempPutStream{
            Server: r.Servers[i],
            UUID:   r.UUIDS[i],
        }
    }
    enc := NewRSEncoder(writers)
    return &RSResumablePutStream{
        RSPutStream:    &RSPutStream{enc},
        resumableToken: &r,
    }, nil
}
