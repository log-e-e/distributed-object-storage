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

type tempInfo struct {
    UUID string
    Name string
    Size int64
}

// 创建存储对象元数据的临时文件，同时返回临时文件的uuid
func post(w http.ResponseWriter, r *http.Request) {
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
    // 缓存对象的临时元数据
    err = temp.writeToFile()
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    // 创建用于存放对象数据的临时文件，这个与对象的临时元数据文件的作用不一样，前者用于标志临时对象所在的服务节点，后者用于存放对象的内容数据
    os.Create(path.Join(global.StoragePath, "temp", temp.UUID + ".dat"))
    w.Write(uuidBytes)
}

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
