package main

import (
    "distributed-object-storage/objects"
    "log"
    "net/http"
    "os"
)

const (
    objectsPattern = "/objects/"
    listenAddress = "LISTEN_ADDRESS"
)

var (

)

func init() {

}

func main() {
    http.HandleFunc(objectsPattern, objects.Handler)
    log.Println(http.ListenAndServe(os.Getenv(listenAddress), nil))
}
