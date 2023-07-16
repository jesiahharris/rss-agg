package main

import (
	"fmt"
	"net/http"

	"github.com/jesiahharris/rss-agg/internal/database"
	"github.com/jesiahharris/rss-agg/internal/database/auth"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

// takes an authedHandler and returns handler func for use in router
func (apiCfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetApiKey(r.Header)
		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("Auth error: %v", err))
			return
		}

		user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Couldn't get user: %v", err))
		}
		handler(w, r, user)
	}
}
