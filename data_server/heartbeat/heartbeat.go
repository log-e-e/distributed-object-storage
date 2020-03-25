package heartbeat

import (
    "distributed-object-storage/data_server/global"
    "distributed-object-storage/src/rabbitmq"
    "os"
    "time"
)

func StartHeartbeat() {
    mq := rabbitmq.NewRabbitMQ(os.Getenv("RABBITMQ_SERVER"))
    defer mq.Close()

    for {
        mq.Publish("apiServers", global.ListenAddr)
        time.Sleep(5 * time.Second)
    }
}
