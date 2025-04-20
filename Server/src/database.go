package src

import (
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"syscall"
)

const dbFolder = "./data"

var (
	dbMutex  sync.Mutex
	activeDB *Database
)

type Database struct {
	Name   string
	Tables map[string]Table
}

type Table struct {
	Name    string
	Columns []string
	Rows    [][]interface{}
}

func init() {
	if err := os.MkdirAll(dbFolder, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create data directory: %v", err))
	}
}

// Windows-safe file creation
func createDBFile(path string) (*os.File, error) {
	// First try normal creation
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err == nil {
		return file, nil
	}

	// If permission error, try to take ownership
	if os.IsPermission(err) {
		if err := takeOwnership(path); err != nil {
			return nil, fmt.Errorf("failed to take ownership: %v", err)
		}
		return os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	}

	return nil, err
}

// Windows ownership helper
func takeOwnership(path string) error {
	var sid *syscall.SID
	err := syscall.AllocateAndInitializeSid(
		&syscall.SECURITY_NT_AUTHORITY,
		2,
		syscall.SECURITY_BUILTIN_DOMAIN_RID,
		syscall.DOMAIN_ALIAS_RID_ADMINS,
		0, 0, 0, 0, 0, 0,
		&sid)
	if err != nil {
		return fmt.Errorf("SID Error: %v", err)
	}

	err = syscall.SetNamedSecurityInfo(
		path,
		syscall.SE_FILE_OBJECT,
		syscall.OWNER_SECURITY_INFORMATION,
		sid,
		nil,
		nil,
		nil)
	
	return err
}

func CreateDatabase(name string) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	path := filepath.Join(dbFolder, name+".db")
	
	file, err := createDBFile(path)
	if err != nil {
		if os.IsExist(err) {
			return errors.New("database already exists")
		}
		return fmt.Errorf("failed to create database file: %v", err)
	}
	defer file.Close()

	db := Database{
		Name:   name,
		Tables: make(map[string]Table),
	}

	if err := gob.NewEncoder(file).Encode(db); err != nil {
		os.Remove(path) // Clean up if encoding fails
		return fmt.Errorf("failed to initialize database: %v", err)
	}

	activeDB = &db
	fmt.Printf("Database created: %s\n", path)
	return nil
}

func saveDB(db Database) error {
	tmpPath := filepath.Join(dbFolder, db.Name+".tmp")
	finalPath := filepath.Join(dbFolder, db.Name+".db")

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
	if err := os.Rename(tmpPath, finalPath); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to commit database: %v", err)
	}

	activeDB = &db
	return nil
}

func loadDB(name string) (*Database, error) {
	path := filepath.Join(dbFolder, name+".db")
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var db Database
	if err := gob.NewDecoder(file).Decode(&db); err != nil {
		return nil, err
	}

	activeDB = &db
	return &db, nil
}

// SetActiveDB ensures proper database is loaded
func SetActiveDB(name string) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	db, err := loadDB(name)
	if err != nil {
		return err
	}
	activeDB = db
	return nil
}

// [Rest of your functions (ReadDatabases, UpdateDatabase, DeleteDatabase) remain the same]