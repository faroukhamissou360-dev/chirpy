package main

import (
	"encoding/json"
	"net/http"

	"github.com/faroukhamissou-dev/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) webhook(w http.ResponseWriter, r *http.Request) {
	type Weebhook struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithJSON(w, http.StatusUnauthorized, nil)
		return
	}
	if apiKey != cfg.POLKA_KEY {
		respondWithJSON(w, http.StatusUnauthorized, nil)
		return
	}

	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	var webhook Weebhook
	err = decoder.Decode(&webhook)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if webhook.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusNoContent, nil)
		return
	}
	userId, err := uuid.Parse(webhook.Data.UserID)
	if err != nil {
		respondWithJSON(w, http.StatusNoContent, nil)
		return
	}
	_, err = cfg.dbQueries.UpgradeRed(r.Context(), userId)

	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	respondWithJSON(w, http.StatusNoContent, nil)
}
