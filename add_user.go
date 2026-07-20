package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/faroukhamissou-dev/chirpy/internal/auth"
	"github.com/faroukhamissou-dev/chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
}

func (cfg *apiConfig) addUser(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	var parameters params
	err := decoder.Decode(&parameters)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}
	hashedPassword, err := auth.HashPassword(parameters.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	user, err := cfg.dbQueries.CreateUser(r.Context(), database.CreateUserParams{Email: parameters.Email, HashedPassword: hashedPassword})
	payload := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	respondWithJSON(w, http.StatusCreated, payload)

}

func (cfg *apiConfig) login(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	var parameters params
	err := decoder.Decode(&parameters)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	user, err := cfg.dbQueries.GetUser(r.Context(), parameters.Email)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "Something went wrong")
		return
	}
	match, err := auth.CheckPasswordHash(parameters.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Something went wrong")
		return
	}
	if !match {
		respondWithError(w, http.StatusUnauthorized, "Email or password incorrect")
		return
	} else {
		exp := parameters.ExpiresInSeconds
		if exp == 0 {
			exp = 3600
		} else if exp > 3600 {
			exp = 3600
		}
		timer := time.Duration(exp) * time.Second
		token, err := auth.MakeJWT(user.ID, cfg.SECRET_KEY, timer)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Something went wrong")
			return
		}
		payload := User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
			Token:     token,
		}

		respondWithJSON(w, http.StatusOK, payload)
	}

}
