package main

import (
	"bookcabin-test/backend-go/handler"
	"bookcabin-test/backend-go/repository"
	"bookcabin-test/backend-go/usecase"
	"log"
	"net/http"
	"time"
)

// main initializes the application, sets up the database, and starts the HTTP server.
// It follows the Clean Architecture principles by separating concerns into data, business, and presentation layers.
func main() {
	// 1. Data Layer: Initialize database
	db, err := repository.InitDB("./db/vouchers.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// 2. Data Layer: Create the repository
	voucherRepo := repository.NewVoucherRepository(db)

	// 3. Business Layer: Create the usecase, injecting the repository
	voucherUsecase := usecase.NewVoucherUsecase(voucherRepo)

	// 4. Presentation Layer: Create the API handler, injecting the usecase
	apiHandler := handler.NewHandler(voucherUsecase)

	// 5. Setup Router and Server
	mux := http.NewServeMux()
	mux.HandleFunc("/api/check", apiHandler.CheckVouchers)
	mux.HandleFunc("/api/generate", apiHandler.GenerateVouchers)

	// Handle undefined paths
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// This check ensures that only truly undefined paths are logged as errors.
		if r.URL.Path != "/api/check" && r.URL.Path != "/api/generate" {
			log.Printf("ERROR: Invalid path accessed: %s", r.URL.Path)
			http.NotFound(w, r)
		}
	})

	// 6. Start the HTTP server with logging middleware
	log.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", loggingMiddleware(mux)); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// loggingMiddleware logs all incoming requests.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("--> %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("<-- %s %s (%s)", r.Method, r.URL.Path, time.Since(start))
	})
}
