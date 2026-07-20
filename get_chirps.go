package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {

	chrips, err := cfg.dbQueries.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	var listOfChirp []Chirp
	for _, c := range chrips {
		listOfChirp = append(listOfChirp, Chirp{ID: c.ID, CreatedAt: c.CreatedAt, UpdatedAt: c.UpdatedAt, Body: c.Body, UserID: c.UserID})
	}
	respondWithJSON(w, http.StatusOK, listOfChirp)

}

func (cfg *apiConfig) getChirp(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid chirp id")
		return
	}

	chirp, err := cfg.dbQueries.GetChirp(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	chp := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	respondWithJSON(w, http.StatusOK, chp)
}
