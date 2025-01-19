function startNewGame() {
    const gameNameInput = document.getElementById("game_name");
    const gameName = gameNameInput.value.trim(); // Remove leading/trailing spaces
    const errorMessage = document.getElementById("error_message");

    if (gameName === "") {
        // Show an error message if the input is empty
        errorMessage.style.display = "block";
        errorMessage.innerText = "Please enter a valid game name.";
    } else {
        // Hide the error message and navigate to game.html
        errorMessage.style.display = "none";
        localStorage.setItem("gameName", gameName);
        localStorage.setItem("loaded", false);
        location.href = `game.html?gameName=${encodeURIComponent(gameName)}`;
    }
}

async function LoadGames(){
    try{
        // get the gamenames from /load and store it in localstorage
        const response = await fetch("http://localhost:8080/load", {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
            },
        });
        if (!response.ok) {
            throw new Error(`Failed to fetch games: ${response.status}`);
        }

        // Parse the JSON response
        const games = await response.json();

        // Store the games list in localStorage
        localStorage.setItem("games", JSON.stringify(games));
    } catch (error) {
        console.error("Error fetching games:", error);
        alert("Failed to load games. Please try again.");
    }
    location.href = `load.html`;
}