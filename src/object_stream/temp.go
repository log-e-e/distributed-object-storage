package object_stream

import (
    "bytes"
    "fmt"
    "io/ioutil"
    "net/http"
    "strings"
)

type TempPutStream struct {
    Server string
    UUID string
}

func NewTempPutStream(server, objectName string, size int64) (*TempPutStream, error) {
    // 创建临时对象
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

// 实现io.Writer的Write接口
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
