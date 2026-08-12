package main

import (
    handler "jurnal-backend/api"
    "log"
    "net/http"
    "strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
    if strings.HasPrefix(r.URL.Path, "/api/") {
        r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api")
    }
    handler.Handler(w, r)
}

func main() {
    http.HandleFunc("/", Handler)
    log.Println("Starting local API server on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
