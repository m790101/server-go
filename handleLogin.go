package main

import (
	"encoding/json"
	"net/http"
	"server/internal/auth"
	"time"
)

func (cfg *apiConfig) Login(w http.ResponseWriter, r *http.Request) {

	type parametersLogin struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	type LoginRes struct {
		Id    int    `json:"id"`
		Email string `json:"email"`
		Token string `json:"token"`
	}

	w.Header().Set("Content-Type", "application/json")
	params := parametersLogin{}
	decoder := json.NewDecoder(r.Body)
	errDecode := decoder.Decode(&params)

	if errDecode != nil {
		responseWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
	}

	users, errUsers := cfg.Db.GetUsers()
	if errUsers != nil {
		responseWithError(w, http.StatusInternalServerError, "Error getting users")
		return
	}

	defaultExpiration := 60 * 60 * 24
	if params.ExpiresInSeconds == 0 {
		params.ExpiresInSeconds = defaultExpiration
	} else if params.ExpiresInSeconds > defaultExpiration {
		params.ExpiresInSeconds = defaultExpiration
	}

	for _, user := range users {
		err := auth.CheckPasswordHash(params.Password, user.Password)

		if err == nil {
			token, err := cfg.MakeJWT(user.Id, cfg.Secret, time.Duration(params.ExpiresInSeconds)*time.Second)

			if err != nil {
				responseWithError(w, http.StatusInternalServerError, "Couldn't create JWT")
				return
			}

			validUser := LoginRes{
				Id:    user.Id,
				Email: user.Email,
				Token: token,
			}
			respondWithJSON(w, http.StatusOK, validUser)
		} else {
			respondWithJSON(w, http.StatusUnauthorized, "")
		}
	}

}
