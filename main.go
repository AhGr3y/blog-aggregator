package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	// Load environment variables
	// by default, gotdotenv will look for a file named .env
	// in the current directory
	godotenv.Load()

	// Load port number
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}

	// Create a ServeMux
	serveMux := http.NewServeMux()

	// Set handler to test server endpoints
	serveMux.HandleFunc("GET /v1/healthz", handlerCheckReadiness)
	serveMux.HandleFunc("GET /v1/err", handlerTestErrorResponse)

	// Create a Server
	server := http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}

	// Start server
	fmt.Printf("Starting server on port: %s", port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
