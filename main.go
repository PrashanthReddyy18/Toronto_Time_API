package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// Init function to establish the database connection
func Init() {
	var err error
	// Connect to the MySQL database
	db, err = sql.Open("mysql", "root:yourpassword@tcp(localhost:3306)/my_toronto_time")
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	// Ensure that the database is accessible
	if err := db.Ping(); err != nil {
		log.Fatalf("Error verifying database connection: %v", err)
	}
	log.Println("Database connection established.")
}

// GetCurrentTime retrieves the current time in Toronto and logs it in the database
func GetCurrentTime(w http.ResponseWriter, r *http.Request) {
	// Get current time in Toronto
	location, err := time.LoadLocation("America/Toronto")
	if err != nil {
		log.Printf("Error loading Toronto timezone: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	currentTime := time.Now().In(location)

	// Insert current time into the database
	_, err = db.Exec("INSERT INTO my_time_log (timestamp) VALUES (?)", currentTime)
	if err != nil {
		log.Printf("Error inserting time into database: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Log the successful request
	log.Printf("Inserted time: %v into database.", currentTime)

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"current_time": "%s"}`, currentTime)
}

// GetLogs retrieves all the logged times from the database
func GetLogs(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, timestamp FROM my_time_log")
	if err != nil {
		log.Printf("Error querying logs: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var logs []string
	for rows.Next() {
		var id int
		var timestamp string
		if err := rows.Scan(&id, &timestamp); err != nil {
			log.Printf("Error scanning log row: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		logs = append(logs, fmt.Sprintf("ID: %d, Timestamp: %s", id, timestamp))
	}

	// Handle any error that occurred during iteration
	if err := rows.Err(); err != nil {
		log.Printf("Error while iterating rows: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Send the logs as a response
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"logs": [%s]}`, logs)
}

// main function to initialize and start the server
func main() {
	// Initialize the database connection
	Init()

	// Setup API routes
	http.HandleFunc("/current-time", GetCurrentTime)
	http.HandleFunc("/logs", GetLogs)

	// Start the server
	port := ":8080"
	log.Printf("Starting server on http://localhost%s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	}
}
