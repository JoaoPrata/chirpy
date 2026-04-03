package main

import (
	"net/http"
	"time"
	"fmt"
	"encoding/json"
	"github.com/google/uuid"

	"github.com/JoaoPrata/chirpy/internal/auth"
	"github.com/JoaoPrata/chirpy/internal/database"
)

const (
	expiresInSecondsDefault int = 3600
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerUsers(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	hashedPass, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}
	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email: params.Email,
		HashedPassword: hashedPass,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	})
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
		ExpiresInSeconds int `josn:"expires_in_seconds"`
	}
	type response struct {
		User
		Token string `json:"token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode login parameters", err)
		return
	}

	user, err := cfg.db.GetUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	validPass, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't validate password", err)
		return
	}
	if !validPass {
		respondWithJSON(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	var expiresIn int = expiresInSecondsDefault
	if params.ExpiresInSeconds != 0 && params.ExpiresInSeconds < expiresInSecondsDefault {
		expiresIn = params.ExpiresInSeconds
	}
	tokenDuration, err := time.ParseDuration(fmt.Sprintf("%ds", expiresIn))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't generate token", err)
		return
	}
	token, err := auth.MakeJWT(user.ID, cfg.tokenSecret, tokenDuration)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't generate token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token: token,
	})
}