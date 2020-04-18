package objects

import (
    "distributed-object-storage/api_server/locate"
    "distributed-object-storage/src/object_stream"
    "fmt"
    "io"
)

func getStream(objectName string) (io.Reader, error) {
    server := locate.Locate(objectName)
    if server == "" {
        return nil, fmt.Errorf("ERROR: object '%s' not found", objectName)
    }
    return object_stream.NewGetStream(server, objectName)
}
