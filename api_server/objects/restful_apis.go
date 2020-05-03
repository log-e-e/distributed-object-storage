package objects

import (
    "distributed-object-storage/api_server/heartbeat"
    "distributed-object-storage/api_server/locate"
    "distributed-object-storage/src/es"
    "distributed-object-storage/src/rs"
    "distributed-object-storage/src/utils"
    "fmt"
    "io"
    "log"
    "net/http"
    "net/url"
    "strconv"
)

// POST操作：创建封装了对象信息及对象的数据节点临时写入流对象的结构体RSResumablePutStream，将其序列化返回
func post(w http.ResponseWriter, r *http.Request) {
    // 从http请求中获取对象信息【对象名、数据总大小、总数据哈希】
    objectName := utils.GetObjectName(r.URL.EscapedPath())
    size, err := strconv.ParseInt(r.Header.Get("size"), 0, 64)
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    hash := utils.GetHashFromHeader(r.Header)
    if hash == "" {
        log.Printf("apiServer Error: missing object [%s] hash\n", objectName)
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    // 哈希校验对象是否已存在，若存在则无需重复上传
    if locate.Exist(url.PathEscape(hash)) {
        err = es.AddVersion(objectName, size, hash)
        if err != nil {
            log.Println(err)
            w.WriteHeader(http.StatusInternalServerError)
            return
        }
        w.WriteHeader(http.StatusOK)
        return
    }
    // 若不存在，则按照以下步骤进行处理：
    // 1. 选取指定数目的可用的数据节点用于存储对象分片
    // 2. 获取封装了存储对象信息和各个数据节点写入流的流对象，将该对象序列化，在header中设置为location参数的值的组成部分
    dataServers := heartbeat.ChooseServers(rs.ALL_SHARDS, nil)
    if len(dataServers) != rs.ALL_SHARDS {
        log.Println("apiServer Error: dataServer is not enough")
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    stream, err := rs.NewRSResumablePutStream(dataServers, objectName, size, hash)
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    // 为什么封装起来？因为客户端的断点续传需要用到该对象的信息【对象名、对象数据大小、对象哈希】
    w.Header().Set("location", stream.ToToken())
    // 将RSResumablePutStream对象序列化后封装到header中返回给客户端
    w.WriteHeader(http.StatusCreated)
}

func put(w http.ResponseWriter, r *http.Request) {
    // 获取hash值
    hashVal := utils.GetHashFromHeader(r.Header)
    if hashVal == "" {
        log.Println("apiServer HTTP Error: missing object hash in request header")
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    // 获取size值
    size := utils.GetSizeFromHeader(r.Header)

    // 存储对象数据
    statusCode, err := storeObject(r.Body, hashVal, size)
    if err != nil {
        log.Println(err)
        w.WriteHeader(statusCode)
    }
    if statusCode != http.StatusOK {
        w.WriteHeader(statusCode)
        return
    }

    // 添加元数据
    name := utils.GetObjectName(r.URL.EscapedPath())
    err = es.AddVersion(name, size, hashVal)
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusInternalServerError)
    }
}

func get(w http.ResponseWriter, r *http.Request) {
    objectName := utils.GetObjectName(r.URL.EscapedPath())
    versionID := r.URL.Query().Get("version")
    var version int
    var err error
    // 如果有version参数则查找指定版本的对象，否则查找最新版本对象
    if len(versionID) != 0 {
        version, err = strconv.Atoi(versionID)
        if err != nil {
            log.Println(err)
            w.WriteHeader(http.StatusBadRequest)
            return
        }
    }
    // 从ES获取对象元数据信息，进而通过元数据新信息的hash值向数据服务层请求对象内容
    metadata, err := es.GetMetadata(objectName, version)
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    if metadata.Hash == "" {
        log.Printf("ES INFO: object [%s] not found", objectName)
        w.WriteHeader(http.StatusNotFound)
        return
    }
    // 通过GetStream()获取对象分片
    stream, err := GetStream(url.PathEscape(metadata.Hash), metadata.Size)
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusNotFound)
        return
    }
    // 获取读取位置
    offset, end := utils.GetOffsetFromHeader(r.Header)
    log.Printf("apiServer INFO: in get(), get object data range [%d, %d]\n", offset, end)
    // 将stream的读指针移动到offset位置，一般的文件读指针最初指向第一个字节
    contentLength := metadata.Size
    if offset != 0 {
        contentLength = end - offset + 1
        stream.Seek(offset, io.SeekCurrent)
        w.Header().Set("content-range", fmt.Sprintf("bytes %d-%d/%d", offset, end, metadata.Size))
        w.WriteHeader(http.StatusPartialContent)
    }
    // 将数据写入响应体中
    written, _ := io.CopyN(w, stream, contentLength)
    log.Println("Wrote to response length:", written)
    // 关闭steam流，并将可能存在的修复数据写入到对应的数据节点中
    stream.Close()
}

// 删除对象【逻辑删除】
func del(w http.ResponseWriter, r *http.Request) {
    name := utils.GetObjectName(r.URL.EscapedPath())
    latestMetadata, err := es.SearchLatestVersion(name)
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    // 逻辑删除：将size和hash置空即可
    err = es.PutMetadata(name, latestMetadata.Version + 1, 0, "")
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusInternalServerError)
    }
}
