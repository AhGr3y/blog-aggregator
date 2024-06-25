package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// respondWithJSON parses code and payload,
// uses http.ResponseWriter to send a response
// with JSON Content-Type back to the client.
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {

	// Convert payload to JSON
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling to JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the response header
	w.Header().Set("Content-Type", "application/json")

	// Set the response status code
	w.WriteHeader(code)

	// Set the response content
	w.Write(dat)
}

// respondWithError parses code and msg,
// calls respondWithJSON to send a response
// with JSON Content-Type back to the client.
func respondWithError(w http.ResponseWriter, code int, msg string) {

	// Log 5xx errors
	if code > 499 {
		log.Printf("Unexpected error: %d - %s", code, msg)
	}

	type errorResponse struct {
		ErrorMsg string `json:"error"`
	}

	// Respond with JSON
	respondWithJSON(w, code, errorResponse{
		ErrorMsg: msg,
	})
}
