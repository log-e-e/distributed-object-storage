package object_operations

import (
    "distributed-object-storage/src/es"
    "distributed-object-storage/src/utils"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "strings"
    "time"
)

const (
    OBJECT_CONTENT = "-content"
    OBJECT_PATH = "-path"
    // PUT操作可接受的对象数据的最大值，超过该值则需使用POST操作
    PUT_MAX_SIZE = 1024 * 1024 * 10
    GET_MAX_SIZE = 1024 * 1024 * 10
)

// PUT操作：上传小数据的对象
func PutObject(apiServer, objectName, ioType, ioValue string) {
    if ioType != OBJECT_CONTENT && ioType != OBJECT_PATH {
        log.Printf("Error: unsupported ioType[%s]\n", ioType)
        return
    }
    // 获取要上传的对象数据的大小
    var dataInFile bool
    var bodySize int64
    if ioType == OBJECT_CONTENT {
        bodySize = int64(len(ioValue))
    } else if ioType == OBJECT_PATH {
        file, err := os.Open(ioValue)
        if err != nil {
            log.Println(err)
            return
        }
        info, _ := file.Stat()
        bodySize = info.Size()
        file.Close()
        dataInFile = true
    }
    // 若是对象数据大小超过PUT可以接受的数据量大小，则不执行请求，并提示使用POST请求上传对象数据大小
    if bodySize > PUT_MAX_SIZE {
        log.Printf("INFO: object data size [%d]KB is too large, the supported max data size of PUT operation is %dKB\n" +
            "Please use POST to upload the object`s data\n", bodySize, PUT_MAX_SIZE)
        return
    }

    // 上传对象数据
    content := ""
    if dataInFile {
        data, err := ioutil.ReadFile(ioValue)
        if err != nil {
            log.Println(err)
            return
        }
        content = string(data)
    } else {
        content = ioValue
    }
    url := fmt.Sprintf("http://%s/objects/%s", apiServer, objectName)
    headerMap := map[string]string{"Digest": fmt.Sprintf("SHA-256=%s", utils.CalculateHash(strings.NewReader(content)))}
    response := doHttpRequest(http.MethodPut, url, content, headerMap)
    if response == nil {
        return
    }
    if response.StatusCode != http.StatusOK {
        log.Printf("clientServer Error: put object failed, STATUSCODE[%d]\n", response.StatusCode)
        return
    }

    echoToTerminal(strings.NewReader(fmt.Sprintf("INFO: object [%s] has been uploaded, object-length[%d]\n", objectName, bodySize)))
}

func GetObject(apiServer, objectName string, version int, destinationPath string) {
    // 获取对象元数据，通过元数据中的对象数据的大小决定是否采用断点下载方式获取对象内容
    metadata, err := es.GetMetadata(objectName, version)
    if err != nil {
        log.Printf("Error: failed to get object[%s] metadata, reason: %s\n", objectName, err)
        return
    }
    if metadata.Hash == "" {
        log.Printf("ES INFO: object [%s] not found", objectName)
        return
    }

    file, err := os.Create(destinationPath)
    if err != nil {
        log.Printf("Error: filepath[%s] is invalid, reason: %s\n", destinationPath, err)
        return
    }

    url := fmt.Sprintf("http://%s/objects/%s?version=%d", apiServer, objectName, version)

    // 若是对象数据的大小在规定阀值内，则不采用断点下载方式获取对象数据
    if metadata.Size <= GET_MAX_SIZE {
        response := doHttpRequest(http.MethodGet, url, "", nil)
        if response == nil {
            return
        }
        if response.StatusCode != http.StatusOK {
            log.Printf("HTTP Error: STATUSCODE[%d]\n", response.StatusCode)
            return
        }
        data, err := ioutil.ReadAll(response.Body)
        if err != nil {
            log.Printf("Error: failed to read content from response body, reason: %s\n", err)
            return
        }
        _, err = file.Write(data)
        if err != nil {
            log.Printf("Error: failed write object[%s] data to file[%s]\n, reason: %s\n", objectName, destinationPath, err)
            return
        }
        echoToTerminal(strings.NewReader(fmt.Sprintf("INFO: object[%s], object-data-path[%s]", objectName, destinationPath)))
        printPartialContent(10, destinationPath)
        return
    }
    // 未完成
    // ------------------------------------------------------------------------------------------------------
    // 对象数据超过阀值，则采用多线程断点下载方式
    sizeList := getBufferSizeList(metadata.Size)
    responseChans := make([]chan *http.Response, len(sizeList))
    var start, end int
    // 多线程断点下载各个数据段
    log.Println("Start send request")
    for index, size := range sizeList {
        // 此处减一，是因为start从0开始，且end需要是个可达上界，start与end的数据长度为(end - start + 1)
        end = start + size - 1

        responseChans[index] = make(chan *http.Response, 1)
        responseChans[index] <- new(http.Response)

        log.Printf("goroutine %d get range [%d, %d]\n", index, start, end)
        headerMap := map[string]string{"range": fmt.Sprintf("bytes=%d-%d", start, end)}
        go runHttpRequest(&responseChans[index], http.MethodGet, url, "", headerMap)

        start += size
    }
    log.Println("Finish sending request")
    time.Sleep(3 * time.Second)
    log.Println("Start write data to file")
    // 将数据写入文件中
    for i := 0; i < len(responseChans); i++ {
        response := <- responseChans[i]
        if response == nil {
            log.Println("clientServer Error: http request error")
            return
        }
        data := readDataFromReader(response.Body)
        file.Write(data)
    }
    log.Println("Finish writing data to file")

    echoToTerminal(strings.NewReader(fmt.Sprintf("INFO: object[%s], object-data-path[%s]", objectName, destinationPath)))
    printPartialContent(10, destinationPath)
    // ------------------------------------------------------------------------------------------------------
}

func deleteObject(apiServer, objectName string) {
    url := fmt.Sprintf("http://%s/objects/%s", apiServer, objectName)

    response := doHttpRequest(http.MethodGet, url, "", nil)
    if response == nil {
        return
    }
    if response.StatusCode != http.StatusOK {
        log.Printf("HTTP Error: STATUSCODE[%d]\n", response.StatusCode)
        return
    }

    log.Printf("INFO: object [%s] has been deleted", objectName)
}
