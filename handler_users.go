package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/ahgr3y/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Name string `json:"name"`
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

	id := database.GenerateUUID()
	currentTime := time.Now().UTC()

	// Save user to database
	_, err = cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        id,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		Name:      params.Name,
	})
	if err != nil {
		log.Printf("Error creating user: %s", err)
		respondWithError(w, http.StatusBadRequest, "Error creating user")
		return
	}

	type responseBody struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Name      string    `json:"name"`
	}

	respondWithJSON(w, http.StatusOK, responseBody{
		ID:        id,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		Name:      params.Name,
	})
}

func (cfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, u database.User) {

	type respBody struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Name      string    `json:"name"`
		ApiKey    string    `json:"api_key"`
	}

	respondWithJSON(w, http.StatusOK, respBody{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Name:      u.Name,
		ApiKey:    u.ApiKey,
	})
}
