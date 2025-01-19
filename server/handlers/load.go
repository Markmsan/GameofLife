package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func LoadGamesHandler(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set response header
		w.Header().Set("Content-Type", "application/json")

		// Only allow GET requests
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Query all game names
		rows, err := database.Query("SELECT name FROM games")
		if err != nil {
			log.Printf("Failed to query game names: %v", err)
			http.Error(w, "Failed to load games", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Collect game names
		var games []string
		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err != nil {
				log.Printf("Failed to scan game name: %v", err)
				http.Error(w, "Failed to load games", http.StatusInternalServerError)
				return
			}
			games = append(games, name)
		}

		// Respond with game names in JSON format
		if err := json.NewEncoder(w).Encode(games); err != nil {
			log.Printf("Failed to encode game names: %v", err)
			http.Error(w, "Failed to encode game names", http.StatusInternalServerError)
		}
	}
}
func LoadGameHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set response headers
		w.Header().Set("Content-Type", "application/json")

		// Only allow POST requests
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse the JSON body to get the game name
		var requestBody struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			log.Printf("Failed to decode request body: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Validate the game name
		if requestBody.Name == "" {
			http.Error(w, "Game name is required", http.StatusBadRequest)
			return
		}
		gameName := requestBody.Name

		// Query the database for the game details
		var gridJSON string
		var height, length int
		err := db.QueryRow(`
			SELECT grid, height, length 
			FROM games 
			WHERE name = ?
		`, gameName).Scan(&gridJSON, &height, &length)

		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Game not found", http.StatusNotFound)
			} else {
				log.Printf("Failed to fetch game: %v", err)
				http.Error(w, "Failed to fetch game", http.StatusInternalServerError)
			}
			return
		}

		// Decode the grid JSON
		var grid [][]int
		err = json.Unmarshal([]byte(gridJSON), &grid)
		if err != nil {
			log.Printf("Failed to decode grid JSON: %v", err)
			http.Error(w, "Failed to decode grid data", http.StatusInternalServerError)
			return
		}

		// Respond with the game details
		response := map[string]interface{}{
			"name":   gameName,
			"grid":   grid,
			"height": height,
			"length": length,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
