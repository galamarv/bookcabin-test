# --- Stage 1: Build the Go application ---
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go application.
# -o /voucher-app places the output binary in the root of the temporary container.
RUN go build -o /voucher-app main.go

# --- Stage 2: Create the final, minimal image ---
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the pre-built binary from the builder stage
COPY --from=builder /voucher-app .

# Copy the db directory which contains the schema.sql.
# The actual database file will be managed by a Docker volume.
COPY db ./db

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the application
CMD ["./voucher-app"]