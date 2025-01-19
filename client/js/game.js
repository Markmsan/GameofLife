let start = false;
let intervalId = null; 
let n = 0


const canvas = document.getElementById('gameCanvas');
const ctx = canvas.getContext('2d');
// Canvas settings
let gridWidth = 10; // Default grid width (columns)
let gridHeight = 10; // Default grid height (rows)
let cellWidth = canvas.width / gridWidth; // Cell width
let cellHeight = canvas.height / gridHeight; // Cell height

let grid; // Grid will be initialized based on game type (new or loaded)

// Determine if the game is new or loaded
const loaded = localStorage.getItem("loaded") === "true";
console.log(loaded)

if (loaded) {
    // Load grid from localStorage
    grid = JSON.parse(localStorage.getItem("grid"));
    gridHeight = grid.length;
    gridWidth = grid[0].length;

    // Recalculate cell sizes based on the loaded grid dimensions
    cellWidth = canvas.width / gridWidth;
    cellHeight = canvas.height / gridHeight;
} else {
    // Initialize a new grid
    grid = Array.from({ length: gridHeight }, () => Array(gridWidth).fill(0));
}
let newGame = !loaded






function onChange() {
    const button = document.getElementById("toggle_button")
    if (!start) {
        button.textContent = "Stop"; 
        start = true;
        intervalId = setInterval(sendAndReceiveGrid, 1000); // Fetch every second
        sendAndReceiveGrid(); // Initial fetch
        saveGame()
    } 
    else {
        button.textContent = "Start"; 
        clearInterval(intervalId); // Clear the interval
        start = false;
    }
}

async function saveGame() {
    
    const name = localStorage.getItem("gameName");
    try {
        await fetch("http://localhost:8080/savegame", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                name: name,
                grid: grid,
                height: grid.length,
                length: grid[0].length,
            }),
        });

    } catch (error) {
        console.error('Error sending or receiving grid:', error);
    }
}

async function reloadGrid(){
    const name = localStorage.getItem("gameName");
    try {
        const response = await fetch('http://localhost:8080/grid', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                name:name,
                 
            }),
        });

        if (!response.ok) throw new Error('Failed to send grid');

        const data = await response.json();
        grid = data.updatedGrid;
        newGame = false; // Set to false after the first request
        drawGrid();
    } catch (error) {
        console.error('Error sending or receiving grid:', error);
    }
}

async function sendAndReceiveGrid() {
    try {
        const response = await fetch('http://localhost:8080/grid', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                grid,
                newGame, 
            }),
        });

        if (!response.ok) throw new Error('Failed to send grid');

        const data = await response.json();
        grid = data.updatedGrid;
        newGame = false; // Set to false after the first request
        drawGrid();
    } catch (error) {
        console.error('Error sending or receiving grid:', error);
    }
}

function updateGridSize() {
    if(!start){
        const widthInput = document.getElementById('grid_width');
        const heightInput = document.getElementById('grid_height');
        const newGridWidth = parseInt(widthInput.value, 10);
        const newGridHeight = parseInt(heightInput.value, 10);

        if (isNaN(newGridWidth) || newGridWidth < 5 || newGridWidth > 50 || 
            isNaN(newGridHeight) || newGridHeight < 5 || newGridHeight > 50) {
            alert("Please enter valid grid sizes between 5 and 50.");
            return;
        }

        // Update grid dimensions
        gridWidth = newGridWidth;
        gridHeight = newGridHeight;

        // Recalculate cell size
        cellWidth = canvas.width / gridWidth;
        cellHeight = canvas.height / gridHeight;

        // Reset the grid
        grid = Array.from({ length: gridHeight }, () => Array(gridWidth).fill(0));
        newGame = true; // Mark as a new game
        drawGrid(); // Redraw with the new grid size
    }
}
function drawGrid() {
    ctx.clearRect(0, 0, canvas.width, canvas.height); // Clear canvas
    for (let row = 0; row < gridHeight; row++) {
        for (let col = 0; col < gridWidth; col++) {
            ctx.fillStyle = grid[row][col] ? 'black' : 'white'; // Black for live cells, white for dead
            ctx.fillRect(col * cellWidth, row * cellHeight, cellWidth, cellHeight);
            ctx.strokeStyle = 'lightgray';
            ctx.strokeRect(col * cellWidth, row * cellHeight, cellWidth, cellHeight);
        }
    }
}
canvas.addEventListener('click', (e) => {
  
    const rect = canvas.getBoundingClientRect();
    const x = e.clientX - rect.left;
    const y = e.clientY - rect.top;

    const row = Math.floor(y / cellHeight);
    const col = Math.floor(x / cellWidth);

    // Toggle the cell state
    grid[row][col] = grid[row][col] ? 0 : 1;
    newGame=true
    drawGrid();
    
});

drawGrid();