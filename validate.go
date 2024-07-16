package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerValidate(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Body string `json:"body"`
	}

	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
	}
	if len(params.Body) > 140 {
		w.WriteHeader(400)
		responseWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"valid":true}`))
}

func responseWithError(w http.ResponseWriter, code int, msg string) {
	type paramError struct {
		error string
	}
	if code > 499 {

		log.Printf("Error decoding parameters: %s", msg)
	}
	w.WriteHeader(500)
	respondWithJSON(w, 500, paramError{error: msg})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}
