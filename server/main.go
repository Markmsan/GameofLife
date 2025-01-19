package main

import (
	"log"
	"net/http"
	"server/db"
	"server/handlers"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Initialize the database
	database, err := db.InitDB("./life.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Middleware to handle CORS
	corsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			if r.Method == http.MethodOptions {
				log.Printf("Received OPTIONS request for %s", r.URL.Path)
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	}

	mux := http.NewServeMux()

	// Set up HTTP handlers
	mux.HandleFunc("/savegame", handlers.SaveorUpdateGameHandler(database))

	mux.HandleFunc("/load", handlers.LoadGamesHandler(database))

	mux.Handle("/loadGame", handlers.LoadGameHandler(database))

	mux.Handle("/grid", handlers.NextStateHandler(database))

	// Start the server
	const addr = ":8080"
	err = http.ListenAndServe(addr, corsMiddleware(mux))
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
