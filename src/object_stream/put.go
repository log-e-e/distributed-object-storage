package object_stream

import (
    "fmt"
    "io"
    "net/http"
)

type PutStream struct {
    writer *io.PipeWriter
    errorChan chan error
}

func NewPutStream(server, objectName string) *PutStream {
    if server == "" || objectName == "" {
        return nil
    }

    r, w := io.Pipe()
    errorChan := make(chan error)

    go func() {
        request, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("http://%s/objects/%s", server, objectName), r)
        httpClient := http.Client{}
        response, err := httpClient.Do(request)
        if err == nil && response.StatusCode != http.StatusOK {
            err = fmt.Errorf("Error: [dataServer] statusCode: %d\n", response.StatusCode)
        }
        errorChan <- err
    }()

    return &PutStream{
        writer:    w,
        errorChan: errorChan,
    }
}

func (putStream *PutStream) Write(p []byte) (n int, err error) {
    return putStream.writer.Write(p)
}

func (putStream *PutStream) Close() error {
    putStream.writer.Close()
    return <- putStream.errorChan
}
