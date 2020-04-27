package rs

const (
    // 数据分片
    DATA_SHARDS = 4
    // 校验分片
    PARITY_SHARDS = 2
    // 校验分片与数据分片的总和，该总和不应超出可用的服务节点
    ALL_SHARDS = DATA_SHARDS + PARITY_SHARDS
    // 每个分片一次性可写入的最大数据量
    BLOCK_PER_SHARD = 8000
    // 每个对象的所有分片一次性可写入的最大数据量，超过该数据量则要分批写入
    BLOCK_SIZE = BLOCK_PER_SHARD * DATA_SHARDS
)
