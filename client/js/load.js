function displayGames() {
    const games = JSON.parse(localStorage.getItem("games")) || [];
    const gamesListDiv = document.getElementById("gamesList");

    if (games.length === 0) {
        gamesListDiv.innerHTML = "<p>No games found.</p>";
        return;
    }

    games.forEach((gameName) => {
        const button = document.createElement("button");
        button.textContent = gameName;
        button.onclick = () => loadGame(gameName); // Function to load the selected game
        gamesListDiv.appendChild(button);
    });
}

// Function to load a selected game
async function loadGame(gameName) {   
    try{
        // get the gamenames from /load and store it in localstorage
        const response = await fetch("http://localhost:8080/loadGame", {
            method: "Post",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                name: gameName,
            }),
        });
        if (!response.ok) {
            throw new Error(`Failed to fetch games: ${response.status}`);
        }

        // Parse the JSON response
        const games = await response.json();

        // Store the games list in localStorage
        localStorage.setItem("grid", JSON.stringify(games.grid));
        localStorage.setItem("loaded", true);
        
        location.href = `game.html?gameName=${encodeURIComponent(gameName)}`;
    } catch (error) {
        console.error("Error fetching games:", error);
        alert("Failed to load games. Please try again.");
    }
  
}

// Run on page load
window.onload = displayGames;