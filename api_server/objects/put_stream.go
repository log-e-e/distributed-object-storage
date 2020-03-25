package objects

import (
    "distributed-object-storage/api_server/heartbeat"
    "distributed-object-storage/src/object_stream"
    "fmt"
    "log"
)

func putStream(objectName string) (*object_stream.PutStream, error) {
    server := heartbeat.ChooseRandomDataServer()
    log.Println("Choose random data server:", server)

    if server == "" {
        return nil, fmt.Errorf("Error: no alive data server\n")
    }
    return object_stream.NewPutStream(server, objectName), nil
}
