package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"server/models"
	"time"
)

func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Create the numbers table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS games (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL unique,
			grid TEXT NOT NULL,
			height INTEGER NOT NULL,
			length INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP);
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %v", err)
	}

	log.Println("Database initialized and table created (if not existing).")
	return db, nil
}

func SaveGameToDB(db *sql.DB, req models.SaveGameRequest) error {
	gridJSON, err := json.Marshal(req.Grid)
	if err != nil {
		return fmt.Errorf("failed to encode grid to JSON: %v", err)
	}

	// First check if the name exists
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM games WHERE name = ?)", req.Name).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check for existing game: %v", err)
	}

	if exists {
		// Update existing record
		_, err = db.Exec(`
            UPDATE games 
            SET grid = ?, height = ?, length = ?, created_at = ?
            WHERE name = ?
        `, string(gridJSON), req.Height, req.Length, time.Now(), req.Name)
		if err != nil {
			return fmt.Errorf("failed to update existing game: %v", err)
		}
		log.Printf("Game '%s' updated successfully.", req.Name)
	} else {
		// Insert new record
		_, err = db.Exec(`
            INSERT INTO games (name, grid, height, length, created_at)
            VALUES (?, ?, ?, ?, ?)
        `, req.Name, string(gridJSON), req.Height, req.Length, time.Now())
		if err != nil {
			return fmt.Errorf("failed to insert new game: %v", err)
		}
		log.Printf("Game '%s' created successfully.", req.Name)
	}

	return nil
}
