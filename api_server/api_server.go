package main

import (
    "distributed-object-storage/api_server/global"
    "distributed-object-storage/api_server/heartbeat"
    "distributed-object-storage/api_server/locate"
    "distributed-object-storage/api_server/objects"
    "flag"
    "log"
    "net/http"
)

func main() {
    flag.StringVar(&global.ListenAddr, "listenAddr", "", "")
    flag.Parse()
    global.CheckSharedVars()

    go heartbeat.ListenHeartbeat()
    http.HandleFunc("/objects/", objects.Handler)
    http.HandleFunc("/locate/", locate.Handler)
    log.Fatalln(http.ListenAndServe(global.ListenAddr, nil))
}
