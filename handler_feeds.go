package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/ahgr3y/blog-aggregator/internal/database"
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

	feedFollowID := database.GenerateUUID()

	feedFollowParams := database.CreateFeedFollowParams{
		ID:        feedFollowID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    u.ID,
		FeedID:    feed.ID,
	}

	feedFollow, err := cfg.DB.CreateFeedFollow(r.Context(), feedFollowParams)
	if err != nil {
		log.Printf("Error creating feed follow: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	type respBody struct {
		Feed       database.Feed       `json:"feed"`
		FeedFollow database.FeedFollow `json:"feed_follow"`
	}

	respondWithJSON(w, http.StatusOK, respBody{
		Feed:       feed,
		FeedFollow: feedFollow,
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
