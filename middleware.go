package main

import (
	"fmt"
	"net/http"

	"rssagg/internal/auth"
	"rssagg/internal/database"
)

type authenticatedHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiconf *apiConfig) middlewareAuth(handler authenticatedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetApiKey(r.Header)
		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("Auth error: %v", err))
		}

		user, err := apiconf.DB.GetUserByApiKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Could not get user: %v", err))
		}
		handler(w, r, user)
	}
}
