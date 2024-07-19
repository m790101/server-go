package main

import (
	"net/http"
	"server/internal/database"
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
