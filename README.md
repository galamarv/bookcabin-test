# Airline Voucher Seat Assignment App

## 1. Project Overview

This application is a full-stack solution designed for an airline's promotional campaign. It allows airline crew to generate three unique, random seat vouchers for a specific flight. The system ensures that vouchers are not generated more than once for the same flight on the same day and persists all assignments to a SQLite database.

The project consists of two main parts:
* A **Go (GoLang) backend** that provides a REST API for managing voucher logic.
* A **React frontend** that provides a user-friendly interface for interacting with the backend.

## 2. Features

* **Voucher Generation:** Dynamically generates 3 unique seat numbers based on the selected aircraft type (ATR, Airbus 320, Boeing 737 Max).
* **Data Persistence:** Saves all voucher assignment details into a `vouchers.db` SQLite database file.
* **Duplicate Prevention:** The backend ensures that vouchers for a specific flight number and date can only be generated once.
* **Clean, Layered Architecture:** The Go backend is built using a clean, layered architecture (Handler, Usecase, Repository) for better maintainability, testability, and separation of concerns.
* **User-Friendly Interface:** The React frontend allows for easy input of crew and flight details.
* **Date Validation:** The frontend prevents users from selecting a date in the past.
* **Containerized Deployment:** The entire application (backend and frontend) can be easily run using Docker and Docker Compose for consistent and portable deployment.

## 3. Architecture Overview

The backend is designed with a clear separation of concerns into distinct layers:

* **`model`**: Defines the core data structures (e.g., `Voucher` struct) used across the application.
* **`repository`**: The data access layer. Its sole responsibility is to communicate directly with the database (executing SQL queries). It knows nothing about business rules.
* **`usecase`**: The business logic layer. It orchestrates the application's features by using the repository. For example, the `GenerateVouchers` use case first checks for existence, then generates seats, and finally saves the result.
* **`handler`**: The presentation layer. It handles incoming HTTP requests, decodes JSON bodies, calls the appropriate use case, and formats the HTTP response.

This layered approach makes the application robust, scalable, and easy to test.

### API Endpoints

The backend exposes the following REST API endpoints:

* **`POST /api/check`**
  * **Purpose:** Checks if vouchers have already been generated for a specific flight.
  * **Request Body:** `{"flightNumber": "GA102", "date": "2025-07-15"}`
  * **Success Response:** `{"exists": true}` or `{"exists": false}`

* **`POST /api/generate`**
  * **Purpose:** Generates and saves 3 new seat vouchers if they don't already exist.
  * **Request Body:** `{"name": "John Doe", "id": "98123", "flightNumber": "GA102", "date": "2025-07-15", "aircraft": "Airbus 320"}`
  * **Success Response:** `{"success": true, "seats": ["3B", "14D", "21A"]}`
  * **Error Response (Conflict):** `{"success": false, "error": "vouchers already generated for this flight/date"}`

## 4. How to Run the Application

There are two ways to run this project: using Docker (recommended for ease of use) or running each service locally.

### Method 1: Running with Docker (Recommended)

This method uses Docker Compose to build and run the backend and frontend containers in an isolated environment.

#### Prerequisites
* [Docker Desktop](https://www.docker.com/products/docker-desktop/) installed and running.

#### Instructions
1.  **Navigate to the project's root directory** (`bookcabin-test/`) in your terminal.
2.  **Build and Start the Containers:**
    Run the following command. This will build the Docker images for both services and start them.
    ```bash
    docker-compose up --build
    ```
    You will see logs from both the backend and frontend services in your terminal.
3.  **Access the Application:**
    Once the containers are running, open your web browser and go to:
    **[http://localhost:3000](http://localhost:3000)**
4.  **Stopping the Application:**
    To stop and remove the containers, press `Ctrl + C` in the terminal where `docker-compose` is running, and then execute:
    ```bash
    docker-compose down
    ```

### Method 2: Running Locally (Without Docker)

This method requires you to run the backend and frontend servers in two separate terminal windows.

#### Prerequisites
* [Go](https://golang.org/doc/install) (version 1.22 or newer)
* [Node.js and npm](https://nodejs.org/en/download/) (version 16 or newer)
* [Make](https://www.gnu.org/software/make/) (optional, for using Makefile shortcuts)

#### Instructions

**Terminal 1: Run the Go Backend**
1.  Navigate to the backend directory.
    ```bash
    cd backend-go
    ```
2.  Ensure all Go dependencies are downloaded.
    ```bash
    go mod tidy
    ```
3.  Start the backend server.
    ```bash
    go run main.go
    ```
    The server will start on `http://localhost:8080`. You will see log messages in this terminal.

**Terminal 2: Run the React Frontend**
1.  Navigate to the frontend directory.
    ```bash
    cd frontend-react
    ```
2.  Install the necessary npm packages (only needs to be done once).
    ```bash
    npm install
    ```
3.  Start the React development server.
    ```bash
    npm start
    ```
    A new browser tab should automatically open to **[http://localhost:3000](http://localhost:3000)**. The frontend is configured to proxy API requests to the backend server.

#### Alternative: Using the Makefile
If you have `make` installed, you can use the provided `Makefile` for simpler commands. From the project's root directory:

* **To run the backend server (Terminal 1):**
    ```bash
    make run-backend
    ```
* **To run the frontend server (Terminal 2):**
    ```bash
    make run-frontend
    ```

## 5. Project Structure
```
bookcabin-test/
├── backend-go/
│   ├── handler/
│   │   └── handler.go
│   ├── model/
│   │   └── voucher.go
│   ├── repository/
│   │   └── voucher_repository.go
│   ├── usecase/
│   │   └── voucher_usecase.go
│   ├── db/
│   │   └── schema.sql
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
│   └── Dockerfile
│
├── frontend-react/
│   ├── public/
│   ├── src/
│   │   ├── App.css
│   │   └── App.js
│   ├── package.json
│   ├── Dockerfile
│   └── nginx.conf
│
├── .dockerignore
├── docker-compose.yml
├── Makefile
└── README.md
```