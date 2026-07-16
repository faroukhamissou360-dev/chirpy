package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func validateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
	}else if len := len(params.Body); len > 140 {
		respondWithError(w, 400, "Chirp is too long")
	} else {
		type jsonRes struct {
			Valid bool `json:"valid"`
		}
		res := jsonRes{
			Valid: true,
		}
		respondWithJSON(w, 200, res)
	}

}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type errorMsg struct {
		ErrorMsg string `json:"error"`
	}
	res := errorMsg{
		ErrorMsg: msg,
	}
	data, err := json.Marshal(res)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {

	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}
