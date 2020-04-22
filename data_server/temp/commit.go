package temp

import (
    "distributed-object-storage/data_server/global"
    "distributed-object-storage/data_server/locate"
    "os"
    "path"
)

// 将临时文件移动至节点内部
func commitTempObject(tempFilePath string, tempinfo *tempInfo) {
    os.Rename(tempFilePath, path.Join(global.StoragePath, "objects", tempinfo.Name))
    locate.AddNewObject(tempinfo.Name)
}
