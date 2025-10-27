/*
 * POWER 4 (CONNECT FOUR) GAME - JAVASCRIPT FRONTEND
 * Version : grille parfaitement espacée et animation ajustée
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

// Retourne la taille (en px) d'une case en fonction du breakpoint.
// On garde les mêmes valeurs que ton CSS : 70px (desktop) / 45px (mobile).
function getCellSize() {
  return window.innerWidth <= 600 ? 45 : 70;
}

async function initGame() {
  await fetchGameState();
  createGrid();
  updateDisplay();
}

// Création de la grille : on positionne explicitement les colonnes à la taille réelle des .cell
function createGrid() {
  if (!gameData || !gameData.board) return;

  const rows = getRows();
  const cols = getCols();
  const cellSize = getCellSize();

  grid.innerHTML = '';
  // on force la grille à utiliser la taille exacte des cellules pour éviter chevauchement
  grid.style.gridTemplateColumns = `repeat(${cols}, ${cellSize}px)`;

  for (let r = 0; r < rows; r++) {
    for (let c = 0; c < cols; c++) {
      const cell = document.createElement("div");
      cell.className = "cell";
      // dataset utiles si besoin futur
      cell.dataset.row = r;
      cell.dataset.col = c;
      cell.style.width = `${cellSize}px`;  // assurer correspondance avec CSS
      cell.style.height = `${cellSize}px`;
      cell.addEventListener("click", () => handleColumnClick(c));
      grid.appendChild(cell);
    }
  }
}

// Gestion du clic sur une colonne
async function handleColumnClick(col) {
  if (isWaiting || (gameData && gameData.gameOver)) return;
  isWaiting = true;

  try {
    const emptyRow = getEmptyRow(col);

    const response = await fetch('/api/make-move', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ column: col })
    });

    const moveResult = await response.json();

    if (moveResult.success) {
      // joueur qui a joué = celui avant la mise à jour de gameData
      const previousPlayer = gameData.turn; // avant update, turn indique le joueur qui devait jouer (celui qui vient d'agir)
      // mettre à jour le state côté frontend
      gameData = moveResult.game;

      // animer en utilisant previousPlayer (couleur correcte)
      await animateTokenFall(col, emptyRow, previousPlayer);

      updateDisplay();

      if (gameData.gameOver) {
        const params = new URLSearchParams(window.location.search);
        const mode = params.get('mode') || gameData.mode || 'normal';
        setTimeout(() => window.location.href = `/result?mode=${encodeURIComponent(mode)}`, 300);
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

// Animation de chute du jeton — calcule taille/position à partir des vraies dimensions des cellules
async function animateTokenFall(col, targetRow, player) {
  return new Promise(resolve => {
    const token = document.createElement("div");
    token.classList.add("falling-token", player);
    document.body.appendChild(token);

    // On récupère la cellule cible
    const cols = getCols();
    const cellIndex = targetRow * cols + col;
    const targetCell = grid.children[cellIndex];
    if (!targetCell) {
      // sécurité : si targetCell manquant, on nettoie et on résoud
      token.remove();
      resolve();
      return;
    }

    const cellRect = targetCell.getBoundingClientRect();
    const gridRect = grid.getBoundingClientRect();

    // taille du token = taille de la cellule (gère responsive)
    const width = Math.round(cellRect.width);
    const height = Math.round(cellRect.height);
    token.style.width = `${width}px`;
    token.style.height = `${height}px`;
    token.style.borderRadius = `${Math.min(width, height) / 2}px`;

    // Position initiale — tenir compte du scroll
    const startLeft = cellRect.left + window.scrollX;
    const startTop = gridRect.top + window.scrollY - cellRect.height;
    token.style.left = `${startLeft}px`;
    token.style.top = `${startTop}px`;

    // Forcer repaint puis animer
    requestAnimationFrame(() => {
      const endTop = cellRect.top + window.scrollY;
      const translateY = endTop - startTop;
      token.style.transform = `translateY(${translateY}px)`;
    });

    token.addEventListener("transitionend", () => {
      token.remove();
      resolve();
    }, { once: true });
  });
}

// Récupère l'état du jeu depuis le backend
async function fetchGameState() {
  try {
    const response = await fetch('/api/game-state');
    gameData = await response.json();
  } catch (error) {
    console.error("Erreur lors de la récupération de l'état du jeu:", error);
  }
}

// Met à jour l'affichage de la grille et du joueur
function updateDisplay() {
  if (!gameData) return;
  const cells = grid.children;

  for (let row = 0; row < getRows(); row++) {
    for (let col = 0; col < getCols(); col++) {
      const idx = row * getCols() + col;
      const cell = cells[idx];
      if (!cell) continue;
      const value = gameData.board[row][col];
      cell.classList.remove("red", "yellow");
      if (value !== "empty") cell.classList.add(value);
    }
  }

  if (!gameData.gameOver) {
    const playerText = gameData.turn === "red" ? "Rouge" : "Jaune";
    playerDisplay.textContent = playerText;
    playerDisplay.className = gameData.turn;
  }
}

// Réinitialisation : conserver le ?mode=... dans l'URL pour demander le bon plateau
async function resetGame() {
  try {
    const params = new URLSearchParams(window.location.search);
    const mode = params.get('mode') || (gameData && gameData.mode) || 'normal';
    const response = await fetch(`/api/reset-game?mode=${encodeURIComponent(mode)}`, { method: 'POST' });
    gameData = await response.json();
    createGrid();
    updateDisplay();
  } catch (error) {
    console.error(error);
  }
}


document.addEventListener("DOMContentLoaded", initGame);
