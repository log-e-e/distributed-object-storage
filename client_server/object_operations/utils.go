package object_operations

import (
    "distributed-object-storage/src/rs"
    "io"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "os/exec"
    "strconv"
    "strings"
)

const (
    // 存储各个节点IP:PORT信息的文件名
    IPADDR_NAME = ".ipAddrs"
    // 数据节点个数
    DATA_SERVER_NUM = 6
    API_SERVER_NUM = 2
    ALL_SERVER_NUM = DATA_SERVER_NUM + API_SERVER_NUM
)

func getBufferSizeList(totalSize int64) []int {
    blockCount, rest := int(totalSize / rs.BLOCK_SIZE), int(totalSize % rs.BLOCK_SIZE)
    buffers := make([]int, 0, blockCount + 1)
    for i := 1; i <= blockCount; i++ {
        buffers = append(buffers, rs.BLOCK_SIZE)
    }
    if rest != 0 {
        buffers = append(buffers, rest)
    }

    return buffers
}

func echoToTerminal(content io.Reader) {
    buf, err := ioutil.ReadAll(content)
    if err != nil {
        log.Print(err)
        return
    }
    cmd := exec.Command("echo", string(buf))
    cmd.Stdout = os.Stdout
    cmd.Run()
}

func printPartialContent(lineNumber int, path string) {
    echoToTerminal(strings.NewReader("-------------------------------------object partial data-------------------------------------"))
    cmd := exec.Command("head", strconv.Itoa(lineNumber), path)
    cmd.Stdout = os.Stdout
    cmd.Run()
    echoToTerminal(strings.NewReader("-------------------------------------object partial data-------------------------------------"))
}

func readDataFromReader(r io.Reader) []byte {
    data, err := ioutil.ReadAll(r)
    if err != nil {
        log.Printf("Error: failed to read data from io.Reader, reason: %s\n", err)
        return nil
    }
    return data
}

func doHttpRequest(method, url string, bodyContent string, headerMap map[string]string) *http.Response {
    var body io.Reader
    if bodyContent != "" {
        body = strings.NewReader(bodyContent)
    }

    request, err := http.NewRequest(method, url, body)
    if err != nil {
        log.Println(err)
        return nil
    }

    if headerMap != nil {
        for key, value := range headerMap {
            request.Header.Set(key, value)
        }
    }

    client := http.Client{}
    response, err := client.Do(request)
    if err != nil {
        log.Println(err)
        return nil
    }

    return response
}

func runHttpRequest(c *chan *http.Response, method, url string, bodyContent string, headerMap map[string]string) {
    response := doHttpRequest(method, url, bodyContent, headerMap)
    *c <- response
}
