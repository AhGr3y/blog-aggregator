package main

import "net/http"

// handlerCheckReadiness responds with a http.StatusOK
// and a 'ok' message to indicate server readiness.
func handlerCheckReadiness(w http.ResponseWriter, r *http.Request) {

	type respBody struct {
		Status string `json:"status"`
	}

	respondWithJSON(w, http.StatusOK, respBody{
		Status: "ok",
	})
}

// handlerTestErrorResponse responds with a http.StatusServerInternalError
// and a error message to indicate error.
func handlerTestErrorResponse(w http.ResponseWriter, r *http.Request) {

	respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
}
