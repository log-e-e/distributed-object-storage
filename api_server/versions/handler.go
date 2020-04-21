package versions

import (
    "distributed-object-storage/src/es"
    "distributed-object-storage/src/utils"
    "encoding/json"
    "log"
    "net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    // from是ES搜索的起始页的数据序号，表示从头开始不跳过任何一条数据，size表示每一页的数据规模
    from, size := 0, 1000
    name := utils.GetObjectName(r.URL.EscapedPath())
    for {
        metadatas, err := es.SearchAllVersions(name, from, size)
        if err != nil {
            log.Println(err)
            w.WriteHeader(http.StatusInternalServerError)
            return
        }
        for index, _ := range metadatas {
            data, _ := json.Marshal(metadatas[index])
            w.Write(data)
            w.Write([]byte("\n"))
        }
        // 若该页的数据不足1000条，说明已经没有了，否则可能还存在数据，继续获取
        if len(metadatas) != size {
            return
        }
        from += size
    }
}
