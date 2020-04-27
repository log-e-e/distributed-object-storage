package rs

import (
    "github.com/klauspost/reedsolomon"
    "io"
)

type rsDecoder struct {
    // 为什么不仅需要readers，还需要writers?因为在读取数据的同时需要进行可能的数据修复
    readers []io.Reader  // 可正常读且数据完好取对象分片的数据节点的文件读对象
    writers []io.Writer  // 不可正常读取或数据缺失的对象分片的数据节点的文件写对象
    rsEnc reedsolomon.Encoder  // reedsolomon中对对象进行分片、编码、解码及数据恢复都需要依靠该对象进行
    size int64  // 对象数据的大小，也就是数据分片中的实际数据量
    cache []byte  // 用于缓存读取的数据
    cacheSize int  // 用于计算缓存了多少数据
    total int64  // 用于读取数据时进行计数
}

func NewDecoder(readers []io.Reader, writers []io.Writer, size int64) *rsDecoder {
    enc, _ := reedsolomon.New(DATA_SHARDS, PARITY_SHARDS)
    return &rsDecoder{
        readers:   readers,
        writers:   writers,
        rsEnc:     enc,
        size:      size,
        cache:     make([]byte, 0, BLOCK_PER_SHARD),
        cacheSize: 0,
        total:     0,
    }
}

// 【读取并解码数据】从cache中读取数据
func (rsDec *rsDecoder) Read(p []byte) (n int, err error) {
    // 当缓存中没有数据时，会通过调用getData()获取数据
    if rsDec.cacheSize == 0 {
        err := rsDec.getData()
        if err != nil {
            return 0, err
        }
    }
    dataLength := len(p)
    if rsDec.cacheSize < dataLength {
        dataLength = rsDec.cacheSize
    }
    rsDec.cacheSize -= dataLength
    copy(p, rsDec.cache[: dataLength])
    rsDec.cache = rsDec.cache[dataLength:]
    return dataLength, nil
}

// 【解码与修复】解码对象数据分片，并将数据分片的数据存放到cache中，如果对象分片存在缺失则进行修复
func (rsDec *rsDecoder) getData() error {
    // 如果当前rsDecoder读取的数据总量total已达到对象数据大小，则直接返回已读完
    if rsDec.total == rsDec.size {
        return io.EOF
    }

    // 读取rsDecoder中的readers序列的文件读对象，将其数据写入到对应的字节序列中
    // 若是某一文件读对象为空则说明该编号的对象分片缺失，需要进行修复，并将其放入repairIds中
    shards := make([][]byte, ALL_SHARDS)
    repairIds := make([]int, 0)
    for i := 0; i < len(shards); i++ {
        if rsDec.readers[i] == nil {
            repairIds = append(repairIds, i)
        } else {
            shards[i] = make([]byte, BLOCK_PER_SHARD)
            readCount, err := io.ReadFull(rsDec.readers[i], shards[i])
            if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
                shards[i] = nil
            } else if readCount != BLOCK_PER_SHARD {
                shards[i] = shards[i][: readCount]
            }
        }
    }
    // 如果存在需要修复的分片，则进行分片修复
    if len(repairIds) > 0 {
        err := rsDec.rsEnc.Reconstruct(shards)
        if err != nil {
            return err
        }
        // 将恢复的分片写入对应的服务节点
        for i := 0; i < len(repairIds); i++ {
            id := repairIds[i]
            rsDec.writers[id].Write(shards[id])
        }
    }
    // 解码数据分片，还原数据
    for i := 0; i < DATA_SHARDS; i++ {
        shardSize := int64(len(shards[i]))
        // 如果处理到最后一块数据分片时，存在数据填充，则只取实际数据
        if rsDec.total + shardSize > rsDec.size {
            shardSize -= rsDec.total + shardSize - rsDec.size
        }
        // 将数据分片的数据存入缓存中，同时计算缓存数据总量
        rsDec.cache = append(rsDec.cache, shards[i][:shardSize]...)
        rsDec.cacheSize += int(shardSize)
        rsDec.total += shardSize
    }
    return nil
}


