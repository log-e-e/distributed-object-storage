package temp

import (
    "distributed-object-storage/data_server/global"
    "distributed-object-storage/src/utils"
    "log"
    "net/http"
    "os"
    "path"
)

func put(w http.ResponseWriter, r *http.Request) {
    uuid := utils.GetObjectName(r.URL.EscapedPath())
    // 从临时缓存区获取存储对象元数据的临时文件
    tempinfo, err := readFromFile(uuid)
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusNotFound)
        return
    }
    // 读取对象数据文件的内容
    infoFile := path.Join(global.StoragePath, "temp", uuid)
    dataFile := infoFile + ".dat"
    file, err := os.OpenFile(dataFile, os.O_WRONLY | os.O_APPEND, 0)
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    defer file.Close()
    // 校验临时缓存区的对象元数据与对象数据的大小是否匹配（一般是匹配的，并且哈希值校验在数据接口层就已完成）
    actualInfo, err := os.Stat(dataFile)
    os.Remove(infoFile)
    if actualInfo.Size() != tempinfo.Size {
        os.Remove(dataFile)
        log.Printf("Error: the actual uploaded file`s size [%d] is dismatched with expected size [%d]\n",
            actualInfo.Size(), tempinfo.Size)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    // 数据转正
    commitTempObject(dataFile, tempinfo)
}
