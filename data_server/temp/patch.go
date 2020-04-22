package temp

import (
    "distributed-object-storage/data_server/global"
    "distributed-object-storage/src/utils"
    "encoding/json"
    "io"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "path"
)

// 缓存对象数据，初步校验数据的大小是否匹配
func patch(w http.ResponseWriter, r *http.Request) {
    uuid := utils.GetObjectName(r.URL.EscapedPath())
    // 获取uuid对应的临时文件存放的对象元数据信息，用于校验实际上传的数据的信息是否正确
    tempinfo, err := readFromFile(uuid)
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusNotFound)
        return
    }
    // 将实际从http中获取的数据写入data文件中
    // infoFile和dataFile是在post请求时创建的，infoFile存放了对象的元数据，dataFile用于存放对象内容数据
    infoFile := path.Join(global.StoragePath, "temp", uuid)
    dataFile := infoFile + ".dat"
    file, err := os.OpenFile(dataFile, os.O_WRONLY | os.O_APPEND, 0)
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    defer file.Close()
    // 写入存储对象数据的临时文件中
    io.Copy(file, r.Body)
    // 比较实际获取的对象数据文件的大小与期望的对象数据大小是否相同
    // 若不相同，则删除创建的临时文件：对象元数据文件、对象数据文件
    actualInfo, err := os.Stat(dataFile)
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    if actualInfo.Size() != tempinfo.Size {
        os.Remove(infoFile)
        os.Remove(dataFile)
        log.Printf("Error: the actual uploaded file`s size [%d] is dismatched with expected size [%d]\n",
            actualInfo.Size(), tempinfo.Size)
        w.WriteHeader(http.StatusInternalServerError)
    }
}

// 读取存放对象元数据的临时文件
func readFromFile(uuid string) (*tempInfo, error) {
    file, err := os.Open(path.Join(global.StoragePath, "temp", uuid))
    if err != nil {
        return nil, err
    }
    defer file.Close()

    tempInfoBytes, _ := ioutil.ReadAll(file)
    var fileInfo *tempInfo
    json.Unmarshal(tempInfoBytes, &fileInfo)
    return fileInfo, nil
}
