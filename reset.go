package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) resetHits(w http.ResponseWriter, r *http.Request) {
	cfg.fileServerHits.Store(0)
	if cfg.PLATFORM == "dev" {
		err := cfg.dbQueries.DeleteAllUsers(r.Context())
		if err != nil {
			log.Fatalf(err.Error())
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hits reset to 0"))
	} else {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("You're not a dev"))
	}

}
