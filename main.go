package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ahgr3y/blog-aggregator/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// Create struct to store config
type apiConfig struct {
	DB *database.Queries
}

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

	// Load database url
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	// Open connection to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	// Create database queries obj
	dbQueries := database.New(db)

	// Create config
	apiCfg := apiConfig{
		DB: dbQueries,
	}

	// Create a ServeMux
	serveMux := http.NewServeMux()

	// Set handler to test server endpoints
	serveMux.HandleFunc("GET /v1/healthz", handlerCheckReadiness)
	serveMux.HandleFunc("GET /v1/err", handlerTestErrorResponse)

	// Set handler for managing users
	serveMux.HandleFunc("POST /v1/users", apiCfg.handlerCreateUser)
	serveMux.HandleFunc("GET /v1/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))

	// Set handler for managing feeds
	serveMux.HandleFunc("POST /v1/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	serveMux.HandleFunc("GET /v1/feeds", apiCfg.handlerGetFeeds)

	// Set handler for managing feed follows
	serveMux.HandleFunc("POST /v1/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
	serveMux.HandleFunc("DELETE /v1/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedByID))
	serveMux.HandleFunc("GET /v1/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollowsByUserID))

	// Set handler for managing posts
	serveMux.HandleFunc("GET /v1/posts", apiCfg.middlewareAuth(apiCfg.handlerGetPostsByUser))

	// Create a Server
	server := http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}

	// Start scraper
	go startScraper(apiCfg.DB)

	// Start server
	fmt.Printf("Starting server on port: %s", port)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
