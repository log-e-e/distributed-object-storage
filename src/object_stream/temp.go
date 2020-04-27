package object_stream

import (
    "bytes"
    "fmt"
    "io/ioutil"
    "net/http"
    "strings"
)

// 对象临时信息
type TempPutStream struct {
    Server string  // 临时信息所在的服务节点
    UUID string  // 临时信息的唯一标识
}

// NewTempPutStream: 通过POST请求，将即将上传的数据的信息（数据大小和数据哈希值）存储在节点的temp下，用于与实际上传的数据进行校验（数据大小和数据哈希校验）
func NewTempPutStream(server, objectName string, size int64) (*TempPutStream, error) {
    // 通过POST请求，将对象的数据大小和哈希值存储在数据节点的temp下
    request, err := http.NewRequest(http.MethodPost, "http://" + server + "/temp/" + objectName, nil)
    if err != nil {
        return nil, err
    }
    request.Header.Set("size", fmt.Sprintf("%d", size))
    httpClient := http.Client{}
    // 执行请求后，正常情况会在在响应中返回临时对象的UUID
    response, err := httpClient.Do(request)
    if err != nil {
        return nil, err
    }
    uuidBytes, err := ioutil.ReadAll(response.Body)
    // 注意：该处读取的uuid的值的末尾会有换行符，因此必须去除，否则会引起http url语法错误
    uuidBytes = bytes.ReplaceAll(uuidBytes, []byte("\n"), []byte(""))
    if err != nil {
        return nil, err
    }

    return &TempPutStream{
        Server: server,
        UUID:   string(uuidBytes),
    }, nil
}

// 该io.Writer接口的实现，主要是将对象数据上传至存放了对象数据临时信息的数据节点的temp下的uuid.dat文件中
// 我们使用patch方法用于局部数据更新，dataServer中的temp包下的patch方法会处理该请求
func (t *TempPutStream) Write(p []byte) (n int, err error) {
    // 通过patch操作，将数据写入临时文件对象，并且校验数据大小，若不符合则删除为该对象创建的临时文件【该步仅校验文件大小】
    request, err := http.NewRequest(http.MethodPatch, "http://" + t.Server + "/temp/" + t.UUID, strings.NewReader(string(p)))
    if err != nil {
        return 0, err
    }
    httpClient := http.Client{}
    response, err := httpClient.Do(request)
    if err != nil {
        return 0, err
    }
    if response.StatusCode != http.StatusOK {
        return 0, fmt.Errorf("dataServer Error: STATUSCODE[%d]\n", response.StatusCode)
    }
    return len(p), nil
}

// 根据positive决定是删除临时对象数据还是转正保存到节点中
func (t *TempPutStream) Commit(positive bool) {
    method := http.MethodDelete
    if positive {
        method = http.MethodPut
    }
    request, _ := http.NewRequest(method, "http://" + t.Server + "/temp/" + t.UUID, nil)
    httpClient := http.Client{}
    httpClient.Do(request)
}

func NewTempGetStream(server, uuid string) (*GetStream, error) {
    return newGetStream("http://" + server + "/temp/" + uuid)
}
