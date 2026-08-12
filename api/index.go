package main

import (
    "net/http"
    "jurnal-backend/handler"
)

func Handler(w http.ResponseWriter, r *http.Request) {
    handler.Handler(w, r)
}
