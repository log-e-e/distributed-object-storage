package utils

import (
    "testing"
)

func TestGetEnvFiles(t *testing.T) {
    for _, filepath := range GetEnvFiles() {
        println(filepath)
    }
}
