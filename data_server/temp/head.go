package temp

import (
    "distributed-object-storage/data_server/global"
    "distributed-object-storage/src/utils"
    "fmt"
    "log"
    "net/http"
    "os"
    "path/filepath"
)

// 获取缓存中指定对象的已上传的数据分片的数据量
func head(w http.ResponseWriter, r *http.Request) {
    uuid := utils.GetObjectName(r.URL.EscapedPath())
    file, err := os.Open(filepath.Join(global.StoragePath, "temp", uuid + ".dat"))
    if err != nil {
        log.Println("dataServer Error:", err)
        w.WriteHeader(http.StatusNotFound)
        return
    }
    defer file.Close()

    // 获取文件信息
    info, err := file.Stat()
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusNotFound)
        return
    }
    w.Header().Set("content-length", fmt.Sprintf("%d", info.Size()))
}
