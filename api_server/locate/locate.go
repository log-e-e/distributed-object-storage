package locate

import (
    "distributed-object-storage/src/rabbitmq"
    "log"
    "os"
    "strconv"
    "time"
)

func Locate(objectName string) string {
    mq := rabbitmq.NewRabbitMQ(os.Getenv("RABBITMQ_SERVER"))

    mq.Publish("dataServers", objectName)
    channel := mq.Consume()
    // Publish()后，设置超时关闭连接，以判断资源是否存在
    go func() {
        time.Sleep(1 * time.Second)
        mq.Close()
    }()
    // 准备接收消息
    msg := <- channel
    result, _ := strconv.Unquote(string(msg.Body))
    if result != "" {
        log.Printf("INFO: object [%s] at server '%s'\n", objectName, result)
    } else {
        log.Printf("INFO: object [%s] not found\n", objectName)
    }

    return result
}

func Exist(objectName string) bool {
    return Locate(objectName) != ""
}
