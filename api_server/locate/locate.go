package locate

import (
    "distributed-object-storage/src/rabbitmq"
    "distributed-object-storage/src/rs"
    "distributed-object-storage/src/types"
    "encoding/json"
    "os"
    "time"
)

func Locate(objectName string) (locateInfo map[int]string) {
    mq := rabbitmq.NewRabbitMQ(os.Getenv("RABBITMQ_SERVER"))

    // 绑定dataServer交换器，通过该交换器向各个数据服务节点投递消息
    mq.Publish("dataServers", objectName)
    // 获取临时消息队列的管道，接收消息
    channel := mq.Consume()
    // Publish()后，设置超时关闭连接，以判断资源是否存在
    go func() {
        time.Sleep(1 * time.Second)
        mq.Close()
    }()
    // 准备接收消息：接收该对象的所有分片（数据分片、校验分片）的消息
    locateInfo = make(map[int]string)
    for i := 0; i < rs.ALL_SHARDS; i++ {
        msg := <- channel
        if len(msg.Body) == 0 {
            return
        }
        var info types.LocateMessage
        json.Unmarshal(msg.Body, &info)
        locateInfo[info.ID] = info.Addr
    }
    return
}

func Exist(objectName string) bool {
    // RS码规则：当获取的分片大于等于数据分片的值时，则可以进行数据恢复，那么我们就认为我们能定位到该数据
    return len(Locate(objectName)) >= rs.DATA_SHARDS
}
