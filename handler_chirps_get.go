package main

import (
	"net/http"
	"time"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

type returnChirp struct {
	ID			uuid.UUID	`json:"id"`
	CreatedAt	time.Time	`json:"created_at"`
	UpdatedAt	time.Time	`json:"updated_at"`
	Body		string		`json:"body"`
	UserID		uuid.UUID	`json:"user_id"`
}

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
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

func (cfg *apiConfig) handlerChirpGet(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	chirp, err := cfg.queries.RetrieveChirp(r.Context(), chirpID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		} else {
			respondWithError(w, http.StatusInternalServerError, "Error retrieving chirp", err)
		}
		return
	}
	respondWithJSON(w, http.StatusOK, returnChirp{
		ID:			chirp.ID,
		CreatedAt:	chirp.CreatedAt,
		UpdatedAt:	chirp.UpdatedAt,
		Body:		chirp.Body,
		UserID:		chirp.UserID,
	})
}