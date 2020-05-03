package objects

import (
    "distributed-object-storage/api_server/heartbeat"
    "distributed-object-storage/src/rs"
    "fmt"
)

func putStream(hash string, size int64) (*rs.RSPutStream, error) {
    // 在putStream()中，我们需要用数据服务节点存储对象的数据片和校验片
    // 因此是初次存储对象数据片，我们需要rs.ALL_SHARDS个服务节点，且不存在有该对象的数据分片，故第二个参数为nil
    servers := heartbeat.ChooseServers(rs.ALL_SHARDS, nil)
    if len(servers) != rs.ALL_SHARDS {
        return nil, fmt.Errorf("apiServer Error: cannot find enugh dataServers\n")
    }

    return rs.NewRSPutStream(servers, hash, size)
}
