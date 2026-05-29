package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	type returnChirp struct {
		ID			uuid.UUID	`json:"id"`
		CreatedAt	time.Time	`json:"created_at"`
		UpdatedAt	time.Time	`json:"updated_at"`
		Body		string		`json:"body"`
		UserID		uuid.UUID	`json:"user_id"`
	}

	chirpsSlice := []returnChirp{}

	chirps, err := cfg.queries.RetrieveChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error retrieving chirps", err)
		return
	}

	for _, chirp := range chirps {
		chirpsSlice = append(chirpsSlice, returnChirp{
			ID:			chirp.ID,
			CreatedAt:	chirp.CreatedAt,
			UpdatedAt:	chirp.UpdatedAt,
			Body:		chirp.Body,
			UserID:		chirp.UserID,
		})
	}
	respondWithJSON(w, http.StatusOK, chirpsSlice)
}