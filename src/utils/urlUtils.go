package utils

import (
    "net/http"
    "strconv"
    "strings"
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
