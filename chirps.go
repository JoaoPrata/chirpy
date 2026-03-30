package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/JoaoPrata/chirpy/internal/database"
	"net/http"
	"time"
	"strings"
    "slices"
)

var profane = []string {"kerfuffle", "sharbert", "fornax"}

type chirpParameters struct {
    Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

type chirpErrorResponse struct {
    Error string `json:"error"`
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID 	  uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := chirpParameters{}
	err := decoder.Decode(&params)
	if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        respBody := chirpErrorResponse{
            Error: "Couldn't decode parameters",
        }
        data, err := json.Marshal(respBody)
        if err != nil {
            return
        }
        w.Write(data)
        return
    }
	cleanChirp := cleanChirp(params.Body)
	chirp, err := cfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: cleanChirp,
		UserID: params.UserID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		respBody := chirpErrorResponse{
			Error: "Couldn't create chirp",
		}
		data, err := json.Marshal(respBody)
		if err != nil {
			return
		}
		w.Write(data)
		return
	}
	w.WriteHeader(http.StatusCreated)
	respBody := Chirp{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	}
	data, err := json.Marshal(respBody)
	if err != nil {
		return
	}
	w.Write(data)
}

func cleanChirp(chirp string) string {
    words := strings.Split(chirp, " ")
    for idx, word := range words {
        if slices.Contains(profane, strings.ToLower(word)){
            words[idx] = "****"
        } 
    }
    return strings.Join(words, " ")
}