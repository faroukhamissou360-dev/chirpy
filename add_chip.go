package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/faroukhamissou-dev/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) addChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
	} else if len := len(params.Body); len > 140 {
		respondWithError(w, 400, "Chirp is too long")
	} else {
		words := strings.Split(params.Body, " ")
		for i, w := range words {
			switch strings.ToLower(w) {
			case "kerfuffle":
				words[i] = "****"
			case "sharbert":
				words[i] = "****"
			case "fornax":
				words[i] = "****"
			default:

			}
		}
		cleaned_body := strings.Join(words, " ")

		chirp, err := cfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{Body: cleaned_body, UserID: params.UserID})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to create chirp")
		}

		res := Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
		respondWithJSON(w, http.StatusCreated, res)
	}

}
