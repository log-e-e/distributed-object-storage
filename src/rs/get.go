package rs

import (
    "distributed-object-storage/src/object_stream"
    "fmt"
    "io"
)

type RSGetStream struct {
    *rsDecoder
}

func NewRSGetStream(locateInfo map[int]string, dataServers []string, hash string, size int64) (*RSGetStream, error) {
    if len(locateInfo) + len(dataServers) != ALL_SHARDS {
        return nil, fmt.Errorf("Error: dataServer number is not equal to %d\n", ALL_SHARDS)
    }
    // 创建对象分片的读对象
    readers := make([]io.Reader, ALL_SHARDS)
    for i := 0; i < ALL_SHARDS; i++ {
        server := locateInfo[i]
        // 如果locateInfo中的数据服务节点为空串，说明该编号对应的对象分片需要修复，因此给该编号分配一个可用的随机服务节点，存储修复的分片
        if server == "" {
            locateInfo[i] = dataServers[0]
            dataServers = dataServers[1:]
        // 否则，创建分片编号对应的读对象，读取对应服务节点中的分片数据
        } else {
            reader, err := object_stream.NewGetStream(server, fmt.Sprintf("%s.%d", hash, i))
            if err == nil {
                readers[i] = reader
            }
        }
    }
    // 为对应编号缺失的分片创建写对象，以便修复后写入数据服务节点
    writers := make([]io.Writer, ALL_SHARDS)
    perShard := (size + DATA_SHARDS - 1) / DATA_SHARDS
    var err error
    for i := 0; i < ALL_SHARDS; i++ {
        if readers[i] == nil {
            writers[i], err = object_stream.NewTempPutStream(locateInfo[i], fmt.Sprintf("%s.%d", hash, i), perShard)
            if err != nil {
                return nil, err
            }
        }
    }

    // 创建解码对象
    dec := NewDecoder(readers, writers, size)
    return &RSGetStream{dec}, nil
}

// 关闭RSGetStream时将修复的数据写入数据节点中
func (s *RSGetStream) Close() {
    for i := 0; i < len(s.writers); i++ {
        if s.writers[i] != nil {
            s.writers[i].(*object_stream.TempPutStream).Commit(true)
        }
    }
}
