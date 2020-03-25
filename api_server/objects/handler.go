package objects

import "net/http"

func Handler(w http.ResponseWriter, r *http.Request) {
    m := r.Method

    if m == http.MethodPut {
        put(w, r)
    } else if m == http.MethodGet {
        get(w, r)
    } else {
        // 如果不是以上请求方法的任一种，则返回405
        w.WriteHeader(http.StatusMethodNotAllowed)
    }
}
