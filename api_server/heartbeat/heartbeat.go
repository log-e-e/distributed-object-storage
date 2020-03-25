package heartbeat

import (
    "distributed-object-storage/src/err_utils"
    "distributed-object-storage/src/rabbitmq"
    "os"
    "strconv"
    "sync"
    "time"
)

var (
    dataServerMap = make(map[string]time.Time)
    mutex sync.Mutex
)

func ListenHeartbeat() {
    mq := rabbitmq.NewRabbitMQ(os.Getenv("RABBITMQ_SERVER"))
    mq.BindExchange("apiServers")
    channel := mq.Consume()
    defer mq.Close()

    go removeExpiredServerNode()
    for msg := range channel {
        server, err := strconv.Unquote(string(msg.Body))
        err_utils.PanicNonNilError(err)
        mutex.Lock()
        dataServerMap[server] = time.Now()
        mutex.Unlock()
    }
}

func removeExpiredServerNode() {
    for {
        time.Sleep(5 * time.Second)
        mutex.Lock()
        for server, heartbeatTime := range dataServerMap {
            if heartbeatTime.Add(10 * time.Second).Before(time.Now()) {
                delete(dataServerMap, server)
            }
        }
        mutex.Unlock()
    }
}

func GetAliveDataServers() []string {
    mutex.Lock()
    defer mutex.Unlock()

    dataServers := make([]string, 0)
    for server, _ := range dataServerMap {
        dataServers = append(dataServers, server)
    }

    return dataServers
}
