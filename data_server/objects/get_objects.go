package objects

import (
    "distributed-object-storage/data_server/global"
    "distributed-object-storage/data_server/locate"
    "distributed-object-storage/src/utils"
    "log"
    "net/url"
    "os"
    "path"
)

func getFilePath(hash string) string {
    filePath := path.Join(global.StoragePath, "objects", hash)
    file, _ := os.Open(filePath)
    // 计算实际存储的文件的哈希值
    storedObjectHash := url.PathEscape(utils.CalculateHash(file))
    file.Close()
    // 校验：校验接口层中ES存储的哈希值与实际存储的内容的哈希值是否一致，若是发生了变化则不一致，并且删除该对象数据
    // 数据存放久了可能会发生数据降解等问题，因此有必要做一致性校验
    if storedObjectHash != hash {
        log.Printf("dataServer INFO: the object`s stored in node %s is broken, we just have removed it from dataServer node.",
            global.ListenAddr)
        locate.Delete(hash)
        os.Remove(filePath)
        return ""
    }
    return filePath
}
