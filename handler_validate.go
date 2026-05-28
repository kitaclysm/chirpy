package main

import (
	"encoding/json"
	"net/http"
	"strings"
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

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
        Body string `json:"body"`
    }
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	respondWithJSON(w, 200, returnVals{CleanedBody:profanityCheck(params.Body)})
}