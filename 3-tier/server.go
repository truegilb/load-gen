package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	// The standard MySQL driver is fully compatible with MariaDB
	_ "github.com/go-sql-driver/mysql" 
)

// Data structure for the payload
type Payload struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

var db *sql.DB

func main() {
	// 1. Initialize Database Connection
	// NOTE: Replace with your actual MariaDB credentials
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST") // e.g., 127.0.0.1:3306 or your MariaDB host
	dbName := os.Getenv("DB_NAME")

	if dbUser == "" || dbPass == "" || dbHost == "" || dbName == "" {
		log.Fatal("DB environment variables not set. Please set DB_USER, DB_PASS, DB_HOST, DB_NAME.")
	}
	
	// The MariaDB connection string uses the same format as MySQL
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPass, dbHost, dbName)
	
	var err error
	db, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}
	defer db.Close()

	// Ping the database to ensure connection is valid
	if err = db.Ping(); err != nil {
		log.Fatalf("Error connecting to MariaDB database: %v", err)
	}
	log.Println("Successfully connected to MariaDB database!")

	// 2. Setup HTTP routes
	http.HandleFunc("/api/data", postHandler)
	http.HandleFunc("/api/data/", getHandler)

	log.Println("Server starting on :8888...")
	// Start the server
	log.Fatal(http.ListenAndServe(":8888", nil))
}

// Handles POST requests to /api/data
func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is supported.", http.StatusMethodNotAllowed)
		return
	}

	var payload Payload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Database call: Write the payload
	stmt, err := db.Prepare("INSERT INTO payloads (key_data, value_data) VALUES (?, ?)")
	if err != nil {
		http.Error(w, "Database preparation error", http.StatusInternalServerError)
		log.Printf("DB Prepare error: %v", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(payload.Key, payload.Value)
	if err != nil {
		http.Error(w, "Database execution error", http.StatusInternalServerError)
		log.Printf("DB Exec error: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "Payload successfully written."}`))
	log.Printf("Wrote payload with key: %s", payload.Key)
}

// Handles GET requests to /api/data/{key}
func getHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is supported.", http.StatusMethodNotAllowed)
		return
	}

	// Extract the key from the URL path (e.g., /api/data/my-key)
	key := r.URL.Path[len("/api/data/"):]
	if key == "" {
		http.Error(w, "Missing key in path", http.StatusBadRequest)
		return
	}

	var result Payload
	// Database call: Read the payload
	row := db.QueryRow("SELECT key_data, value_data FROM payloads WHERE key_data = ?", key)
	err := row.Scan(&result.Key, &result.Value)

	if err == sql.ErrNoRows {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Database read error", http.StatusInternalServerError)
		log.Printf("DB Scan error: %v", err)
		return
	}

	// Respond with the retrieved data
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
	log.Printf("Read payload with key: %s", result.Key)
}
