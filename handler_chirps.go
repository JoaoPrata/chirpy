package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
	"sort"

	"github.com/JoaoPrata/chirpy/internal/database"
	"github.com/JoaoPrata/chirpy/internal/auth"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	Body      string    `json:"body"`
}

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get token", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.tokenSecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	cleaned, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleaned,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("author_id")
	sorting := r.URL.Query().Get("sort")
	var data []database.Chirp
	var err error
	if id != "" {
		userID, err := uuid.Parse(id)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't parse author id", err)
			return
		}
		data, err = cfg.db.GetUserChirps(r.Context(), userID)
	} else {
		data, err = cfg.db.GetChirps(r.Context())
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
		return
	}
	chirps := mapChirps(data)
	if sorting == "desc" {
		sort.Slice(
			chirps, 
			func(i, j int) bool { return chirps[i].CreatedAt.After(chirps[j].CreatedAt) })
	}
	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse id", err)
		return
	}
	data, err := cfg.db.GetChirp(r.Context(), id)
	if err != nil {
		respondWithJSON(w, http.StatusNotFound, "Couldn't get chirp")
		return
	}
	respondWithJSON(w, http.StatusOK, mapChirp(data))
}

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get token", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.tokenSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't parse id", err)
		return
	}
	chirp, err := cfg.db.GetChirp(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't find chirp", err)
		return
	}
	if chirp.UserID != userID {
		respondWithJSON(w, http.StatusForbidden, "Chirp dosen't belong to user")
		return
	}

	err = cfg.db.DeleteChirp(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete chirp", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func validateChirp(body string) (string, error) {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleaned := getCleanedBody(body, badWords)
	return cleaned, nil
}

func getCleanedBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}

func mapChirps(data []database.Chirp) (chirps []Chirp) {
	for _, chirp := range data {
		chirps = append(chirps, Chirp{
			ID: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			UserID: chirp.UserID,
			Body: chirp.Body,
		})
	}
	return chirps
}

func mapChirp(data database.Chirp) (chirps Chirp) {
	return Chirp{
			ID: data.ID,
			CreatedAt: data.CreatedAt,
			UpdatedAt: data.UpdatedAt,
			UserID: data.UserID,
			Body: data.Body,
		}
}
