package main

import (
	"encoding/json"
	"net/http"
	"server/internal/auth"
	"server/internal/database"
)

type LoginRes struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
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

	w.Header().Set("Content-Type", "application/json")
	db, _ := database.NewDB("./database.json")
	params := database.ParametersLogin{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)

	hashPassword, err := auth.HashPassword(params.Password)

	params.Password = hashPassword

	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
	}
	users, err := db.CreateUser(params)

	if err != nil {

		responseWithError(w, http.StatusInternalServerError, "Error getting chirps")
		return
	}

	respondWithJSON(w, http.StatusCreated, users)

}

func (cfg *apiConfig) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db, _ := database.NewDB("./database.json")
	params := database.ParametersLogin{}
	decoder := json.NewDecoder(r.Body)
	errDecode := decoder.Decode(&params)

	if errDecode != nil {
		responseWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
	}

	users, errUsers := db.GetUsers()
	if errUsers != nil {

		responseWithError(w, http.StatusInternalServerError, "Error getting users")
		return
	}

	for _, user := range users {
		err := auth.CheckPasswordHash(params.Password, user.Password)

		if err == nil {
			validUser := LoginRes{
				Id:    user.Id,
				Email: user.Email,
			}
			respondWithJSON(w, http.StatusOK, validUser)
		} else {
			respondWithJSON(w, http.StatusUnauthorized, "")
		}
	}

}
