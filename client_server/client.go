package main

import (
    "distributed-object-storage/client_server/object_operations"
    "flag"
    "log"
)

const (
    PUT = "put"
    GET = "get"
    DELETE = "delete"
)

var (
    apiServer = flag.String("apiServer", "", "")
    operationType = flag.String("operationType", "", "")
    objectName = flag.String("objectName", "", "")
    desFilepath = flag.String("file", "", "")
    ioType = flag.String("ioType", "", "")
    ioValue = flag.String("ioValue", "", "")
    version = flag.Int("version", 0, "")
)


func main() {
    flag.Parse()
    if *operationType == "" || *objectName == "" || *apiServer == "" {
        log.Printf("Error: operationType[%s] or objectName[%s] or apiServer[%s] is invalid",
            *operationType, *objectName, *apiServer)
        return
    }

    if *operationType == PUT {
        // 获取PUT操作的对象内容形式及其值
        if (*ioType == "" || *ioValue == "") || (*ioType != "-path" && *ioType != "-content") {
            log.Printf("Error: ioType[%s] or ioValue[%s] is invalid", *ioType, *ioValue)
            return
        }
        log.Printf("INFO: apiServer[%s] method[%s], object[%s], ioType[%s], ioValue[%s]\n", *apiServer, PUT, *objectName, *ioType, *ioValue)
        object_operations.PutObject(*apiServer, *objectName, *ioType, *ioValue)

        return
    }

    if *operationType == GET {
        log.Printf("INFO: apiServer[%s] method[%s], object[%s]-version[%d] saveTo[%s]\n", *apiServer, GET, *objectName, *version, *desFilepath)
        object_operations.GetObject(*apiServer, *objectName, *version, *desFilepath)
        return
    }

    if *operationType == DELETE {
        //deleteObject(*objectName)
        log.Printf("INFO: apiServer[%s] method[%s], object[%s]\n", *apiServer, DELETE, *objectName)
        return
    }

    // 如果不是以上任意一种情况，则不进行任何处理
    log.Printf("Error: unsupported operation type [%s]\n", *operationType)
}
