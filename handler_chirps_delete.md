package main

import (
	"net/http"
	"strconv"

	"github.com/jkvyff/simple-server/internal/auth"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := strconv.Atoi(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}
	subject, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}

	userIDInt, err := strconv.Atoi(subject)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse user ID")
		return
	}

	chirp, err := cfg.DB.GetChirp(chirpID)
    if err != nil {
        respondWithError(w, http.StatusNotFound, "Chirp not found")
        return
    }

    if chirp.AuthorID != userIDInt { 
        respondWithError(w, http.StatusForbidden, "You are not authorized to delete this chirp")
        return
    }

    err = cfg.DB.DeleteChirp(chirpID)
    if err != nil { 
        respondWithError(w, http.StatusInternalServerError, "Couldn't delete chirp")
        return
    }

    respondWithJSON(w, http.StatusOK, struct{}{}) 
}
