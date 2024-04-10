package main

import (
	"net/http"
	"sort"
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

func (cfg *apiConfig) handlerUsersRetrieve(w http.ResponseWriter, r *http.Request) {
	dbUsers, err := cfg.DB.GetUsers()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve users")
		return
	}

	users := []User{}
	for _, dbUser := range dbUsers {
		users = append(users, User{
			Email: dbUser.Email,
			ID:   dbUser.ID,
		})
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].ID < users[j].ID
	})

	respondWithJSON(w, http.StatusOK, users)
}