package utils

import (
    "io/ioutil"
    "log"
    "path/filepath"
    "strings"
)

var (
    configDirPath, _ = filepath.Abs("config")
    envSuffix = ".env"
)

func GetEnvFiles() []string {
    fileInfos, err := ioutil.ReadDir(configDirPath)
    if err != nil {
        log.Fatalf("Error: config dir path '%s' is invalid", configDirPath)
    }

    envFilePaths := make([]string, 0)
    for _, fi := range fileInfos {
        if !fi.IsDir() && strings.HasSuffix(fi.Name(), envSuffix) {
            envFilePath, err := filepath.Abs(filepath.Join(configDirPath, fi.Name()))
            if err == nil {
                envFilePaths = append(envFilePaths, envFilePath)
            }
        }
    }

    return envFilePaths
}
