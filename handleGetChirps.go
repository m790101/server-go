package main

import (
	"log"
	"net/http"
	"server/internal/database"
	"strconv"
)

func (cfg *apiConfig) handleGetChirps(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	db, _ := database.NewDB("./database.json")
	chirps, err := db.GetChirps()
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Error getting chirps")
		return
	}
	respondWithJSON(w, http.StatusOK, chirps)

}

func (cfg *apiConfig) handleGetOne(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("chirpID")
	pathNum, err := strconv.Atoi(path)
	if err != nil {
		log.Fatal(path)
	}

	w.Header().Set("Content-Type", "application/json")
	db, _ := database.NewDB("./database.json")
	chirps, err := db.GetChirps()

	if err != nil {

		responseWithError(w, http.StatusInternalServerError, "Error getting chirps")
		return
	}

	if len(chirps) <= pathNum {
		respondWithJSON(w, http.StatusNotFound, "")
		return
	}
	chirp := chirps[pathNum-1]

	respondWithJSON(w, http.StatusOK, chirp)

}
