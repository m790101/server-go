package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func (cfg *apiConfig) handleDeleteOne(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("chirpID")
	pathNum, err := strconv.Atoi(path)
	if err != nil {
		log.Fatal(path)
	}

	jwtToken, err := GetBearerToken(r.Header)

	if err != nil {
		responseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	userIDString, _ := ValidateJWT(jwtToken, cfg.Secret)
	fmt.Sprintln(userIDString)
	userId, _ := strconv.Atoi(userIDString)

	chirp, err := cfg.Db.GetChirp(userId)

	if err != nil {
		respondWithJSON(w, http.StatusForbidden, chirp)
		return
	}

	cfg.Db.DeleteChirp(pathNum)

	w.Header().Set("Content-Type", "application/json")

	respondWithJSON(w, http.StatusNoContent, "")

}
