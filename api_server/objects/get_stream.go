package objects

import (
    "distributed-object-storage/api_server/heartbeat"
    "distributed-object-storage/api_server/locate"
    "distributed-object-storage/src/rs"
    "fmt"
    "log"
)

func GetStream(objectName string, size int64) (*rs.RSGetStream, error) {
    locateInfo := locate.Locate(objectName)
    if len(locateInfo) < rs.DATA_SHARDS {
        return nil, fmt.Errorf("Error: object %s locate failed, the data shards located is not enough: %v\n",
            objectName, locateInfo)
    }
    // 若是获取的对象分片不足ALL_SHARDS，说明需要进行修复
    dataServers := make([]string, 0)
    if len(locateInfo) < rs.ALL_SHARDS {
        log.Printf("INFO: some of shards need to repair\n")
        dataServers = heartbeat.ChooseServers(rs.ALL_SHARDS - len(locateInfo), locateInfo)
    }
    return rs.NewRSGetStream(locateInfo, dataServers, objectName, size)
}
