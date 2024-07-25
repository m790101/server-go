package main

import (
	"encoding/json"
	"net/http"
	"server/internal/auth"
	"server/internal/database"
	"strconv"
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

func (cfg *apiConfig) handleUpdateUser(w http.ResponseWriter, r *http.Request) {

	type parametersLogin struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	params := parametersLogin{}
	decoder := json.NewDecoder(r.Body)
	errDecode := decoder.Decode(&params)

	if errDecode != nil {
		responseWithError(w, http.StatusInternalServerError, "Error decode")
	}

	hashPassword, _ := auth.HashPassword(params.Password)

	w.Header().Set("Content-Type", "application/json")

	token, err := GetBearerToken(r.Header)
	if err != nil {
		respondWithJSON(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}

	id, err := ValidateJWT(token, cfg.Secret)

	if err != nil {
		respondWithJSON(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}

	idNum, _ := strconv.Atoi(id)
	user, _ := cfg.Db.GetUser(idNum)

	userModified, _ := cfg.Db.UpdateUser(user.Id, params.Email, hashPassword)

	respondWithJSON(w, http.StatusOK, userModified)

}
