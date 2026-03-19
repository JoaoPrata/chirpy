package main

import (
    "net/http"
    "encoding/json"
    "strings"
    "slices"
)

var profane = []string {"kerfuffle", "sharbert", "fornax"}

type parameters struct {
    Body string `json:"body"`
}
type returnVals struct {
    CleanedBody string `json:"cleaned_body"`
}

type errorResponse struct {
    Error string `json:"error"`
}

func (cfg *apiConfig) handlerValidate(w http.ResponseWriter, r *http.Request) {
    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    w.Header().Add("Content-Type", "application/json")
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        respBody := errorResponse{
            Error: "Couldn't decode parameters",
        }
        data, err := json.Marshal(respBody)
        if err != nil {
            return
        }
        w.Write(data)
        return
    }
    
    if len(params.Body) < 141 {
        w.WriteHeader(http.StatusOK)
        respBody := returnVals{
            CleanedBody: cleanChirp(params.Body),
        }
        data, err := json.Marshal(respBody)
        if err != nil {
            return
        }
        w.Write(data)
    } else {
        w.WriteHeader(http.StatusBadRequest)
        respBody := errorResponse{
            Error: "Chirp is too long",
        }
        data, err := json.Marshal(respBody)
        if err != nil {
            return
        }
        w.Write(data)
    }
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