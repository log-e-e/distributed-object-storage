package temp

import (
    "distributed-object-storage/data_server/global"
    "distributed-object-storage/src/utils"
    "net/http"
    "os"
    "path"
)

// 若是对象数据校验未通过则删除缓存区的临时文件
func del(w http.ResponseWriter, r *http.Request) {
    uuid := utils.GetObjectName(r.URL.EscapedPath())
    infoFile := path.Join(global.StoragePath, "temp", uuid)
    dataFile := infoFile + ".dat"
    os.Remove(infoFile)
    os.Remove(dataFile)
}
