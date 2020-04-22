package temp

import (
    "distributed-object-storage/data_server/global"
    "io/ioutil"
    "os"
    "path"
    "time"
)

// 每隔12小时清理一次临时文件
func CleanTemp() {
    time.Sleep(12 * time.Hour)
    tmpDir := path.Join(global.StoragePath, "temp")
    files, _ := ioutil.ReadDir(tmpDir)
    for i := 0; i < len(files); i++ {
        dif := int(files[i].ModTime().Sub(time.Now()).Minutes())
        if dif >= 30 {
            os.Remove(path.Join(global.StoragePath, "temp", files[i].Name()))
        }
    }
}
