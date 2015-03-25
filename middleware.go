package main
import (
    "net/http"
    "log"
)

// handler middleware which ensures that the path is valid
func validateURL(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
    m := validPath.FindStringSubmatch(r.URL.Path)
    if m == nil {
        log.Println("middleware being a boss")
        http.NotFound(w, r)
        return
    }
    next(w, r)
}