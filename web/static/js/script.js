/*
 * POWER 4 (CONNECT FOUR) GAME - JAVASCRIPT FRONTEND
 * Ajout : animation de chute du jeton
 */

const grid = document.getElementById("grid");
const playerDisplay = document.querySelector("#current-player span");

let gameData = null;
let isWaiting = false;


function getRows() {
  return gameData.board.length;
}

function getCols() {
  return gameData.board[0].length;
}

async function initGame() {
  await fetchGameState(); // d'abord récupérer le plateau
  createGrid();           // puis créer la grille selon gameData
  updateDisplay();
}


function createGrid() {
  if (!gameData || !gameData.board) return;

  const rows = gameData.board.length;
  const cols = gameData.board[0].length;

  grid.innerHTML = '';
  for (let r = 0; r < rows; r++) {
    for (let c = 0; c < cols; c++) {
      const cell = document.createElement("div");
      cell.className = "cell";
      cell.addEventListener("click", () => handleColumnClick(c));
      grid.appendChild(cell);
    }
  }
}


async function handleColumnClick(col) {
  if (isWaiting || (gameData && gameData.gameOver)) return;
  isWaiting = true;

  try {
    // On récupère la ligne de destination avant d’envoyer le coup
    const emptyRow = getEmptyRow(col);

    const response = await fetch('/api/make-move', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ column: col })
    });

    const moveResult = await response.json();

    if (moveResult.success) {
      const player = gameData ? gameData.turn : "red";
      await animateTokenFall(col, emptyRow, player);

      gameData = moveResult.game;
      updateDisplay();

      if (gameData.gameOver) {
        setTimeout(() => {
          window.location.href = '/result';
        }, 300);
      }
    } else {
      console.warn("Coup invalide:", moveResult.message);
    }
  } catch (error) {
    console.error("Erreur lors du coup:", error);
  } finally {
    isWaiting = false;
  }
}

// Renvoie la première ligne vide d’une colonne
function getEmptyRow(col) {
  for (let row = getRows() - 1; row >= 0; row--) {
    if (gameData.board[row][col] === "empty") return row;
  }
  return null;
}

// Animation de chute
async function animateTokenFall(col, targetRow, player) {
  return new Promise(resolve => {
    const token = document.createElement("div");
    token.classList.add("falling-token", player);
    document.body.appendChild(token);

    // On récupère la cellule cible pour position exacte
    const cellIndex = targetRow * getCols() + col;
    const targetCell = grid.children[cellIndex];
    const cellRect = targetCell.getBoundingClientRect();

    // Position initiale (au-dessus de la grille)
    const gridRect = grid.getBoundingClientRect();
    token.style.left = `${cellRect.left}px`;
    token.style.top = `${gridRect.top - cellRect.height}px`;

    // Déclenche l’animation
    requestAnimationFrame(() => {
      const endY = cellRect.top;
      token.style.transform = `translateY(${endY - (gridRect.top - cellRect.height)}px)`;
    });

    token.addEventListener("transitionend", () => {
      token.remove();
      resolve();
    }, { once: true });
  });
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
  for (let row = 0; row < getRows(); row++) {
    for (let col = 0; col < getCols(); col++) {
      const cellIndex = row * getCols() + col;
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

async function resetGame() {
  try {
    const response = await fetch('/api/reset-game' + window.location.search, { method: 'POST' });
    gameData = await response.json();
    createGrid();
    updateDisplay();
  } catch (error) {
    console.error(error);
  }
}

document.addEventListener("DOMContentLoaded", initGame);