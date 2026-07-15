package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) countHits(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hits: %d", cfg.fileServerHits.Load())
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}
