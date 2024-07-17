package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func (cfg *apiConfig) handlerValidate(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Body string `json:"body"`
	}
	type responseType struct {
		Cleaned_body string `json:"cleaned_body"`
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
	content := params.Body
	spliteV := strings.Split(content, " ")
	res := ""
	for i, word := range spliteV {
		wordCheck := strings.ToLower(word)
		if wordCheck == "kerfuffle" || wordCheck == "sharbert" || wordCheck == "fornax" {
			spliteV[i] = "****"
		}
	}
	res = strings.Join(spliteV, " ")
	respondWithJSON(w, http.StatusOK, responseType{Cleaned_body: res})
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
