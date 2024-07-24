package main

import (
	"encoding/json"
	"net/http"
	"server/internal/auth"
	"server/internal/database"
)

func (cfg *apiConfig) Login(w http.ResponseWriter, r *http.Request) {

	type parametersLogin struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type LoginRes struct {
		Id    int    `json:"id"`
		Email string `json:"email"`
		Token string `json:"token"`
	}

	w.Header().Set("Content-Type", "application/json")
	db, _ := database.NewDB("./database.json")
	params := parametersLogin{}
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
				Token: cfg.Secret,
			}
			respondWithJSON(w, http.StatusOK, validUser)
		} else {
			respondWithJSON(w, http.StatusUnauthorized, "")
		}
	}

}
