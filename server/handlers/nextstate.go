package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

func NextStateHandler(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var request struct {
				Grid    [][]int `json:"grid"`
				NewGame bool    `json:"newGame"`
			}
			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}

			fmt.Printf("Received grid from client\n")

			// Process the grid by incrementing each element by 1
			request.Grid = computeNextGeneration(request.Grid)

			w.WriteHeader(http.StatusOK) // Explicitly set the response status code
			// Respond with the updated grid
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(struct {
				UpdatedGrid [][]int `json:"updatedGrid"`
			}{
				UpdatedGrid: request.Grid,
			})
		}
	}
}
func countLiveNeighbors(grid [][]int, x, y, rows, cols int) int {
	directions := []struct{ dx, dy int }{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}

	liveNeighbors := 0
	for _, d := range directions {
		nx, ny := x+d.dx, y+d.dy
		if nx >= 0 && nx < rows && ny >= 0 && ny < cols && grid[nx][ny] == 1 {
			liveNeighbors++
		}
	}
	return liveNeighbors
}
func computeNextGeneration(grid [][]int) [][]int {
	rows := len(grid)
	cols := len(grid[0])
	nextGrid := make([][]int, rows)
	for i := range nextGrid {
		nextGrid[i] = make([]int, cols)
	}

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			liveNeighbors := countLiveNeighbors(grid, i, j, rows, cols)

			if grid[i][j] == 1 {
				if liveNeighbors < 2 || liveNeighbors > 3 {
					nextGrid[i][j] = 0
				} else {
					nextGrid[i][j] = 1
				}
			} else {
				if liveNeighbors == 3 {
					nextGrid[i][j] = 1
				}
			}
		}
	}

	return nextGrid
}
