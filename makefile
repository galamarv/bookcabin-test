# Makefile for the Airline Voucher Seat Assignment App

# --- Variables ---
# Define the name for the backend binary
BINARY_NAME=voucher-app

# --- Go Backend Commands ---

## run-backend: Runs the Go backend server in development mode.
run-backend:
	@echo "Starting Go backend server on http://localhost:8080..."
	@cd backend-go && go run main.go

## build-backend: Compiles the Go application into a binary.
build-backend:
	@echo "Building Go backend application..."
	@cd backend-go && go build -o ../$(BINARY_NAME) main.go
	@echo "Backend binary created: $(BINARY_NAME)"

# --- React Frontend Commands ---

## setup-frontend: Installs npm dependencies for the frontend.
setup-frontend:
	@echo "Installing frontend dependencies..."
	@cd frontend-react && npm install

## run-frontend: Runs the React frontend in development mode.
run-frontend:
	@echo "Starting React frontend server on http://localhost:3000..."
	@cd frontend-react && npm start

## build-frontend: Creates a production build of the React app.
build-frontend:
	@echo "Building React frontend for production..."
	@cd frontend-react && npm run build

# --- General Commands ---

## run: Runs both backend and frontend servers concurrently (requires two terminals).
run:
	@echo "This command requires two separate terminals."
	@echo "In one terminal, run: make run-backend"
	@echo "In another terminal, run: make run-frontend"

## build: Builds both the backend and frontend for production.
build: build-backend build-frontend

## clean: Removes generated files (binary and build artifacts).
clean:
	@echo "Cleaning up build artifacts..."
	@rm -f $(BINARY_NAME)
	@rm -rf frontend-react/build
	@echo "Cleanup complete."

## help: Shows this help message.
help:
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: run-backend build-backend setup-frontend run-frontend build-frontend run build clean help