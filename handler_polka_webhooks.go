package main

import(
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/JoaoPrata/chirpy/internal/auth"
	"github.com/JoaoPrata/chirpy/internal/database"
)

const (
	EventUpgradeUser string = "user.upgraded"
)

func (cfg *apiConfig) handlerPolkaWebHooks(w http.ResponseWriter, r *http.Request) {
	type parameters struct{
		Event string `json:"event"`
		Data struct{
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get api key", err)
		return
	}
	if apiKey != cfg.polkaKey {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if params.Event != EventUpgradeUser {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	id, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse user id", err)
		return
	}

	_, err = cfg.db.UpdateUserChirpyRed(r.Context(), database.UpdateUserChirpyRedParams{
		ID: id,
		IsChirpyRed: true,
	})
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}