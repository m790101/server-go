package main

import (
	"net/http"
	"time"
)

func (cfg *apiConfig) handleRefresh(w http.ResponseWriter, r *http.Request) {

	type Response struct {
		Token string `json:"token"`
	}

	refreshToken, err := GetBearerToken(r.Header)
	if err != nil {
		respondWithJSON(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}
	user, err := cfg.Db.UserForRefreshToken(refreshToken)
	if err != nil {
		respondWithJSON(w, http.StatusUnauthorized, "Couldn't get user for refresh token")
		return
	}
	accessToken, err := cfg.MakeJWT(
		user.Id,
		cfg.Secret,
		time.Hour,
	)

	if err != nil {
		respondWithJSON(w, http.StatusUnauthorized, "Couldn't validate token")
		return
	}

	respondWithJSON(w, http.StatusOK, Response{
		Token: accessToken,
	})

}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := GetBearerToken(r.Header)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, "Couldn't find token")
		return
	}

	err = cfg.Db.RevokeRefreshToken(refreshToken)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, "Couldn't revoke session")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
