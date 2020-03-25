package main

import (
    "distributed-object-storage/data_server/global"
    "distributed-object-storage/data_server/heartbeat"
    "distributed-object-storage/data_server/locate"
    "distributed-object-storage/data_server/objects"
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
    http.HandleFunc("/objects/", objects.Handler)
    log.Fatalln(http.ListenAndServe(global.ListenAddr, nil))
}
