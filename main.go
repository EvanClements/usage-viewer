package main

import (
	"database/sql"
	"encoding/csv"
	"encoding/xml"
	_ "github.com/mattn/go-sqlite3"
	"github.com/wailsio/wails/v2"
	"strings"
)

// YourDataStruct represents the structure of your data
type YourDataStruct struct {
	// Define the fields based on your XML/CSV schema
	Field1 string `json:"field1"`
	Field2 string `json:"field2"`
	// Add more fields as needed
}

// App struct holds the Wails app instance
type App struct {
	runtime *wails.Runtime
	db      *sql.DB
}

// WailsInit is called on application startup
func (a *App) WailsInit(runtime *wails.Runtime) error {
	a.runtime = runtime
	return nil
}

// WailsAppInit initializes the Wails app
func (a *App) WailsAppInit(runtime *wails.Runtime) error {
	a.runtime = runtime

	// Connect to SQLite database
	db, err := sql.Open("sqlite3", "./your_database.db")
	if err != nil {
		return err
	}
	a.db = db

	// Create table if not exists
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS your_table (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			field1 TEXT,
			field2 TEXT
			-- Add more columns as needed
		)
	`)
	if err != nil {
		return err
	}

	return nil
}

// HandleFileUpload is a Wails handler for file upload
func (a *App) HandleFileUpload(fileData []byte, filename string) (string, error) {
	var data []YourDataStruct

	if isXMLFile(filename) {
		err := xml.Unmarshal(fileData, &data)
		if err != nil {
			return "Error parsing XML file", err
		}
	} else if isCSVFile(filename) {
		reader := csv.NewReader(strings.NewReader(string(fileData)))
		records, err := reader.ReadAll()
		if err != nil {
			return "Error parsing CSV file", err
		}

		// Convert CSV data to YourDataStruct format
		// Append to data slice
		// ...

	} else {
		return "Unsupported file type", nil
	}

	err := a.insertDataIntoDatabase(data)
	if err != nil {
		return "Error inserting data into database", err
	}

	return "File successfully uploaded and data inserted into the database", nil
}

// Helper function to check if the file is XML
func isXMLFile(filename string) bool {
	return strings.HasSuffix(filename, ".xml")
}

// Helper function to check if the file is CSV
func isCSVFile(filename string) bool {
	return strings.HasSuffix(filename, ".csv")
}

// Helper function to insert data into SQLite database
func (a *App) insertDataIntoDatabase(data []YourDataStruct) error {
	for _, entry := range data {
		_, err := a.db.Exec(`
			INSERT INTO your_table (field1, field2)
			VALUES (?, ?)
		`, entry.Field1, entry.Field2)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	app := &App{}
	err := wails.Run(&wails.AppConfig{
		Width:     800,
		Height:    600,
		Title:     "Wails with SQLite and Svelte",
		JS:        "./frontend/dist/App.js",
		Resizable: true,
		CSS:       "./frontend/dist/App.css",
		Colour:    "#131313",
	}, app)
	if err != nil {
		panic(err)
	}
}
