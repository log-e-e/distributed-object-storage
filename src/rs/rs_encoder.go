package rs

import (
    "github.com/klauspost/reedsolomon"
    "io"
)

// 封装了reedsolomon的encoder对象
// rsEncoder: rsEncoder的写数据的机制是首先将要写的数据缓存在cache中，
// 当cache中的数据量为BLOCK_SIZE(rs包中的global_vars)时便将数据写入到指定的位置中，否则先缓存在cache中，等待最终写入
type rsEncoder struct {
    writers []io.Writer  // io.Writer接口，用于将对象分片数据写入到指定的存储位置
    rsEncode reedsolomon.Encoder  // reedsolomon用于编码的调用对象
    cache []byte  // 缓存待写入的数据，大小一般为BLOCK_SIZE，一次性可写入ALL_SHARDS个切片的数据
}

func NewRSEncoder(writers []io.Writer) *rsEncoder {
    rsEnc, _ := reedsolomon.New(DATA_SHARDS, PARITY_SHARDS)
    return &rsEncoder{
        writers: writers,
        rsEncode: rsEnc,
        cache: make([]byte, 0),
    }
}

// rsEncoder的规则：当数据p超过可支持的一次性写入的最大数据量时，就需要分批次写入
// 需要注意的是，如果p的量不足一次性写入的最大量，则会延缓到Commit时才刷新写入对应的temp中
func (rsEnc *rsEncoder) Write(p []byte) (n int, err error) {
    dataLength := len(p)
    start := 0
    for dataLength != 0 {
        end := BLOCK_SIZE - len(rsEnc.cache)
        if end > dataLength {
            end = dataLength
        }
        rsEnc.cache = append(rsEnc.cache, p[start: end]...)
        // 若是cache的数据量已达到一次性可写入的最大的数据量，则先将该部分数据写入数据服务节点中，然后再读取下一批
        if len(rsEnc.cache) == BLOCK_SIZE {
            rsEnc.Flush()
        }
        start += end
        dataLength -= end
    }

    return len(p), nil
}

// 通过rsEnc中的各个writer对象将cache中的数据写入到对应的数据节点中
// Flush()的操作会将数据写入各个节点对应的temp下的缓存文件，
// temp下的缓存文件没有转正之前是不会计算分片内容哈希值的，因此不会因为多次调用Flush()时导致内容的不断变化而造成分片哈希值校验不通过的问题
func (rsEnc *rsEncoder) Flush() {
    if len(rsEnc.cache) == 0 {
        return
    }
    // 分片、编码、写入对应的服务节点的文件中
    shards, _ := rsEnc.rsEncode.Split(rsEnc.cache)
    rsEnc.rsEncode.Encode(shards)
    for i := 0; i < len(shards); i++ {
        rsEnc.writers[i].Write(shards[i])
    }
    rsEnc.cache = []byte{}
}
