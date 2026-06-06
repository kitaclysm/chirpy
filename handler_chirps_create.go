package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/kitaclysm/chirpy/internal/database"
	"github.com/kitaclysm/chirpy/internal/auth"
)

func profanityCheck(msg string) string {
	words := strings.Split(msg, " ")
	returnWords := []string{}
	for _, word := range words {
		if strings.ToLower(word) == "kerfuffle" || strings.ToLower(word) == "sharbert" || strings.ToLower(word) == "fornax" {
			word = "****"
		}
		returnWords = append(returnWords, word)
	}
	return strings.Join(returnWords, " ")
}

func validate(s string) (string, error) {
	if len(s) > 140 {
		return "", errors.New("Chirp exceeds 140 character limit")
	}
	return profanityCheck(s), nil
}

func (cfg *apiConfig) handlerChirpCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
        Body	string		`json:"body"`
    }
	type returnChirp struct {
		ID			uuid.UUID	`json:"id"`
		CreatedAt	time.Time	`json:"created_at"`
		UpdatedAt	time.Time	`json:"updated_at"`
		Body		string		`json:"body"`
		UserID	uuid.UUID	`json:"user_id"`
	}

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error getting token", err)
		return
	}
	
	userID, err := auth.ValidateJWT(tokenString, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error validating token", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	validated, err := validate(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error validating: ", err)
		return
	}

	chirp, err := cfg.queries.CreateChirp(r.Context(), database.CreateChirpParams{Body: validated, UserID: userID})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error creating chirp: ", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, returnChirp{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	})
}