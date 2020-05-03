package objects

import (
    "net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
    m := r.Method
    // POST：上传大对象
    if m == http.MethodPost {
        post(w, r)
        return
    }
    // PUT：用于存储小对象
    if m == http.MethodPut {
        put(w, r)
        return
    }
    // GET：获取对象数据
    if m == http.MethodGet {
        get(w, r)
        return
    }
    // DELETE：删除对象元数据
    if m == http.MethodDelete {
        del(w, r)
        return
    }
    // 如果不是以上请求方法的任一种，则返回405
    w.WriteHeader(http.StatusMethodNotAllowed)
}
