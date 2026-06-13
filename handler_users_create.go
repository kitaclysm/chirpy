package main

import (
	"net/http"
	"encoding/json"

	"github.com/kitaclysm/chirpy/internal/database"
	"github.com/kitaclysm/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerUserCreate(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := Parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Error decoding JSON",err)
		return
	}

	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error hashing password", err)
		return
	}

	user, err := cfg.queries.CreateUser(r.Context(), database.CreateUserParams{
		Email: params.Email,
		HashedPassword:	hash,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating user", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, ReturnUser{
		ID:				user.ID,
		CreatedAt: 		user.CreatedAt,
		UpdatedAt: 		user.UpdatedAt,
		Email: 			user.Email,
		IsChirpyRed:	user.IsChirpyRed,
	})
}