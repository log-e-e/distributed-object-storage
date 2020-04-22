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
        log.Println("API-Server HTTP Error: missing object hash in request header")
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
    // 存储对象数据
    name := url.PathEscape(metadata.Hash)
    stream, err := getStream(name)
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusNotFound)
        return
    }
    io.Copy(w, stream)
}

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
