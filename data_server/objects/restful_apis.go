package objects

import (
    "distributed-object-storage/src/utils"
    "net/http"
)

func get(w http.ResponseWriter, r *http.Request) {
    objectName := utils.GetObjectName(r.URL.EscapedPath())
    filePath := getFilePath(objectName)
    if filePath == "" {
        w.WriteHeader(http.StatusNotFound)
        return
    }
    sendFile(w, filePath)
}
