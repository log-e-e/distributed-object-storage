package objects

import (
    "distributed-object-storage/utils"
    "github.com/joho/godotenv"
    "log"
)

const (
    storageRootEnvName = "STORAGE_ROOT"
    objectParentDirName = "objects"
    uriSep = "/"
)

func init() {
    err := godotenv.Load(utils.GetEnvFiles()...)
    if err != nil {
        log.Fatalln("godotenv Error: env files load failed")
    }
}
