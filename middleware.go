package main
import (
    "net/http"
    "log"
)

// handler middleware which ensures that the path is valid
func validateURL(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
    m := validPath.FindStringSubmatch(r.URL.Path)
    log.Println("running middleware")
    if m == nil {
        log.Println("middleware being a boss")
        http.NotFound(w, r)
        return
    }
    next(w, r)
}

func myMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
    // do some stuff before
    log.Println("I am generic middleware")
    next(rw, r)
    // do some stuff after
}