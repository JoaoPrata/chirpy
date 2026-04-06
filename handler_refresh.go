package main

import(
	"time"

	"github.com/JoaoPrata/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get refresh token from header", err)
		return
	}

	dbRefreshToken, err := cfg.db.GetRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get refresh token", err)
		return
	}
	if dbRefreshToken.ExpiresAt > time.Now().UTC() {
		respondWithJSON(w, http.StatusUnauthorized, "Token has expired")
		return
	}
	if dbRefreshToken.RevokedAt != nil {
		respondWithJSON(w, http.StatusUnauthorized, "Token is revoked")
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get user", err)
		return
	}
	tokenDuration, err := time.ParseDuration("1h")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't generate token", err)
		return
	}
	token, err := auth.MakeJWT(user.ID, cfg.tokenSecret, tokenDuration)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't generate token", err)
		return
	}
	respondWithJSON(w, http.StatusOK, response{
		Token: token,
	})
}