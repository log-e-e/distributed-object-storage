package objects

import (
    "io"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "strings"
)

func put(w http.ResponseWriter, r *http.Request) {
    objectName := strings.Split(r.RequestURI, uriSep)[2]
    file, err := os.Create(filepath.Join(os.Getenv(storageRootEnvName), objectParentDirName, objectName))
    if err != nil {
        log.Println("PUT FAILED:", err.Error())
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    defer file.Close()

    io.Copy(file, r.Body)
    log.Printf("PUT SUCCESS: object '%s'\n", file.Name())
}

func get(w http.ResponseWriter, r *http.Request) {
    objectName := strings.Split(r.RequestURI, uriSep)[2]
    file, err := os.Open(filepath.Join(os.Getenv(storageRootEnvName), objectParentDirName, objectName))
    if err != nil {
        log.Println("GET FAILED:", err.Error())
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    defer file.Close()

    io.Copy(w, file)
    log.Printf("GET SUCCESS: object '%s'\n", file.Name())
}
