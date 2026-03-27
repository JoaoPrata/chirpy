package main

import (
	"net/http"
	"fmt"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	err := cfg.dbQueries.DeleteUsers(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Failed to reset users: %s", err)))
		return
	}
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}

func (cfg *apiConfig) handlerDisabledReset(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusForbidden)
}