package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type returnVals struct {
		Err string `json:"error"`
	}
	respondWithJSON(w, code, returnVals{Err:msg})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(dat)
}

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
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	respondWithJSON(w, 200, returnVals{CleanedBody:profanityCheck(params.Body)})
}