package temp

import "net/http"

func Handler(w http.ResponseWriter, r *http.Request) {
    method := r.Method
    if method == http.MethodHead {
        head(w, r)
        return
    }
    if method == http.MethodPut {
        put(w, r)
        return
    }

    w.WriteHeader(http.StatusMethodNotAllowed)
}
