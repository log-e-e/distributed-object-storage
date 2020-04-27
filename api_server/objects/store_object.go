package objects

import (
    "distributed-object-storage/api_server/locate"
    "distributed-object-storage/src/utils"
    "fmt"
    "io"
    "net/http"
    "net/url"
)

func storeObject(reader io.Reader, hash string, size int64) (statusCode int, err error) {
    escapedHash := url.PathEscape(hash)
    // 若是对象的内容数据已存在，则不用重复上传，否则将对象数据保存到临时缓存中等待校验
    if locate.Exist(escapedHash) {
        return http.StatusOK, nil
    }

    // 对象数据是第一次存储时，执行以下代码

    // 保存对象数据到数据节点临时缓存
    // 实际上，通过执行putStream，我们已经将对象数据的信息（大小和哈希值）缓存到数据节点的缓存区中的临时文件中
    // putStream返回的io.Writer对象stream会在后续将数据写入到临时文件中
    stream, err := putStream(escapedHash, size)
    if err != nil {
        return http.StatusInternalServerError, err
    }

    // 使用TeeReader实现在读取reader的数据的同时，将读取的数据写入stream中
    // 也就是说，我们通过io.TeeReader()返回的io.Reader对象读取数据，同时将读取的数据通过io.Writer对象stream将写入临时文件中
    // 这样，我们就实现了在apiServer层对数据进行边读边校验的操作
    r := io.TeeReader(reader, stream)
    actualHash := utils.CalculateHash(r)
    if actualHash != hash {
        stream.Commit(false)
        err = fmt.Errorf("apiServer Error: object hash value is not match, actualHash=[%s], expectedHash=[%s]\n", actualHash, hash)
        return http.StatusBadRequest, err
    }
    stream.Commit(true)

    return http.StatusOK, nil
}
