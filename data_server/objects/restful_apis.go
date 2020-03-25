package objects

import (
    "distributed-object-storage/data_server/global"
    "distributed-object-storage/src/utils"
    "io"
    "net/http"
    "os"
    "path/filepath"
    "strings"
)

func get(w http.ResponseWriter, r *http.Request) {
    objectName := strings.Split(r.RequestURI, "/")[2]
    file, err := os.Open(filepath.Join(global.StoragePath, "objects", objectName))
    defer file.Close()
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    io.Copy(w, file)
}

func put(w http.ResponseWriter, r *http.Request) {
    objectName := utils.GetObjectName(r.URL.EscapedPath())
    file, err := os.Create(filepath.Join(global.StoragePath, "objects", objectName))
    defer file.Close()
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    io.Copy(file, r.Body)
}
