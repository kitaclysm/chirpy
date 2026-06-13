package main

import (
	"net/http"
	"encoding/json"
	"database/sql"
	"errors"
)

func (cfg *apiConfig) handlerUpgrade(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := upgradeRequest{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Error fetching request",err)
		return
	}
	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	err := cfg.queries.UpgradeUser(r.Context(), params.Data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "error finding user", err)
		} else {
			respondWithError(w, http.StatusInternalServerError, "error upgrading user", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}