package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

var DB *sql.DB

// Initialize the database connection
func Init() {
	var err error
	// Replace the connection string with your MySQL credentials
	DB, err = sql.Open("mysql", "root:Reddy@123@tcp(localhost:3306)/my_toronto_time")
	if err != nil {
		log.Fatalf("Error opening database connection: %v\n", err)
	}

	// Verify the database connection
	if err := DB.Ping(); err != nil {
		log.Fatalf("Error verifying database connection: %v\n", err)
	}
	log.Println("Successfully connected to MySQL database.")
}

// TimeData struct to structure the JSON response
type TimeData struct {
	CurrentTime string `json:"current_time"`
}

// Handler for the '/current-time' endpoint
func currentTimeHandler(w http.ResponseWriter, r *http.Request) {
	// Get the current time in Toronto (Eastern Standard Time)
	location, err := time.LoadLocation("America/Toronto")
	if err != nil {
		http.Error(w, "Unable to load Toronto time zone", http.StatusInternalServerError)
		return
	}
	torontoTime := time.Now().In(location)

	// Insert the current time into the database
	_, err = DB.Exec("INSERT INTO my_time_log (timestamp) VALUES (?)", torontoTime)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error logging time to database: %v", err), http.StatusInternalServerError)
		return
	}

	// Prepare the response data
	timeData := TimeData{
		CurrentTime: torontoTime.Format("2006-01-02 15:04:05"),
	}

	// Set the response header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Send the JSON response
	if err := json.NewEncoder(w).Encode(timeData); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding JSON response: %v", err), http.StatusInternalServerError)
	}
}

func main() {
	// Initialize the database connection
	Init()

	// Setup API routes
	http.HandleFunc("/current-time", currentTimeHandler)

	// Start the server
	port := ":8080"
	log.Printf("Starting server on http://localhost%s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	}
}
