package utils

import "strings"

func GetObjectName(url string) string {
    url = strings.TrimSpace(url)
    components := strings.Split(url, "/")

    return components[len(components) - 1]
}
