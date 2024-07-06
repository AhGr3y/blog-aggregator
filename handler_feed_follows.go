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

	feedFollowID := database.GenerateUUID()

	dbFeedFollow, err := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        feedFollowID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    u.ID,
		FeedID:    params.FeedID,
	})
	if err != nil {
		log.Printf("Error creating feed follow: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	respondWithJSON(w, http.StatusOK, databaseFeedFollowToFeedFollow(dbFeedFollow))
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
	feedFollow, err := cfg.DB.GetFeedFollowByID(r.Context(), feedFollowID)
	if err != nil {
		log.Printf("Error getting feed follow: %s", err)
		respondWithError(w, http.StatusBadRequest, "Feed follow does not exist")
		return
	}

	// Only able to delete feed follows of authenticated user
	if feedFollow.UserID != u.ID {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	err = cfg.DB.DeleteFeedFollowByID(r.Context(), database.DeleteFeedFollowByIDParams{
		ID:     feedFollowID,
		UserID: u.ID,
	})
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

	respondWithJSON(w, http.StatusOK, databaseFeedFollowsToFeedFollows(dbFeedFollows))
}
