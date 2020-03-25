package object_stream

import (
    "fmt"
    "io"
    "net/http"
)

type GetStream struct {
    reader io.Reader
}

func NewGetStream(server, objectName string) (*GetStream, error) {
    if server == "" || objectName == "" {
        return nil, fmt.Errorf("Value Error: server='%s', objectName='%s'\n", server, objectName)
    }
    return newGetStream(fmt.Sprintf("http://%s/objects/%s", server, objectName))
}

func newGetStream(url string) (*GetStream, error) {
    response, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    if response.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("Error: [dataServer] statusCode: %d\n", response.StatusCode)
    }

    return &GetStream{reader: response.Body}, nil
}

func (getStream *GetStream) Read(p []byte) (n int, err error) {
    return getStream.reader.Read(p)
}
