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
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
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
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}

	respondWithJSON(w, http.StatusCreated, payload)

}

func (cfg *apiConfig) login(w http.ResponseWriter, r *http.Request) {
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

		accesstoken, err := auth.MakeJWT(user.ID, cfg.SECRET_KEY)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Something went wrong")
			return
		}
		refreh_token := auth.MakeRefreshToken()
		ref_tok, err := cfg.dbQueries.AddRefreshToken(r.Context(), database.AddRefreshTokenParams{
			Token:     refreh_token,
			UserID:    user.ID,
			ExpiresAt: time.Now().Add(60 * 24 * time.Hour)})
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Something went wrong")
			return
		}
		payload := User{
			ID:           user.ID,
			CreatedAt:    user.CreatedAt,
			UpdatedAt:    user.UpdatedAt,
			Email:        user.Email,
			Token:        accesstoken,
			RefreshToken: ref_tok.Token,
			IsChirpyRed:  user.IsChirpyRed,
		}

		respondWithJSON(w, http.StatusOK, payload)
	}

}

func (cfg *apiConfig) refresh(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	refresh_token, err := cfg.dbQueries.GetRefreshToken(r.Context(), token)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Not Authorized")
		return
	}

	if refresh_token.ExpiresAt.Before(time.Now()) || refresh_token.RevokedAt.Valid {

		respondWithError(w, http.StatusUnauthorized, "Not Authorized")
		return
	}
	user, err := cfg.dbQueries.GetUserFromRefreshToken(r.Context(), refresh_token.UserID)
	accesstoken, err := auth.MakeJWT(user.ID, cfg.SECRET_KEY)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	type AccessToken struct {
		Token string `json:"token"`
	}
	payload := AccessToken{Token: accesstoken}

	respondWithJSON(w, http.StatusOK, payload)
	return

}

func (cfg *apiConfig) revoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	refresh_token, err := cfg.dbQueries.GetRefreshToken(r.Context(), token)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Not Authorized")
		return
	}
	err = cfg.dbQueries.RevokeRefreshToken(r.Context(), refresh_token.Token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Not Authorized")
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
	return

}

func (cfg *apiConfig) updateUser(w http.ResponseWriter, r *http.Request) {

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

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Something went wrong")
		return
	}
	id, err := auth.ValidateJWT(token, cfg.SECRET_KEY)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Not Authotized")
		return
	}
	hashed_pass, err := auth.HashPassword(parameters.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	new_user, err := cfg.dbQueries.UpdateUser(r.Context(), database.UpdateUserParams{Email: parameters.Email, HashedPassword: hashed_pass, ID: id})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}
	refreshToken, err := cfg.dbQueries.GetRTFromUserID(r.Context(), new_user.ID)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Not Authotized")
		return
	}
	payload := User{
		Email:        new_user.Email,
		CreatedAt:    new_user.CreatedAt,
		UpdatedAt:    new_user.UpdatedAt,
		ID:           new_user.ID,
		Token:        token,
		RefreshToken: refreshToken,
		IsChirpyRed:  new_user.IsChirpyRed,
	}
	respondWithJSON(w, http.StatusOK, payload)
}
