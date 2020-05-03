package temp

import (
    "distributed-object-storage/src/rs"
    "distributed-object-storage/src/utils"
    "fmt"
    "log"
    "net/http"
)

// 主要用于获取已上传的数据的大小
func head(w http.ResponseWriter, r *http.Request) {
    token := utils.GetObjectName(r.URL.EscapedPath())
    stream, err := rs.NewRSResumablePUtStreamFromToken(token)
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusForbidden)
        return
    }
    uploadedSize := stream.CurrentSize()
    if uploadedSize < 0 {
        w.WriteHeader(http.StatusNotFound)
        return
    }
    w.Header().Set("content-length", fmt.Sprintf("%d", uploadedSize))
}
