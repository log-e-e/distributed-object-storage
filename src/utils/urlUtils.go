package utils

import (
    "log"
    "net/http"
    "strconv"
    "strings"
)

const (
    PARAMENTER_LENGTH = 6
)

func GetObjectName(url string) string {
    url = strings.TrimSpace(url)
    components := strings.Split(url, "/")

    return components[len(components) - 1]
}

func GetHashFromHeader(h http.Header) string {
    digest := h.Get("digest")
    // 存放hash值的参数名设为SHA-256，因此若是hash值为空或者参数名对应不上，则直接返回空串
    if len(digest) < 9 || digest[:8] != "SHA-256=" {
        return ""
    }
    return digest[8:]
}

func GetSizeFromHeader(h http.Header) int64 {
    size, _ := strconv.ParseInt(h.Get("content-length"), 0, 64)
    return size
}

func GetOffsetFromHeader(h http.Header) (offset, end int64) {
    byteRange := h.Get("range")
    log.Printf("GetOffsetFromHeader(): range content[%s]\n", byteRange)
    if len(byteRange) < PARAMENTER_LENGTH {
        return 0, 0
    }
    if byteRange[:PARAMENTER_LENGTH] != "bytes=" {
        return 0, 0
    }
    bytesPositions := strings.Split(byteRange[PARAMENTER_LENGTH:], "-")
    log.Printf("%v\n", bytesPositions)
    offset, _ = strconv.ParseInt(bytesPositions[0], 0, 64)
    end, _ = strconv.ParseInt(bytesPositions[1], 0, 64)
    log.Printf("offset[%d]-end[%d]\n", offset, end)

    return offset, end
}
