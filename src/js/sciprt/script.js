/*
 * POWER 4 (CONNECT FOUR) GAME - JAVASCRIPT FRONTEND
 * Support des grilles de différentes tailles selon le mode
 */
console.log("🚀 Script chargé !");

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
  await fetchGameState(); // récupère l'état du jeu depuis le serveur
  createGrid(); // crée visuellement la grille HTML
  updateDisplay(); // affiche les jetons déjà placés (si partie en cours)
}


function createGrid() {
  if (!gameData || !gameData.board) return; // Si pas de données, on quitte

  const rows = getRows();
  const cols = getCols();

  // ✅ Configuration dynamique du CSS Grid avec !important via style
  grid.style.cssText = `
    display: grid !important;
    grid-template-columns: repeat(${cols}, var(--cell-size)) !important;
    grid-template-rows: repeat(${rows}, var(--cell-size)) !important;
    gap: var(--cell-gap);
    background-color: var(--board-bg);
    border-radius: 18px;
    padding: 16px;
    box-shadow: 0 6px 18px var(--shadow-light);
    width: fit-content;
    margin: 0 auto;
  `;

  grid.innerHTML = ''; // Vide la grille avant de la recréer
  for (let r = 0; r < rows; r++) {           // Pour chaque ligne
  for (let c = 0; c < cols; c++) {         // Pour chaque colonne
    const cell = document.createElement("div");  // Crée un <div>
    cell.className = "cell";                     // Ajoute la classe CSS "cell"
    cell.dataset.row = r;                        // data-row="0", "1", etc.
    cell.dataset.col = c;                        // data-col="0", "1", etc.
      cell.addEventListener("click", () => handleColumnClick(c));
      grid.appendChild(cell);
    }
  }

  console.log(`✅ Grille créée: ${rows} lignes × ${cols} colonnes (mode: ${gameData.mode})`);
  console.log(`📐 Style grid appliqué: ${cols} colonnes × ${rows} lignes`);
}


async function handleColumnClick(col) {
  if (isWaiting || (gameData && gameData.gameOver)) return; // Ignore si en attente ou partie terminée
  isWaiting = true;

  try {
    // On récupère la ligne de destination avant d'envoyer le coup
    const emptyRow = getEmptyRow(col);

    if (emptyRow === null) {
      console.warn("Colonne pleine");
      isWaiting = false;
      return;
    }

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

// Renvoie la première ligne vide d'une colonne
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
    
    if (!targetCell) {
      console.error("Cellule cible introuvable");
      token.remove();
      resolve();
      return;
    }

    const cellRect = targetCell.getBoundingClientRect();
    const gridRect = grid.getBoundingClientRect();

    // Position initiale (au-dessus de la grille)
    token.style.left = `${cellRect.left}px`;
    token.style.top = `${gridRect.top - cellRect.height}px`;
    token.style.width = `${cellRect.width * 0.8}px`;
    token.style.height = `${cellRect.height * 0.8}px`;

    // Déclenche l'animation
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
    console.log("📊 État du jeu récupéré:", gameData.mode, `${getRows()}x${getCols()}`);
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
      if (!cell) continue;
      
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
    const response = await fetch('/api/reset-game', { method: 'POST' }); // Envoie une requête POST pour réinitialiser
    gameData = await response.json();
    createGrid();
    updateDisplay();
  } catch (error) {
    console.error(error);
  }
}

document.addEventListener("DOMContentLoaded", initGame);