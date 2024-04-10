package main

import (
	"encoding/json"
	"net/http"
)

type User struct {
	ID   int    `json:"id"`
	Email string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID   int    `json:"id"`
	Email string `json:"email"`
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	user, err := cfg.DB.CreateUser(params.Email, []byte(params.Password))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	respondWithJSON(w, http.StatusCreated, UserResponse{
		Email: user.Email,
		ID:   user.ID,
	})
}

func (cfg *apiConfig) handlerUsersLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	user, err := cfg.DB.LoginUser(params.Email, []byte(params.Password))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized , "Unauthorized to login")
		return
	}

	respondWithJSON(w, http.StatusOK, UserResponse{
		Email: user.Email,
		ID:   user.ID,
	})
}


// func validateUser(email string) (string, error) {
// 	return email, nil
// }
