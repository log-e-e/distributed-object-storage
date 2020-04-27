package rs

import (
    "distributed-object-storage/src/object_stream"
    "fmt"
    "io"
)

//RSPutStream: 组合rsEncoder，在此基础上实现将对象分片数据写入到数据节点中
type RSPutStream struct {
    *rsEncoder
}

// NewRSPutStream： 封装TempPutStream，根据提供的数据节点IP:PORT将对象的分片数据信息（分片数据大小、分片数据哈希值）保存在节点的缓存中，等待上传分片数据时与实际的分片数据进行哈希校验
// 写入流的文件格式为：["对象哈希" + "." + "分片编号"]
func NewRSPutStream(dataServers []string, objectHash string, objectSize int64) (rsPutStream *RSPutStream, err error) {
    // 创建一个对象分片的写入流，因此存储其分片的数据节点数目应等于ALL_SHARDS
    if len(dataServers) != ALL_SHARDS {
        return nil, fmt.Errorf("Error: dataServer number is not enough\n")
    }

    // 向上取整，计算出每个分片的大小
    perShardSize := (objectSize + DATA_SHARDS - 1) / DATA_SHARDS
    // 创建分片信息（分片数据大小、分片数据哈希值）写入流，每个写入流的名字格式为["对象哈希" + "." + "分片编号"]
    writers := make([]io.Writer, ALL_SHARDS)
    for i := 0; i < len(writers); i++ {
        writers[i], err = object_stream.NewTempPutStream(dataServers[i], fmt.Sprintf("%s.%d", objectHash, i), perShardSize)
        if err != nil {
            return nil, err
        }
    }
    // 创建封装了reedsolomon的Encoder对象，用于对对象数据进行编码与写入
    rsEnc := NewRSEncoder(writers)
    rsPutStream = &RSPutStream{rsEnc}
    return rsPutStream, nil
}

// Commit()方法会将缓存在temp下的相关分片数据文件转正或删除，
// temp下的数据文件的名字格式由["对象哈希" + "." + "分片ID"]，变为["对象哈希" + "." + "分片ID" + "." + "分片哈希"]
func (rsPutStream *RSPutStream) Commit(positive bool) {
    // 将缓存中的数据刷新写入到对应的temp中，rsPutStream的cache中的数据可能是第一次写入，也可能是剩下未写入的数据
    // 总之，无论如何，在提交之前都需要将数据写入到对应节点的temp中
    // 这是最后一次刷新数据到rsEncoder中对应的各个io.Writer中，接下来就是执行TempPutStream的Commit方法将流中的数据写入到文件中
    rsPutStream.Flush()
    // 根据positive的值决定temp中的数据是转正还是删除
    for i := 0; i < len(rsPutStream.writers); i++ {
        rsPutStream.writers[i].(*object_stream.TempPutStream).Commit(positive)
    }
}
