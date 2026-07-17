package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}



func (cfg *apiConfig) addUser(w http.ResponseWriter, r *http.Request) {
		type params struct {
			Email string `json:"email"`
		}
		defer r.Body.Close()
		decoder := json.NewDecoder(r.Body)
		var email params
		err := decoder.Decode(&email)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Something went wrong")
		}

		user, err := cfg.dbQueries.CreateUser(r.Context(), email.Email)
		payload := User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		}

		respondWithJSON(w, http.StatusCreated, payload)
	
}
