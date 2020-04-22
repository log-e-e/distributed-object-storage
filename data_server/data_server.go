package main

import (
    "distributed-object-storage/data_server/global"
    "distributed-object-storage/data_server/heartbeat"
    "distributed-object-storage/data_server/locate"
    "distributed-object-storage/data_server/objects"
    "distributed-object-storage/data_server/temp"
    "flag"
    "log"
    "net/http"
)

func main() {
    flag.StringVar(&global.ListenAddr, "listenAddr", "", "")
    flag.StringVar(&global.StoragePath, "storageRoot", "", "")
    flag.Parse()
    global.CheckSharedVars()

    go heartbeat.StartHeartbeat()
    go locate.ListenLocate()
    go temp.CleanTemp()
    http.HandleFunc("/objects/", objects.Handler)
    http.HandleFunc("/temp/", temp.Handler)
    log.Fatalln(http.ListenAndServe(global.ListenAddr, nil))
}
