package locate

import (
    "distributed-object-storage/data_server/global"
    "distributed-object-storage/src/err_utils"
    "distributed-object-storage/src/rabbitmq"
    "distributed-object-storage/src/types"
    "log"
    "os"
    "path"
    "path/filepath"
    "strconv"
    "strings"
    "sync"
)

var (
    // key为对象分片的组合哈希名[对象哈希]，value为该分片的编号
    objectMap = make(map[string]int)
    mutex sync.Mutex
)

const (
    // 用-1表示对象分片不在该服务节点上
    NOT_FOUND = -1
    // 对象分片名字的组成："分片所属的对象哈希" + "." + "当前分片的编号" + "." + "当前分片数据的哈希"
    SHARD_NAME_COMPONENT_NUM = 3
)

func init() {
    // 扫描磁盘，将分片数据信息载入内存
    ScanObjects()
}

func ObjectExists(hash string) bool {
    return objectMap[hash] != NOT_FOUND
}

func AddNewObject(hash string, id int) {
    mutex.Lock()
    defer mutex.Unlock()
    objectMap[hash] = id
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
        // 哈希值是对象的内容数据哈希值，在apiServer层的ES中保存的还是对象整体内容数据的哈希值
        hash, err := strconv.Unquote(string(msg.Body))
        err_utils.PanicNonNilError(err)
        if ObjectExists(hash) {
            mq.Send(msg.ReplyTo, types.LocateMessage{
                Addr: global.ListenAddr,
                ID: objectMap[hash],
            })
        }
    }
}

// 扫描数据服务节点上已有的对象文件，载入内存中
func ScanObjects() {
    files, _ := filepath.Glob(path.Join(global.StoragePath, "objects", "*"))
    for i := 0; i < len(files); i++ {
        // shardNameComponents: 分为三部分[对象哈希，分片编号，分片哈希]
        shardNameComponents := strings.Split(filepath.Base(files[i]), ".")
        if len(shardNameComponents) != SHARD_NAME_COMPONENT_NUM {
            log.Fatalf("Error: shard %v name is invalid, it should be 3 compoments [objectHash.ID.shardHash]\n", shardNameComponents)
        }
        objectHash := shardNameComponents[0]
        shardID, err := strconv.Atoi(shardNameComponents[1])
        err_utils.PanicNonNilError(err)
        objectMap[objectHash] = shardID
    }
}
