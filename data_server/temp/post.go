package temp

import (
    "distributed-object-storage/data_server/global"
    "distributed-object-storage/src/utils"
    "encoding/json"
    "log"
    "net/http"
    "os"
    "os/exec"
    "path"
    "strconv"
    "strings"
)

// tempInfo保存的是对象数据的信息（数据哈希、数据大小）和该临时信息在缓存区中的唯一标识
// 对象数据的信息用于后续上传实际的对象数据时，与实际的数据进行校验（大小和哈希值）
// UUID用于在后续的临时数据请求中在临时数据中找到该临时文件
type tempInfo struct {
    UUID string  // 标志对象数据的临时信息的唯一值UUID
    Name string  // 对象分片数据的哈希值 + 分片编号：shardHash.ID
    Size int64  // 对象数据的大小
}

// post用于处理存放对象数据的信息的请求，该请求会创建相关临时文件用于存放对象数据的信息，这些信息用于后续的数据上传校验
// 临时文件：
// 1. uuid，该文件存放了对象数据的临时信息tempInfo
// 2. uuid.dat，该文件用于存放上传的数据，但是目前在该请求中为空
func post(w http.ResponseWriter, r *http.Request) {
    // 生成UUID并从请求中获取对象哈希和数据大小
    // 注意：产生的uuid值末尾会携带一个换行符，因此必须去除换行符
    uuidBytes, _ := exec.Command("uuidgen").Output()
    uuid := string(uuidBytes)
    uuid = strings.ReplaceAll(uuid, "\n", "")
    name := utils.GetObjectName(r.URL.EscapedPath())
    size, err := strconv.ParseInt(r.Header.Get("size"), 0, 64)
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    temp := tempInfo{
        UUID: uuid,
        Name: name,
        Size: size,
    }
    // 缓存用于对象数据校验的校验信息
    err = temp.writeToFile()
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    // 创建用于存放对象数据的临时文件，但目前该文件为空，需要由后续的patch请求往该文件写数据
    os.Create(path.Join(global.StoragePath, "temp", temp.UUID + ".dat"))

    // 将存放对象数据的临时文件的唯一标识uuid响应给请求方，请求方将根据该uuid在后续的请求中处理该uuid标志的临时文件
    w.Write(uuidBytes)
}

// writeToFile: 将用于校验对象数据的临时信息存放到文件中
// 通过序列化结构体数据存放到文件中，可以反序列化后相对较为方便地访问其中的信息
func (t *tempInfo) writeToFile() error {
    file, err := os.Create(path.Join(global.StoragePath, "temp", t.UUID))
    if err != nil {
        return err
    }
    defer file.Close()

    bytesData, _ := json.Marshal(t)
    file.Write(bytesData)
    return nil
}
