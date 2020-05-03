package temp

import (
    "distributed-object-storage/data_server/global"
    "distributed-object-storage/src/utils"
    "io"
    "log"
    "net/http"
    "os"
    "path"
)

func get(w http.ResponseWriter, r *http.Request) {
    uuid := utils.GetObjectName(r.URL.EscapedPath())
    file, err := os.Open(path.Join(global.StoragePath, "temp", uuid + ".dat"))
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusNotFound)
        return
    }
    defer file.Close()

    // 将文件数据复制到响应体中
    io.Copy(w, file)
}
