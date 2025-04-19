package go_src

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Path where all database binary files will be stored
const dbFolder = "./data"

// Ensure database folder exists
func init() {
	if _, err := os.Stat(dbFolder); os.IsNotExist(err) {
		os.Mkdir(dbFolder, os.ModePerm)
	}
}

// CreateDatabase creates a new binary file for the database
func CreateDatabase(name string) error {
	path := filepath.Join(dbFolder, name+".db")

	if _, err := os.Stat(path); err == nil {
		return errors.New("database already exists")
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Println("Database created:", path)
	return nil
}

// ReadDatabases lists all existing databases
func ReadDatabases() ([]string, error) {
	files, err := os.ReadDir(dbFolder)
	if err != nil {
		return nil, err
	}

	var dbs []string
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".db" {
			dbs = append(dbs, file.Name())
		}
	}
	return dbs, nil
}

// UpdateDatabase renames an existing database
func UpdateDatabase(oldName, newName string) error {
	oldPath := filepath.Join(dbFolder, oldName+".db")
	newPath := filepath.Join(dbFolder, newName+".db")

	if _, err := os.Stat(oldPath); os.IsNotExist(err) {
		return errors.New("old database does not exist")
	}
	if _, err := os.Stat(newPath); err == nil {
		return errors.New("new database name already exists")
	}

	return os.Rename(oldPath, newPath)
}

// DeleteDatabase removes a database file
func DeleteDatabase(name string) error {
	path := filepath.Join(dbFolder, name+".db")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return errors.New("database does not exist")
	}

	return os.Remove(path)
}
