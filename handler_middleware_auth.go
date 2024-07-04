package main

import (
	"log"
	"net/http"

	"github.com/ahgr3y/blog-aggregator/internal/auth"
	"github.com/ahgr3y/blog-aggregator/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		apiKey, err := auth.ExtractApiKeyFromRequest(r)
		if err != nil {
			log.Printf("Error extracting api key: %s", err)
			respondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		user, err := cfg.DB.GetUser(r.Context(), apiKey)
		if err != nil {
			log.Printf("Error getting user: %s", err)
			respondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		handler(w, r, user)
	})
}
