package main

import (
	 "encoding/json"
	"fmt"
	"log"
	"net/http"
	"server/go-src" // Import your package where the database functions are located
	"github.com/rs/cors"
)

// Function to create a database using go_src package
func createDB(w http.ResponseWriter, r *http.Request) {
	// Extract dbName from query parameters
	dbName := r.URL.Query().Get("dbName")
	if dbName == "" {
		http.Error(w, "dbName is required", http.StatusBadRequest)
		return
	}

	// Call the CreateDatabase function from go_src
	err := go_src.CreateDatabase(dbName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating database: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	log.Printf("Database %s created successfully\n", dbName)
	w.Write([]byte(fmt.Sprintf("Database %s created successfully", dbName)))
}

// Function to create a table in an existing database using the go_src/table.go package
func createTable(w http.ResponseWriter, r *http.Request) {
    // Extract dbName and tableName from query parameters
    dbName := r.URL.Query().Get("dbName")
    tableName := r.URL.Query().Get("tableName")
    columns := r.URL.Query()["columns"] 
    if dbName == "" || tableName == "" || len(columns) == 0 {
        http.Error(w, "dbName, tableName, and columns are required", http.StatusBadRequest)
        return
    }

    // Call CreateTable from the go_src/table.go package
    err := go_src.CreateTable(tableName, columns)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error creating table: %v", err), http.StatusInternalServerError)
        return
    }

    // Respond with a success message
    log.Printf("Table %s created in database %s with columns %v\n", tableName, dbName, columns)
    w.Write([]byte(fmt.Sprintf("Table %s created in database %s with columns %v", tableName, dbName, columns)))
}


// Function to read the list of existing databases
func readDatabases(w http.ResponseWriter, r *http.Request) {
	// Call the ReadDatabases function from go_src
	dbs, err := go_src.ReadDatabases()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading databases: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with the list of databases
	w.Write([]byte(fmt.Sprintf("Databases: %v", dbs)))
}

// Function to update (rename) an existing database
func updateDatabase(w http.ResponseWriter, r *http.Request) {
	log.Printf("Table")
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var payload struct {
		OldName string `json:"oldName"`
		NewName string `json:"newName"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if payload.OldName == "" || payload.NewName == "" {
		http.Error(w, "Both oldName and newName are required", http.StatusBadRequest)
		return
	}

	err := go_src.UpdateDatabase(payload.OldName, payload.NewName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating database: %v", err), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(fmt.Sprintf("Database renamed from %s to %s", payload.OldName, payload.NewName)))
}


// Function to delete a database
func deleteDatabase(w http.ResponseWriter, r *http.Request) {
	// Extract dbName from query parameters
	dbName := r.URL.Query().Get("dbName")
	if dbName == "" {
		http.Error(w, "dbName is required", http.StatusBadRequest)
		return
	}

	// Call the DeleteDatabase function from go_src
	err := go_src.DeleteDatabase(dbName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting database: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.Write([]byte(fmt.Sprintf("Database %s deleted successfully", dbName)))
}
func handleLoadDatabase(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	err := go_src.LoadFromDB()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading database: %v", err), http.StatusInternalServerError)
		return
	}

	log.Println("Loaded table data:")
	for name, table := range go_src.Tables {
		log.Printf("Table Name: %s\n", name)
		log.Printf("Columns: %v\n", table.Columns)
	}
	

	// Encode the exported go_src.Tables
	if err := json.NewEncoder(w).Encode(go_src.Tables); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding JSON: %v", err), http.StatusInternalServerError)
	}
}

func main() {
	// Set up CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"}, // Allow frontend from localhost:3000
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	// Create the default HTTP mux (router)
	mux := http.NewServeMux()

	// Define your routes
	mux.HandleFunc("/create_db", createDB)
	mux.HandleFunc("/create_table", createTable)
	mux.HandleFunc("/read_databases", readDatabases)
	mux.HandleFunc("/modify_db", updateDatabase)
	mux.HandleFunc("/delete_db", deleteDatabase)
	mux.HandleFunc("/load-database", handleLoadDatabase)
	// Wrap the mux with the CORS handler
	handler := c.Handler(mux)

	// Start the server
	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", handler)) // Use the CORS handler
}
