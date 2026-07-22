package main

import (
	"net/http"
	"sort"

	"github.com/faroukhamissou-dev/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	sortQuery := r.URL.Query().Get("sort")
	authorQuery := r.URL.Query().Get("author_id")
	if authorQuery != "" {
		authorId, err := uuid.Parse(authorQuery)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		chrips, err := cfg.dbQueries.GetChirpsByAuthor(r.Context(), authorId)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		var listOfChirp []Chirp
		for _, c := range chrips {
			listOfChirp = append(listOfChirp, Chirp{ID: c.ID, CreatedAt: c.CreatedAt, UpdatedAt: c.UpdatedAt, Body: c.Body, UserID: c.UserID})
		}
		if sortQuery == "desc" {
			sort.Slice(listOfChirp, func(i, j int) bool { return listOfChirp[i].CreatedAt.After(listOfChirp[j].CreatedAt) })

		}
		respondWithJSON(w, http.StatusOK, listOfChirp)
		return
	}

	chrips, err := cfg.dbQueries.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var listOfChirp []Chirp
	for _, c := range chrips {
		listOfChirp = append(listOfChirp, Chirp{ID: c.ID, CreatedAt: c.CreatedAt, UpdatedAt: c.UpdatedAt, Body: c.Body, UserID: c.UserID})
	}
	if sortQuery == "desc" {
		sort.Slice(listOfChirp, func(i, j int) bool { return listOfChirp[i].CreatedAt.After(listOfChirp[j].CreatedAt) })

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
	return
}

func (cfg *apiConfig) deleteChirp(w http.ResponseWriter, r *http.Request) {

	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid chirp ID")
		return
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Something went wrong")
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.SECRET_KEY)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Not Authotized")
		return
	}
	chirp, err := cfg.dbQueries.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Oops Chirp Not Found!")
		return
	}
	if userID != chirp.UserID {
		respondWithError(w, http.StatusForbidden, "Not Authotized")
		return
	}
	err = cfg.dbQueries.DeleteChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Oops Something went wrong!")
		return
	}
	respondWithJSON(w, http.StatusNoContent, nil)

}
