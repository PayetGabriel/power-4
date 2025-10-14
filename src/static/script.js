/*
 * POWER 4 (CONNECT FOUR) GAME - JAVASCRIPT FRONTEND
 * Ajout : animation de chute du jeton
 */

const COLS = 7;
const ROWS = 6;

const grid = document.getElementById("grid");
const playerDisplay = document.querySelector("#current-player span");

let gameData = null;
let isWaiting = false;

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
        if (moveResult.tokens) {
          drawStars(moveResult.tokens);
        }
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
  for (let row = ROWS - 1; row >= 0; row--) {
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
    const cellIndex = targetRow * COLS + col;
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

async function resetGame() {
  try {
    const response = await fetch('/api/reset-game', { method: 'POST' });
    gameData = await response.json();
    updateDisplay();
  } catch (error) {
    console.error("Erreur lors de la réinitialisation du jeu:", error);
  }
}

/**
 * DESSINER DES ÉTOILES BLEUES
 * 
 * Ajoute des étoiles bleues au centre des jetons alignés.
 * 
 * @param {Array} tokens - Coordonnées des jetons alignés (ex: [[0,1], [0,2], [0,3], [0,4]])
 */
function drawStars(tokens) {
  console.log("drawStars called with tokens:", tokens); // Log the function call

  const cells = grid.children;

  tokens.forEach(([row, col]) => {
    console.log(`Processing token at row ${row}, col ${col}`); // Log each token
    const index = row * COLS + col;
    const cell = cells[index];

    if (!cell) {
      console.warn(`No cell found at index ${index} for row ${row}, col ${col}`);
      return;
    }

    // Vérifier si une étoile existe déjà
    const existingStar = cell.querySelector('.star');
    if (!existingStar) {
      // Ajouter une étoile bleue avec des styles temporaires
      const star = document.createElement('div');
      star.className = 'star';
      star.style.backgroundColor = 'red'; // Temporary bright color for visibility
      star.style.border = '2px solid yellow'; // Temporary border for debugging
      cell.appendChild(star);
      console.log("Étoile ajoutée à la cellule :", cell); // Log de confirmation

      // Log computed styles
      const computedStyles = window.getComputedStyle(star);
      console.log("Computed styles for .star:", computedStyles);
    } else {
      console.log("Étoile déjà présente dans la cellule :", cell); // Log if star already exists
    }
  });
}

// Ajouter un style CSS pour les étoiles
const style = document.createElement('style');
style.textContent = `
  .star {
    width: 20px;
    height: 20px;
    background-color: blue;
    clip-path: polygon(50% 0%, 61% 35%, 98% 35%, 68% 57%, 79% 91%, 50% 70%, 21% 91%, 32% 57%, 2% 35%, 39% 35%);
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    z-index: 10; /* S'assurer que l'étoile est au-dessus des jetons */
    border: 2px solid red; /* Test visuel temporaire */
  }
`;
document.head.appendChild(style);

// ===== INITIALISATION =====
// Démarrer le jeu quand la page est chargée
document.addEventListener("DOMContentLoaded", initGame);

// ===== FONCTIONS UTILITAIRES POUR LE DÉBOGAGE =====
// Ces fonctions peuvent être appelées depuis la console du navigateur

/**
 * Afficher l'état actuel du jeu dans la console
 */
function debugGameState() {
  console.log("État du jeu:", gameData);
}

/**
 * Réinitialiser manuellement le jeu (pour le débogage)
 */
function debugReset() {
  resetGame();
}

// Test manuel pour vérifier l'affichage des étoiles
function testStarDisplay() {
  const testCell = grid.children[0]; // Première cellule
  const testStar = document.createElement('div');
  testStar.className = 'star';
  testCell.appendChild(testStar);
  console.log("Étoile ajoutée manuellement à la première cellule.");
}

testStarDisplay();
