package main

import (
	"net/http"
	"time"

	"github.com/kitaclysm/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	refreshString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error getting token", err)
		return
	}

	userID, err := cfg.queries.GetUserFromRefreshToken(r.Context(), refreshString)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error invalid token", err)
		return
	}

	jwtToken, err := auth.MakeJWT(userID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating token", err)
		return
	}

	type response struct {
		Token	string	`json:"token"`
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: jwtToken,
	})
}