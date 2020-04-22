package locate

import (
    "distributed-object-storage/data_server/global"
    "distributed-object-storage/src/err_utils"
    "distributed-object-storage/src/rabbitmq"
    "os"
    "path"
    "path/filepath"
    "strconv"
    "sync"
)

var (
    objectMap = make(map[string]bool)
    mutex sync.Mutex
)

func ObjectExists(hash string) bool {
    mutex.Lock()
    _, ok := objectMap[hash]
    mutex.Unlock()
    return ok
}

func AddNewObject(hash string) {
    mutex.Lock()
    defer mutex.Unlock()
    objectMap[hash] = true
}

func Delete(hash string) {
    mutex.Lock()
    defer mutex.Unlock()
    delete(objectMap, hash)
}

func ListenLocate() {
    mq := rabbitmq.NewRabbitMQ(os.Getenv("RABBITMQ_SERVER"))
    defer mq.Close()

    mq.BindExchange("dataServers")
    channel := mq.Consume()
    for msg := range channel {
        hash, err := strconv.Unquote(string(msg.Body))
        err_utils.PanicNonNilError(err)
        if ObjectExists(hash) {
            mq.Send(msg.ReplyTo, global.ListenAddr)
        }
    }
}

// 扫描节点上已有的对象文件，载入内存中
func ScanObjects() {
    files, _ := filepath.Glob(path.Join(global.StoragePath, "objects", "*"))
    for i := 0; i < len(files); i++ {
        hash := filepath.Base(files[i])
        objectMap[hash] = true
    }
}
