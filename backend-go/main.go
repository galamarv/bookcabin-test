package main

import (
	"bookcabin-test/backend-go/handler"
	"bookcabin-test/backend-go/store"
	"log"
	"net/http"
	"time"
)

// loggingMiddleware logs the details of each request and its processing time.
// It is used to track the performance and flow of requests through the application.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("--> %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("<-- %s %s (%s)", r.Method, r.URL.Path, time.Since(start))
	})
}

// router is the main HTTP request handler that routes requests to the appropriate handlers.
// It handles the paths for checking and generating vouchers, and logs errors for any undefined paths.
func router(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/api/check":
		handler.CheckHandler(w, r)
	case "/api/generate":
		handler.GenerateHandler(w, r)
	default:
		// Log an error for any undefined path accessed
		log.Printf("ERROR: Invalid path accessed: %s", r.URL.Path)
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}

// main initializes the database and starts the HTTP server.
// It sets up the logging middleware and registers the main handler for all paths.
func main() {
	store.InitDB("./db/vouchers.db")

	// Create the main HTTP handler with logging middleware.
	// This handler will route requests to the appropriate handlers based on the path.
	mainHandler := http.HandlerFunc(router)

	// Wrap the main handler with the logging middleware and register it for all paths.
	http.Handle("/", loggingMiddleware(mainHandler))

	log.Println("Server starting on port 8080...")
	// Start the HTTP server on port 8080 and log any errors that occur.
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
