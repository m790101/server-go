package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerWebhooks(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	type parametersWebhook struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		} `json:"data"`
	}

	params := parametersWebhook{}
	decoder := json.NewDecoder(r.Body)
	errDecode := decoder.Decode(&params)

	if errDecode != nil {
		responseWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
	}

	if params.Event == "user.upgraded" {
		id := params.Data.UserID
		user, err := cfg.Db.GetUser(id)
		if err != nil {
			respondWithJSON(w, http.StatusNotFound, "")
		}

		user.IsChirpyRed = true

		updateUser, err := cfg.Db.UpdateUser(user.Id, user.Email, user.Password, user.IsChirpyRed)

		fmt.Sprintln(updateUser)
		if err != nil {
			respondWithJSON(w, http.StatusNotFound, "")
		}
	}

	respondWithJSON(w, http.StatusNoContent, "")
}
