package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/ahgr3y/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

// handlerCreateFeedFollow -
func (cfg apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, u database.User) {

	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	// Extract JSON from request body into parameters struct
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding JSON: %s", err)
		respondWithError(w, http.StatusBadRequest, "Missing or invalid request payload")
		return
	}

	// Check if FeedID is a valid id
	_, err = cfg.DB.GetFeedByID(r.Context(), params.FeedID)
	if err != nil {
		log.Printf("Error getting feed by id: %s", err)
		respondWithError(w, http.StatusBadRequest, "Feed does not exist")
		return
	}

	feedID := database.GenerateUUID()

	feedFollowParams := database.CreateFeedFollowParams{
		ID:        feedID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    u.ID,
		FeedID:    params.FeedID,
	}

	feedFollow, err := cfg.DB.CreateFeedFollow(r.Context(), feedFollowParams)
	if err != nil {
		log.Printf("Error creating feed follow: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	type respBody struct {
		ID        uuid.UUID `json:"id"`
		FeedID    uuid.UUID `json:"feed_id"`
		UserID    uuid.UUID `json:"user_id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	respondWithJSON(w, http.StatusOK, respBody{
		ID:        feedFollow.ID,
		FeedID:    feedFollow.FeedID,
		UserID:    feedFollow.UserID,
		CreatedAt: feedFollow.CreatedAt,
		UpdatedAt: feedFollow.UpdatedAt,
	})
}

// handlerDeleteFeedByID -
func (cfg apiConfig) handlerDeleteFeedByID(w http.ResponseWriter, r *http.Request, u database.User) {

	// Extract the feed follows id to delete from URL
	feedFollowIDString := r.PathValue("feedFollowID")
	if feedFollowIDString == "" {
		log.Printf("Invalid/missing feed follow id: %s", feedFollowIDString)
		respondWithError(w, http.StatusBadRequest, "Inavlid or missing feed follow in url path")
		return
	}

	// Convert string to UUID
	feedFollowID, err := uuid.Parse(feedFollowIDString)
	if err != nil {
		log.Printf("Error parsing string to UUID: %s", err)
		respondWithError(w, http.StatusBadRequest, "Invalid feed follow id format")
		return
	}

	// Check if feedFollowID exist
	_, err = cfg.DB.GetFeedFollowByID(r.Context(), feedFollowID)
	if err != nil {
		log.Printf("Error getting feed follow: %s", err)
		respondWithError(w, http.StatusBadRequest, "Feed follow does not exist")
		return
	}

	err = cfg.DB.DeleteFeedFollowByID(r.Context(), feedFollowID)
	if err != nil {
		log.Printf("Error deleting feed follow: %s", err)
		respondWithError(w, http.StatusBadRequest, "Unable to delete feed follow")
		return
	}

	w.WriteHeader(http.StatusOK)
}

// handlerGetFeedFollowsByUserID -
func (cfg apiConfig) handlerGetFeedFollowsByUserID(w http.ResponseWriter, r *http.Request, u database.User) {

	dbFeedFollows, err := cfg.DB.GetFeedFollowsByUserID(r.Context(), u.ID)
	if err != nil {
		log.Printf("Error getting feed follow by user id: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	feedFollows := []database.FeedFollow{}
	feedFollows = append(feedFollows, dbFeedFollows...)

	respondWithJSON(w, http.StatusOK, feedFollows)
}
