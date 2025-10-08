/*
 * POWER 4 (CONNECT FOUR) GAME - JAVASCRIPT FRONTEND
 * Modifié : redirige vers /result quand la partie est terminée
 */

// ===== CONFIGURATION DU JEU =====
const COLS = 7;  // Nombre de colonnes dans la grille
const ROWS = 6;  // Nombre de lignes dans la grille

// ===== ÉLÉMENTS DOM =====
const grid = document.getElementById("grid");
const playerDisplay = document.querySelector("#current-player span");

// ===== ÉTAT DU JEU =====
let gameData = null;  // État du jeu reçu du backend
let isWaiting = false; // Empêcher les clics multiples pendant les requêtes

async function initGame() {
  createGrid();
  await fetchGameState();
  updateDisplay();
}

function createGrid() {
  grid.innerHTML = '';
  for (let i = 0; i < ROWS * COLS; i++) {
    const cell = document.createElement("div");
    cell.className = "cell";
    const col = i % COLS;
    cell.addEventListener("click", () => handleColumnClick(col));
    grid.appendChild(cell);
  }
}

async function handleColumnClick(col) {
  if (isWaiting || (gameData && gameData.gameOver)) {
    return;
  }

  isWaiting = true;

  try {
    const response = await fetch('/api/make-move', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ column: col })
    });

    const moveResult = await response.json();

    if (moveResult.success) {
      gameData = moveResult.game;
      updateDisplay();

      // Si le jeu est terminé, rediriger vers la page de résultat
      if (gameData.gameOver) {
        // petit délai pour voir le dernier jeton posé
        setTimeout(() => {
          window.location.href = '/result';
        }, 300);
      }
    } else {
      console.warn("Coup invalide:", moveResult.message);
    }
  } catch (error) {
    console.error("Erreur lors de l'envoi du coup:", error);
  } finally {
    isWaiting = false;
  }
}

async function fetchGameState() {
  try {
    const response = await fetch('/api/game-state');
    gameData = await response.json();
  } catch (error) {
    console.error("Erreur lors de la récupération de l'état du jeu:", error);
  }
}

function updateDisplay() {
  if (!gameData) return;

  const cells = grid.children;
  for (let row = 0; row < ROWS; row++) {
    for (let col = 0; col < COLS; col++) {
      const cellIndex = row * COLS + col;
      const cell = cells[cellIndex];
      const cellValue = gameData.board[row][col];

      cell.classList.remove("red", "yellow");

      if (cellValue !== "empty") {
        cell.classList.add(cellValue);
      }
    }
  }

  if (!gameData.gameOver) {
    const playerText = gameData.turn === "red" ? "Rouge" : "Jaune";
    playerDisplay.textContent = playerText;
    playerDisplay.className = gameData.turn;
  }
}

// suppression de la boîte confirm(); on redirige vers /result à la fin du coup
function showGameOverMessage(message) {
  // plus utilisé : conservé pour debug si besoin
  console.log("Fin de partie:", message);
}

async function resetGame() {
  try {
    const response = await fetch('/api/reset-game', {
      method: 'POST'
    });

    gameData = await response.json();
    updateDisplay();
  } catch (error) {
    console.error("Erreur lors de la réinitialisation du jeu:", error);
  }
}

document.addEventListener("DOMContentLoaded", initGame);

function debugGameState() {
  console.log("État du jeu:", gameData);
}

function debugReset() {
  resetGame();
}
