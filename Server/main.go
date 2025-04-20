package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"go_src/src"
	"github.com/rs/cors"
)

func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func createDB(w http.ResponseWriter, r *http.Request) {
	dbName := r.URL.Query().Get("dbName")
	if dbName == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "dbName is required"})
		return
	}

	err := src.CreateDatabase(dbName)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Error creating database: %v", err),
		})
		return
	}

	log.Printf("Database %s created successfully\n", dbName)
	jsonResponse(w, http.StatusOK, map[string]string{
		"message": fmt.Sprintf("Database %s created successfully", dbName),
	})
}

func createTable(w http.ResponseWriter, r *http.Request) {
	dbName := r.URL.Query().Get("dbName")
	tableName := r.URL.Query().Get("tableName")
	columns := r.URL.Query()["columns"]
	
	if dbName == "" || tableName == "" || len(columns) == 0 {
		jsonResponse(w, http.StatusBadRequest, map[string]string{
			"error": "dbName, tableName, and columns are required",
		})
		return
	}

	err := src.CreateTable(dbName, tableName, columns)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Error creating table: %v", err),
		})
		return
	}

	jsonResponse(w, http.StatusOK, map[string]string{
		"message": fmt.Sprintf("Table %s created in database %s", tableName, dbName),
	})
}

func readDatabases(w http.ResponseWriter, r *http.Request) {
	dbs, err := src.ReadDatabases()
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Error reading databases: %v", err),
		})
		return
	}

	jsonResponse(w, http.StatusOK, dbs)
}

func readTables(w http.ResponseWriter, r *http.Request) {
	dbName := r.URL.Query().Get("dbName")
	if dbName == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]string{
			"error": "dbName is required",
		})
		return
	}

	tables, err := src.ReadTables(dbName)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Error reading tables: %v", err),
		})
		return
	}

	jsonResponse(w, http.StatusOK, tables)
}

func updateDatabase(w http.ResponseWriter, r *http.Request) {
	oldName := r.URL.Query().Get("oldName")
	newName := r.URL.Query().Get("newName")

	if oldName == "" || newName == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]string{
			"error": "Both oldName and newName are required",
		})
		return
	}

	err := src.UpdateDatabase(oldName, newName)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Error updating database: %v", err),
		})
		return
	}

	jsonResponse(w, http.StatusOK, map[string]string{
		"message": fmt.Sprintf("Database renamed from %s to %s", oldName, newName),
	})
}

func updateTable(w http.ResponseWriter, r *http.Request) {
	dbName := r.URL.Query().Get("dbName")
	tableName := r.URL.Query().Get("tableName")
	columns := r.URL.Query()["columns"]

	if dbName == "" || tableName == "" || len(columns) == 0 {
		jsonResponse(w, http.StatusBadRequest, map[string]string{
			"error": "dbName, tableName, and columns are required",
		})
		return
	}

	err := src.UpdateTable(dbName, tableName, columns)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Error updating table: %v", err),
		})
		return
	}

	jsonResponse(w, http.StatusOK, map[string]string{
		"message": fmt.Sprintf("Table %s updated in database %s", tableName, dbName),
	})
}

func deleteDatabase(w http.ResponseWriter, r *http.Request) {
	dbName := r.URL.Query().Get("dbName")
	if dbName == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]string{
			"error": "dbName is required",
		})
		return
	}

	err := src.DeleteDatabase(dbName)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Error deleting database: %v", err),
		})
		return
	}

	jsonResponse(w, http.StatusOK, map[string]string{
		"message": fmt.Sprintf("Database %s deleted successfully", dbName),
	})
}

func deleteTable(w http.ResponseWriter, r *http.Request) {
	dbName := r.URL.Query().Get("dbName")
	tableName := r.URL.Query().Get("tableName")

	if dbName == "" || tableName == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]string{
			"error": "dbName and tableName are required",
		})
		return
	}

	err := src.DeleteTable(dbName, tableName)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Error deleting table: %v", err),
		})
		return
	}

	jsonResponse(w, http.StatusOK, map[string]string{
		"message": fmt.Sprintf("Table %s deleted from database %s", tableName, dbName),
	})
}

func insertRecord(w http.ResponseWriter, r *http.Request) {
	dbName := r.URL.Query().Get("dbName")
	tableName := r.URL.Query().Get("tableName")
	record := r.URL.Query()["record"]

	if dbName == "" || tableName == "" || len(record) == 0 {
		jsonResponse(w, http.StatusBadRequest, map[string]string{
			"error": "dbName, tableName, and record are required",
		})
		return
	}

	err := src.InsertRecord(dbName, tableName, record)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Error inserting record: %v", err),
		})
		return
	}

	jsonResponse(w, http.StatusOK, map[string]string{
		"message": fmt.Sprintf("Record inserted into table %s in database %s", tableName, dbName),
	})
}

func handleLoadDatabase(w http.ResponseWriter, r *http.Request) {
	dbName := r.URL.Query().Get("dbName")
	if dbName == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]string{
			"error": "dbName is required",
		})
		return
	}

	tables, err := src.LoadFromDB(dbName)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Error loading database: %v", err),
		})
		return
	}

	jsonResponse(w, http.StatusOK, tables)
}

func main() {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/create_db", createDB)
	mux.HandleFunc("/create_table", createTable)
	mux.HandleFunc("/read_databases", readDatabases)
	mux.HandleFunc("/read_tables", readTables)
	mux.HandleFunc("/update_db", updateDatabase)
	mux.HandleFunc("/update_table", updateTable)
	mux.HandleFunc("/delete_db", deleteDatabase)
	mux.HandleFunc("/delete_table", deleteTable)
	mux.HandleFunc("/insert_record", insertRecord)
	mux.HandleFunc("/load-database", handleLoadDatabase)

	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", c.Handler(mux)))
}