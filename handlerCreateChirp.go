package main

import (
	"encoding/json"
	"log"
	"net/http"
	"server/internal/database"
	"strconv"
	"strings"
)

type parameters struct {
	Body string `json:"body"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {

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
	db, dbErr := database.NewDB("./database.json")
	if dbErr != nil {
		log.Fatal(dbErr)
	}

	// get jwt token
	token, _ := GetBearerToken(r.Header)
	userIDString, _ := ValidateJWT(token, cfg.Secret)
	id, _ := strconv.Atoi(userIDString)

	chirp, err := db.CreateChirp(res, id)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Error creating chirp")
		return
	}
	respondWithJSON(w, http.StatusCreated, chirp)
}
