package main

import (
	"net/http"
	"strconv"
)

func (cfg *apiConfig) handlerUsersRetrieveByID(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("userId")
	intPath, err := strconv.Atoi(path)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't find Id")
		return
	}
	dbUser, err := cfg.DB.GetUserByID(intPath)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User Id could not be found")
		return
	}

	respondWithJSON(w, http.StatusOK, dbUser)
}
