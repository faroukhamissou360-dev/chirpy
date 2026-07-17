package main

import "net/http"

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
