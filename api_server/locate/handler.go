package locate

import (
    "distributed-object-storage/src/err_utils"
    "encoding/json"
    "net/http"
    "strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    location := Locate(strings.Split(r.URL.EscapedPath(), "/")[2])
    if location == "" {
        w.WriteHeader(http.StatusNotFound)
        return
    }

    locationJson, err := json.Marshal(location)
    err_utils.PanicNonNilError(err)
    w.Write(locationJson)
}
