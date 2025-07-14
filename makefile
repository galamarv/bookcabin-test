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

## run: Runs both backend and frontend servers.
# Note: This command requires two separate terminals to run the backend and frontend servers concurrently.	
run:
	@echo "This command requires two separate terminals."
	@echo "In one terminal, run: make run-backend"
	@echo "In another terminal, run: make run-frontend"

## build: Builds both backend and frontend applications for production.
build: build-backend build-frontend

##  clean: Cleans up build artifacts and generated files.
clean:
	@echo "Cleaning up build artifacts..."
	@rm -f $(BINARY_NAME)
	@rm -rf frontend-react/build
	@echo "Cleanup complete."

## help: Displays this help message.
help:	
	@echo "Makefile Commands:"
	@echo "  make run-backend       - Start the Go backend server"
	@echo "  make build-backend     - Build the Go backend application"
	@echo "  make setup-frontend    - Install npm dependencies for the frontend"
	@echo "  make run-frontend      - Start the React frontend server"
	@echo "  make build-frontend    - Build the React frontend for production"
	@echo "  make run               - Run both backend and frontend servers (requires two terminals)"
	@echo "  make build             - Build both backend and frontend for production"
	@echo "  make clean             - Remove generated files and build artifacts"
	@echo "  make help              - Display this help message"

# --- Phony Targets ---
# These targets are not actual files, but commands to be executed.
.PHONY: run-backend build-backend setup-frontend run-frontend build-frontend run build clean help