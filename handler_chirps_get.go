package main

import (
	"net/http"
	"database/sql"
	"errors"
	"sort"

	"github.com/google/uuid"
	"github.com/kitaclysm/chirpy/internal/database"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	chirpsSlice := []returnChirp{}
	chirps := []database.Chirp{}
	var err error
	var authID uuid.UUID

	stringID := r.URL.Query().Get("author_id")
	if stringID != "" {
		authID, err = uuid.Parse(stringID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "error parsing author id", err)
			return
		}
		chirps, err = cfg.queries.GetChirpsByAuthor(r.Context(), authID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "error retrieving chirps", err)
			return
		}
	} else {
		chirps, err = cfg.queries.RetrieveChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error retrieving chirps", err)
			return
		}
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

	sorty := r.URL.Query().Get("sort")
	if sorty == "desc" {
		sort.Slice(chirpsSlice, func(i, j int) bool {
			return chirpsSlice[i].CreatedAt.After(chirpsSlice[j].CreatedAt)
		})
	} else if sorty == "asc" || sorty == "" {
		sort.Slice(chirpsSlice, func(i, j int) bool {
			return chirpsSlice[i].CreatedAt.Before(chirpsSlice[j].CreatedAt)
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

	chirp, err := cfg.queries.GetChirp(r.Context(), chirpID)
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