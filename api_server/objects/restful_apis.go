package objects

import (
    "distributed-object-storage/src/utils"
    "io"
    "log"
    "net/http"
)

func put(w http.ResponseWriter, r *http.Request) {
    objectName := utils.GetObjectName(r.URL.EscapedPath())
    statusCode, err := storeObject(r.Body, objectName)
    if err != nil {
        log.Println(err)
    }
    w.WriteHeader(statusCode)
}

func get(w http.ResponseWriter, r *http.Request) {
    objectName := utils.GetObjectName(r.URL.EscapedPath())
    stream, err := getStream(objectName)
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    io.Copy(w, stream)
}
