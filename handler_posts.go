package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/ahgr3y/blog-aggregator/internal/database"
)

func (cfg apiConfig) handlerGetPostsByUser(w http.ResponseWriter, r *http.Request, u database.User) {

	// Get 'limit' query parameter
	limitString := r.URL.Query().Get("limit")
	if limitString == "" {
		limitString = "10"
	}
	limit, err := strconv.Atoi(limitString)
	if err != nil {
		limit = 10
	}

	dbPosts, err := cfg.DB.GetPostsByUser(r.Context(), database.GetPostsByUserParams{
		UserID: u.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		log.Printf("Error from GetPostsByUser: %s", err)
		respondWithError(w, http.StatusBadRequest, "Unable to get posts")
		return
	}

	respondWithJSON(w, http.StatusOK, databasePostsToPosts(dbPosts))
}
