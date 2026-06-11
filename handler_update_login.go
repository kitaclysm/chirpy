package main

import (
	"net/http"
	"encoding/json"

	"github.com/kitaclysm/chirpy/internal/auth"
	"github.com/kitaclysm/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUpdateLogin(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := Parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Error decoding JSON",err)
		return
	}

	hashedPass, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error hashing password", err)
		return
	}

	updatedUser, err := cfg.queries.UpdateUser(r.Context(), database.UpdateUserParams{
		Email:			params.Email,
		HashedPassword:	hashedPass,
		ID:				userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error updating user", err)
		return
	}

	respondWithJSON(w, http.StatusOK, ReturnUser{
		ID:			updatedUser.ID,
		Email:		updatedUser.Email,
		UpdatedAt:	updatedUser.UpdatedAt,
	})
}