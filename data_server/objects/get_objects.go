package objects

import (
    "crypto/sha256"
    "distributed-object-storage/data_server/global"
    "distributed-object-storage/data_server/locate"
    "encoding/base64"
    "fmt"
    "io"
    "log"
    "net/url"
    "os"
    "path"
    "path/filepath"
    "strings"
)

func getFilePath(hash string) string {
    // 模糊搜索文件名以["对象哈希" + "." + "切片编号"]开头的文件，并且一个数据服务节点上至多只有一个文件相匹配
    files, _ := filepath.Glob(path.Join(global.StoragePath, "objects", fmt.Sprintf("%s.*", hash)))
    // 在一个数据节点中，一般最多只有一个某一对象的分片文件
    if len(files) != 1 {
        return ""
    }
    // 将分片文件对象复制到sha256进行哈希计算
    shardFileName := files[0]
    h := sha256.New()
    shardFile, _ := os.Open(shardFileName)
    io.Copy(h, shardFile)
    // 哈希校验
    shardFileHash := url.PathEscape(base64.StdEncoding.EncodeToString(h.Sum(nil)))
    expectedHash := strings.Split(shardFileName, ".")[2]
    if shardFileHash != expectedHash {
        log.Printf("Shard file content hash: %s, expected hash: %s\n", shardFileHash, expectedHash)
        // 删除已损坏的分片的定位信息
        locate.Delete(hash)
        // 删除文件
        os.Remove(shardFileName)
        return ""
    }
    return shardFileName
}
