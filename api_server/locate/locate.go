package locate

import (
    "distributed-object-storage/src/err_utils"
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
    result, err := strconv.Unquote(string(msg.Body))
    err_utils.PanicNonNilError(err)
    log.Printf("INFO: object at server '%s'\n", result)
    
    return result
}

func Exist(objectName string) bool {
    return Locate(objectName) != ""
}
