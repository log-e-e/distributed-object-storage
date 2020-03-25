package heartbeat

import (
    "log"
    "math/rand"
    "strings"
)

func ChooseRandomDataServer() string {
    dataServers := GetAliveDataServers()
    serverCount := len(dataServers)

    if serverCount == 0 {
        return ""
    }

    log.Println("Alive data servers:", strings.Join(dataServers, ", "))
    return dataServers[rand.Intn(serverCount)]
}
