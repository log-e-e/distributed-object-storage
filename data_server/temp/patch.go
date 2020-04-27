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

// 接收来自apiServer上传的对象数据，并将对象数据根据apiServer请求中的uuid存放到指定的uuid.dat文件中
// 在将数据写入uuid.dat后校验该文件的大小与uuid文件中存放的tempInfo序列化信息中的size是否相等，若相等则完成该请求，否则说明数据被修改，直接删除相关临时文件
func patch(w http.ResponseWriter, r *http.Request) {
    uuid := utils.GetObjectName(r.URL.EscapedPath())
    // 反序列化uuid标志的临时存放了对象数据的校验信息文件
    tempinfo, err := readFromFile(uuid)
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusNotFound)
        return
    }
    // 将实际从http中获取的数据写入data文件中
    // infoFile存放了tempinfo序列化数据的文件，该文件在后续会删除，故拼出该文件的路径
    infoFile := path.Join(global.StoragePath, "temp", uuid)
    // dataFile是用于存放对象实际上传数据的临时文件
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

// 反序列化结构体tempInfo的数据
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
