package objects

import (
    "distributed-object-storage/api_server/heartbeat"
    "distributed-object-storage/src/object_stream"
    "fmt"
    "log"
)

func putStream(hash string, size int64) (*object_stream.TempPutStream, error) {
    server := heartbeat.ChooseRandomDataServer()
    log.Println("Choose random data server:", server)

    if server == "" {
        return nil, fmt.Errorf("Error: no alive data server\n")
    }
    return object_stream.NewTempPutStream(server, hash, size)
}
