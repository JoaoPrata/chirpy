package main

import (
	"net/http"
	"time"
	"encoding/json"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type userParameters struct {
    Email string `json:"email"`
}

type userErrorResponse struct {
    Error string `json:"error"`
}

func (cfg *apiConfig) handlerUsers(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := userParameters{}
	err := decoder.Decode(&params)
	if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        respBody := userErrorResponse{
            Error: "Couldn't decode parameters",
        }
        data, err := json.Marshal(respBody)
        if err != nil {
            return
        }
        w.Write(data)
        return
    }
	user, err := cfg.dbQueries.CreateUser(r.Context(), params.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		respBody := userErrorResponse{
            Error: "Couldn't create user",
        }
        data, err := json.Marshal(respBody)
        if err != nil {
            return
        }
        w.Write(data)
        return
	}
	userResp := User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	}
	data, err := json.Marshal(userResp)
	if err != nil {
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
}