package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/ahgr3y/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func (cfg apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, u database.User) {

	type parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	// Parse request to parameters
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	feedID := database.GenerateUUID()

	feedParams := database.CreateFeedParams{
		ID:        feedID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.Url,
		UserID:    u.ID,
	}

	feed, err := cfg.DB.CreateFeed(r.Context(), feedParams)
	if err != nil {
		log.Printf("Error creating feed: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	type respBody struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Name      string    `json:"name"`
		Url       string    `json:"url"`
		UserID    uuid.UUID `json:"user_id"`
	}

	respondWithJSON(w, http.StatusOK, respBody{
		ID:        feed.ID,
		CreatedAt: feed.CreatedAt,
		UpdatedAt: feed.UpdatedAt,
		Name:      feed.Name,
		Url:       feed.Url,
		UserID:    feed.UserID,
	})
}

func (cfg apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {

	dbFeeds, err := cfg.DB.GetFeeds(r.Context())
	if err != nil {
		log.Printf("Error collecting feeds: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't get feeds")
		return
	}

	feeds := []database.Feed{}
	feeds = append(feeds, dbFeeds...)

	respondWithJSON(w, http.StatusOK, feeds)
}
