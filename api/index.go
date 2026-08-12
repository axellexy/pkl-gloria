package handler

import (
    handler "jurnal-backend/api"
    "net/http"
    "strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
    if strings.HasPrefix(r.URL.Path, "/api/") {
        r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api")
    }
    handler.Handler(w, r)
}

func main() {}
