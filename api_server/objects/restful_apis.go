package objects

import (
    "distributed-object-storage/src/es"
    "distributed-object-storage/src/utils"
    "io"
    "log"
    "net/http"
    "net/url"
    "strconv"
)

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
    // 在io.Copy()过程中，stream会执行Read()方法，而该Read()方法会对数据分片进行解码操作，如果该操作无误则可以正常读取数据，并将数据复制到w中
    _, err = io.Copy(w, stream)
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusNotFound)
        return
    }

    // 保证能够正常解码数据后，在将修复的数据保存到数据节点中
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
