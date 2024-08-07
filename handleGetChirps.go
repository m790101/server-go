package main

import (
	"log"
	"net/http"
	"server/internal/database"
	"sort"
	"strconv"
)

func (cfg *apiConfig) handleGetChirps(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	// db, _ := database.NewDB("./database.json")
	chirps, err := cfg.Db.GetChirps()
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Error getting chirps")
		return
	}

	sortingOption := r.URL.Query().Get("sort")

	// default sorting
	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})

	if sortingOption == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].ID > chirps[j].ID
		})
	}

	authorIDString := r.URL.Query().Get("author_id")
	if authorIDString != "" {
		authorID, err := strconv.Atoi(authorIDString)
		if err != nil {
			respondWithJSON(w, http.StatusBadRequest, "Invalid author ID")
			return
		}
		filteredChirps := []database.Chirp{}
		for index, chirp := range chirps {
			if chirp.AuthorId == authorID {
				filteredChirps = append(filteredChirps, chirps[index])
			}
		}
		// sort.Slice(filteredChirps, func(i, j int) bool {
		// 	return chirps[i].ID < chirps[j].ID
		// })
		respondWithJSON(w, http.StatusOK, filteredChirps)
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
