# GameofLife
This project implements Conway's Game of Life as a no-player game using a database-server-client architecture. The application allows users to create and manage games, edit the grid, and simulate the game's evolution according to its rules.

## Installation
1. Clone the Repository:
```bash
git clone https://github.com/your-username/game-of-life.git
cd game-of-life
```
2. Install Dependencies: Ensure you have Go installed. Run:
```bash
go mod tidy
```
3. The database will initialize automatically on the first server run.

## Usage 

1. Navigate to the server directory and run: The server will start on http://localhost:8080.
```bash
cd server
go run main.go
```

2. Open the Client
```bash
Open the client/start.html file in your browser.
```