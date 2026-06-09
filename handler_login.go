package main

import (
	"net/http"
	"encoding/json"
	"time"

	"github.com/kitaclysm/chirpy/internal/auth"
	"github.com/kitaclysm/chirpy/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := Parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Error decoding JSON",err)
		return
	}

	user, err := cfg.queries.GetUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error finding user", err)
		return
	}

	accuratePass, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error checking password", err)
		return
	}
	if !accuratePass {
		respondWithError(w, http.StatusUnauthorized, "Error checking password", err)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating token", err)
		return
	}

	refreshToken, err := cfg.queries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:		auth.MakeRefreshToken(),
		UserID:		user.ID,
		ExpiresAt:	time.Now().UTC().Add(time.Duration(time.Hour * 24 * 60)),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating refresh token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, LoginResponse{
		ID:				user.ID,
		CreatedAt: 		user.CreatedAt,
		UpdatedAt: 		user.UpdatedAt,
		Email: 			user.Email,
		Token:			token,
		RefreshToken:	refreshToken.Token,
	})
}