package main

import(
	"time"
	"net/http"
	"database/sql"

	"github.com/JoaoPrata/chirpy/internal/auth"
	"github.com/JoaoPrata/chirpy/internal/database"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find token", err)
		return
	}

	_, err = cfg.db.UpdateRefreshToken(r.Context(), database.UpdateRefreshTokenParams{
		Token: refreshToken,
		RevokedAt: sql.NullTime{
			Time: time.Now().UTC(),
			Valid: true,
		},
	})
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't revoke session", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}