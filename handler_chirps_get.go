package main

import (
	"net/http"
	"sort"
	"strconv"
)

func (cfg *apiConfig) handlerChirpsRetrieveByID(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("chirpId")
	intPath, err := strconv.Atoi(path)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't find Id")
		return
	}
	dbChirp, err := cfg.DB.GetChirpByID(intPath)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp Id could not be found")
		return
	}

	respondWithJSON(w, http.StatusOK, dbChirp)
}

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
		return
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			Body: dbChirp.Body,
			ID:   dbChirp.ID,
		})
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})

	respondWithJSON(w, http.StatusOK, chirps)
}