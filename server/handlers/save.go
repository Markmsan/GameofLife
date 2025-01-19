package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"server/db"
	"server/models"
)

func SaveorUpdateGameHandler(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set content type header first
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodPost {
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Method not allowed",
			})
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var req models.SaveGameRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("Failed to decode request: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Invalid request body",
			})
			return
		}

		// Validate input
		if req.Name == "" || req.Height <= 0 || req.Length <= 0 || len(req.Grid) != req.Height {
			log.Printf("Invalid input data received: Name='%s', Height=%d, Length=%d, InitialGrid=%v",
				req.Name, req.Height, req.Length, req.Grid)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Invalid input data",
			})
			return
		}

		if err := db.SaveGameToDB(database, req); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Failed to save game",
			})
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"message": "Game saved successfully",
		})
	}
}
