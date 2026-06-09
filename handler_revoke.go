package main

import (
	"net/http"

	"github.com/kitaclysm/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	refreshString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error getting token", err)
		return
	}

	err = cfg.queries.RevokeRefreshToken(r.Context(), refreshString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error updating token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}