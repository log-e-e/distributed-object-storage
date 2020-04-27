package temp

import (
    "distributed-object-storage/data_server/global"
    "distributed-object-storage/data_server/locate"
    "distributed-object-storage/src/utils"
    "net/url"
    "os"
    "path"
    "strconv"
    "strings"
)

// 将形式为["对象哈希" + "." + "分片编号"]的临时分片文件转为名字形式为["对象哈希" + "." + "分片编号" + "." + "分片哈希"]的正式文件
func commitTempObject(tempFilePath string, tempinfo *tempInfo) {
    shardFile, _ := os.Open(tempFilePath)
    shardHash := url.PathEscape(utils.CalculateHash(shardFile))
    shardFile.Close()
    // 转正重命名
    os.Rename(tempFilePath, path.Join(global.StoragePath, "objects", tempinfo.Name + "." + shardHash))
    // 将分片编号信息加入到当前数据节点的内存中
    // map: key为对象总数据的哈希值，value为数据分片的编号ID
    locate.AddNewObject(tempinfo.hash(), tempinfo.id())
}

func (t *tempInfo) hash() string {
    return strings.Split(t.Name, ".")[0]
}

func (t *tempInfo) id() int {
    id, _ := strconv.Atoi(strings.Split(t.Name, ".")[1])
    return id
}
