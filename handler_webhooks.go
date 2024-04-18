package main

import (
	"encoding/json"
	"net/http"

	"github.com/jkvyff/simple-server/internal/auth"
)

func (cfg *apiConfig) handlerPolka(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		} `json:"data"`
	}

	apiKey, err := auth.GetApiKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find auth token")
		return
	}

	if apiKey != cfg.polkaAPIKey {
		respondWithError(w, http.StatusUnauthorized, "Not authorized for request")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	if params.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusOK, struct{}{})
		return
	}

	user, err := cfg.DB.GetUserByID(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get user")
		return
	}

	_, err = cfg.DB.UpdateChirpyRed(user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user")
		return
	}

	respondWithJSON(w, http.StatusOK, struct{}{})
}