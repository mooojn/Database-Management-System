package src

import (
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"syscall"
)

var (
	tableMu sync.Mutex
)

type Table struct {
	Name    string     `json:"name"`
	Columns []string   `json:"columns"`
	Records [][]string `json:"rows"`
}

// SetActiveDB must be called before any table operations
func SetActiveDB(dbName string) error {
	tableMu.Lock()
	defer tableMu.Unlock()

	path := filepath.Join("data", dbName+".db")
	
	// Convert to absolute path for Windows reliability
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %v", err)
	}

	// Verify/create directory with proper permissions
	if err := os.MkdirAll(filepath.Dir(absPath), 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %v", err)
	}

	// Verify we can access the file
	if err := verifyFileAccess(absPath); err != nil {
		return fmt.Errorf("database access error: %v", err)
	}

	return nil
}

func verifyFileAccess(path string) error {
	// Try opening the file
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		if os.IsNotExist(err) {
			// Create new file if it doesn't exist
			file, err = os.Create(path)
			if err != nil {
				return fmt.Errorf("failed to create database file: %v", err)
			}
			defer file.Close()
			
			// Initialize empty database
			db := Database{
				Name:   filepath.Base(path),
				Tables: make(map[string]Table),
			}
			if err := gob.NewEncoder(file).Encode(db); err != nil {
				return fmt.Errorf("failed to initialize database: %v", err)
			}
			return nil
		}
		
		// Try fixing permissions if access denied
		if os.IsPermission(err) {
			if err := fixWindowsPermissions(path); err != nil {
				return fmt.Errorf("permission fix failed: %v", err)
			}
			return verifyFileAccess(path) // Retry
		}
		return err
	}
	file.Close()
	return nil
}

func fixWindowsPermissions(path string) error {
	// Try taking ownership of the file
	owner, err := syscall.StringToSid("S-1-5-32-544") // Administrators group
	if err != nil {
		return err
	}

	dacl := &syscall.ACL{}
	if err := syscall.SetNamedSecurityInfo(
		path,
		syscall.SE_FILE_OBJECT,
		syscall.OWNER_SECURITY_INFORMATION|syscall.DACL_SECURITY_INFORMATION,
		owner,
		nil,
		dacl,
		nil,
	); err != nil {
		return fmt.Errorf("failed to set Windows permissions: %v", err)
	}
	return nil
}

// All table operations now work with the currently active database
func CreateTable(dbName, tableName string, columns []string) error {
	if tableName == "" {
		return errors.New("table name cannot be empty")
	}
	if len(columns) == 0 {
		return errors.New("table must have at least one column")
	}

	tableMu.Lock()
	defer tableMu.Unlock()

	path := filepath.Join("data", dbName+".db")
	db, err := loadDatabase(path)
	if err != nil {
		return err
	}

	if _, exists := db.Tables[tableName]; exists {
		return fmt.Errorf("table '%s' already exists", tableName)
	}

	db.Tables[tableName] = Table{
		Name:    tableName,
		Columns: columns,
		Records: [][]string{},
	}

	log.Printf("Created table '%s' in database '%s'", tableName, dbName)
	return saveDatabase(path, db)
}

// [Similar modifications for other functions...]
// UpdateTable, DeleteTable, InsertRecord etc. should follow same pattern
func UpdateTable(dbName, tableName string, newColumns []string) error {
    if dbName == "" || tableName == "" {
        return errors.New("database and table names cannot be empty")
    }
    if len(newColumns) == 0 {
        return errors.New("table must have at least one column")
    }

    tableMu.Lock()
    defer tableMu.Unlock()

    path := filepath.Join("data", dbName+".db")
    db, err := loadDatabase(path)
    if err != nil {
        return err
    }

    table, exists := db.Tables[tableName]
    if !exists {
        return fmt.Errorf("table '%s' not found", tableName)
    }

    // Prevent schema changes if table has data
    if len(table.Records) > 0 && !columnsMatch(table.Columns, newColumns) {
        return errors.New("cannot modify columns when table contains data")
    }

    table.Columns = newColumns
    db.Tables[tableName] = table

    log.Printf("Updated table '%s' in database '%s'", tableName, dbName)
    return saveDatabase(path, db)
}

func DeleteTable(dbName, tableName string) error {
    if dbName == "" || tableName == "" {
        return errors.New("database and table names cannot be empty")
    }

    tableMu.Lock()
    defer tableMu.Unlock()

    path := filepath.Join("data", dbName+".db")
    db, err := loadDatabase(path)
    if err != nil {
        return err
    }

    if _, exists := db.Tables[tableName]; !exists {
        return fmt.Errorf("table '%s' not found", tableName)
    }

    delete(db.Tables, tableName)
    log.Printf("Deleted table '%s' from database '%s'", tableName, dbName)
    return saveDatabase(path, db)
}

func InsertRecord(dbName, tableName string, record []string) error {
    if dbName == "" || tableName == "" {
        return errors.New("database and table names cannot be empty")
    }

    tableMu.Lock()
    defer tableMu.Unlock()

    path := filepath.Join("data", dbName+".db")
    db, err := loadDatabase(path)
    if err != nil {
        return err
    }

    table, exists := db.Tables[tableName]
    if !exists {
        return fmt.Errorf("table '%s' not found", tableName)
    }

    if len(record) != len(table.Columns) {
        return fmt.Errorf("record has %d fields, expected %d", len(record), len(table.Columns))
    }

    table.Records = append(table.Records, record)
    db.Tables[tableName] = table

    log.Printf("Inserted record into table '%s' in database '%s'", tableName, dbName)
    return saveDatabase(path, db)
}

func GetTable(dbName, tableName string) (*Table, error) {
    if dbName == "" || tableName == "" {
        return nil, errors.New("database and table names cannot be empty")
    }

    tableMu.Lock()
    defer tableMu.Unlock()

    path := filepath.Join("data", dbName+".db")
    db, err := loadDatabase(path)
    if err != nil {
        return nil, err
    }

    table, exists := db.Tables[tableName]
    if !exists {
        return nil, fmt.Errorf("table '%s' not found", tableName)
    }

    // Return a copy to prevent external modifications
    tableCopy := table
    return &tableCopy, nil
}

func columnsMatch(oldCols, newCols []string) bool {
    if len(oldCols) != len(newCols) {
        return false
    }
    for i := range oldCols {
        if oldCols[i] != newCols[i] {
            return false
        }
    }
    return true
}

func loadDatabase(path string) (*Database, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %v", err)
    }
    defer file.Close()

    var db Database
    if err := gob.NewDecoder(file).Decode(&db); err != nil {
        return nil, fmt.Errorf("failed to decode database: %v", err)
    }
    return &db, nil
}

func saveDatabase(path string, db *Database) error {
    tmpPath := path + ".tmp"
    
    file, err := os.Create(tmpPath)
    if err != nil {
        return fmt.Errorf("failed to create temp file: %v", err)
    }

    err = gob.NewEncoder(file).Encode(db)
    file.Close()
    if err != nil {
        os.Remove(tmpPath)
        return fmt.Errorf("failed to encode database: %v", err)
    }

    if err := os.Rename(tmpPath, path); err != nil {
        os.Remove(tmpPath)
        return fmt.Errorf("failed to commit database: %v", err)
    }
    return nil
}
func ReadTables(dbName string) ([]string, error) {
    if dbName == "" {
        return nil, errors.New("database name cannot be empty")
    }

    tableMu.Lock()
    defer tableMu.Unlock()

    path := filepath.Join("data", dbName+".db")
    db, err := loadDatabase(path)
    if err != nil {
        return nil, err
    }

    tableNames := make([]string, 0, len(db.Tables))
    for name := range db.Tables {
        tableNames = append(tableNames, name)
    }
    return tableNames, nil
}
func loadDatabase(path string) (*Database, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}
	defer file.Close()

	var db Database
	if err := gob.NewDecoder(file).Decode(&db); err != nil {
		return nil, fmt.Errorf("failed to decode database: %v", err)
	}
	return &db, nil
}

func saveDatabase(path string, db *Database) error {
	tmpPath := path + ".tmp"
	
	// 1. Save to temporary file
	file, err := os.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %v", err)
	}

	err = gob.NewEncoder(file).Encode(db)
	file.Close()
	if err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to encode database: %v", err)
	}

	// 2. Atomic rename
	if err := os.Rename(tmpPath, path); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to commit database: %v", err)
	}

	return nil
}