package heartbeat

import (
    "distributed-object-storage/src/rs"
    "log"
    "math/rand"
    "strings"
)

func ChooseRandomDataServer() string {
    dataServers := GetAliveDataServers()
    serverCount := len(dataServers)

    if serverCount == 0 {
        return ""
    }

    log.Println("Alive data servers:", strings.Join(dataServers, ", "))
    return dataServers[rand.Intn(serverCount)]
}

// ChooseServers: 选取dataServersNum个数据服务节点用于存放分片数据，该函数有两种使用方式：
// 1. 第一次存储对象分片数据，则dataServersNum等于ALL_SHARDS，unbrokenShardServerMap为nil
// 2. 从可用的数据服务节点中排除对象分片数据正常的节点，获取可用于存储修复的分片数据的数据服务节点
func ChooseServers(dataServersNum int, unbrokenShardServerMap map[int]string) (dataServers []string) {
    // 所需的用于存储的分片的节点数与已存放正常分片数据的节点数之和应等于一个对象的分片数之和，否则应直接中断程序执行
    if dataServersNum + len(unbrokenShardServerMap) != rs.ALL_SHARDS {
        panic("apiServer Error: the sum of brokenShards number and unbrokenShards number is not equal to ALL_SHARDS\n")
    }
    // 用于存放分片数据的候选服务节点，该切片长度不小于dataServersNum时，才可以进行随机选择
    candidateServers := make([]string, 0, dataServersNum)
    // 将分片ID与所在服务节点交换key-value，以便找出哪些编号的分片数据需要修复
    reversedUnbrokenShardServerMap := make(map[string]int)
    for id, serverAddr := range unbrokenShardServerMap {
        reversedUnbrokenShardServerMap[serverAddr] = id
    }

    // 获取可以作为分片存储节点的服务节点，实际上就是获取存储了需要修复分片数据的服务节点
    aliveServers := GetAliveDataServers()
    for i := 0; i < len(aliveServers); i++ {
        if _, in := reversedUnbrokenShardServerMap[aliveServers[i]]; !in {
            candidateServers = append(candidateServers, aliveServers[i])
        }
    }
    // 若是候选服务节点数小于所需的数据服务节点数在直接返回空，这说明没有足够的服务节点满足对象分片的存储需求
    if len(candidateServers) < dataServersNum {
        return
    }

    // 打乱并随机选择所需的数目
    randomIds := rand.Perm(len(candidateServers))
    for i := 0; i < dataServersNum; i++ {
        dataServers = append(dataServers, candidateServers[randomIds[i]])
    }
    return
}
