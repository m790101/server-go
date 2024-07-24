package main

import (
	"encoding/json"
	"net/http"
	"server/internal/auth"
	"server/internal/database"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

func (cfg *apiConfig) handleGetUsers(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	db, _ := database.NewDB("./database.json")
	users, err := db.GetUsers()
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Error getting chirps")
		return
	}
	respondWithJSON(w, http.StatusOK, users)

}

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {

	type parametersLogin struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type response struct {
		User
	}

	w.Header().Set("Content-Type", "application/json")
	db, _ := database.NewDB("./database.json")
	params := parametersLogin{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)

	hashPassword, err := auth.HashPassword(params.Password)

	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
	}
	users, err := db.CreateUser(params.Email, hashPassword)

	if err != nil {

		responseWithError(w, http.StatusInternalServerError, "Error getting chirps")
		return
	}

	respondWithJSON(w, http.StatusCreated, users)

}
