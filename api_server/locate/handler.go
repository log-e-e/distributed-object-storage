package locate

import (
    "distributed-object-storage/src/utils"
    "encoding/json"
    "net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    location := Locate(utils.GetObjectName(r.URL.EscapedPath()))
    if len(location) == 0 {
        w.WriteHeader(http.StatusNotFound)
        return
    }

    locationJson, _ := json.Marshal(location)
    w.Write(locationJson)
}
