package go_src

import (
	"encoding/gob"
	"fmt"
	"os"
	"sync"
	"log"
)
// Must be capitalized to export
var Tables map[string]Table

type Table struct {
	Name    string   `json:"name"`
	Columns []string `json:"columns"`
	Records [][]string    `json:"rows"` // assuming you have a Row struct
}

// In-memory map for storing tables
var (
	tables   = make(map[string]Table)
	mu       sync.Mutex
	// dbFile   string 
	dbFile = "data/renameddb.db"
)

// SetDatabaseFile sets the database file path that the user provides
func SetDatabaseFile(filePath string) {
	dbFile = filePath
}

// CreateTable creates a new table with user-defined columns
func CreateTable(tableName string, columns []string) error {
	mu.Lock()
	defer mu.Unlock()

	if dbFile == "" {
		return fmt.Errorf("database file path is not set")
	}

	// Check if table already exists
	if _, exists := tables[tableName]; exists {
		return fmt.Errorf("table %s already exists", tableName)
	}

	// Create new table with the user-defined columns
	tables[tableName] = Table{Name: tableName, Columns: columns}
	log.Printf("Tables: %v\n", tables)
	return saveToDB()
}

// ReadTables retrieves the list of table names
func ReadTables() ([]string, error) {
	mu.Lock()
	defer mu.Unlock()

	var tableNames []string
	for tableName := range tables {
		tableNames = append(tableNames, tableName)
	}
	return tableNames, nil
}

// UpdateTable updates columns of an existing table
func UpdateTable(tableName string, newColumns []string) error {
	mu.Lock()
	defer mu.Unlock()

	if table, exists := tables[tableName]; exists {
		table.Columns = newColumns
		tables[tableName] = table
		return saveToDB()
	}
	return fmt.Errorf("table %s not found", tableName)
}

// DeleteTable removes a table from the map
func DeleteTable(tableName string) error {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := tables[tableName]; exists {
		delete(tables, tableName)
		return saveToDB()
	}
	return fmt.Errorf("table %s not found", tableName)
}

func saveToDB() error {
    log.Printf("Saving tables to dbFile: %s", dbFile)

    // Open the file in read-write mode or create it if it doesn't exist
    file, err := os.OpenFile(dbFile, os.O_RDWR|os.O_CREATE, 0666)
    if err != nil {
        return fmt.Errorf("failed to open db file: %v", err)
    }
    defer file.Close()

    // Create a new Gob encoder
    encoder := gob.NewEncoder(file)

    // Encode the tables map and write it to the file
    err = encoder.Encode(tables)
    if err != nil {
        return fmt.Errorf("failed to encode tables: %v", err)
    }

    log.Println("Tables successfully saved to file")
    return nil
}

func LoadFromDB() error {
    file, err := os.OpenFile(dbFile, os.O_RDWR, 0666)
    if err != nil {
        return fmt.Errorf("failed to open db file: %v", err)
    }
    defer file.Close()

    decoder := gob.NewDecoder(file)

    // Decode into the global variable!
    err = decoder.Decode(&Tables)
    if err != nil {
        return fmt.Errorf("failed to decode tables: %v", err)
    }

    log.Printf("Tables loaded: %v", Tables)
    return nil
}


