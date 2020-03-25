package objects

import (
    "io"
    "net/http"
)

func storeObject(reader io.Reader, objectName string) (statusCode int, err error) {
    stream, err := putStream(objectName)
    if err != nil {
        return http.StatusInternalServerError, err
    }

    io.Copy(stream, reader)
    err = stream.Close()
    if err != nil {
        return http.StatusInternalServerError, err
    }

    return http.StatusOK, nil
}
