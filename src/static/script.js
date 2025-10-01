/*
 * POWER 4 (CONNECT FOUR) GAME - JAVASCRIPT FRONTEND
 * 
 * Ce fichier gère l'interface utilisateur et communique avec le backend Go
 * pour la logique du jeu et la détection de victoire.
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

/**
 * INITIALISER LE JEU
 * 
 * Cette fonction démarre le jeu en créant la grille
 * et en récupérant l'état initial du serveur.
 */
async function initGame() {
  createGrid();
  await fetchGameState();
  updateDisplay();
}

/**
 * CRÉER LA GRILLE DE JEU
 * 
 * Crée 42 cellules cliquables (7 colonnes × 6 lignes)
 * avec les gestionnaires d'événements appropriés.
 */
function createGrid() {
  grid.innerHTML = ''; // Vider la grille existante
  
  for (let i = 0; i < ROWS * COLS; i++) {
    const cell = document.createElement("div");
    cell.className = "cell";
    
    // Calculer la colonne à partir de l'index
    const col = i % COLS;
    
    // Ajouter un gestionnaire de clic pour cette colonne
    cell.addEventListener("click", () => handleColumnClick(col));
    
    grid.appendChild(cell);
  }
}

/**
 * GÉRER LE CLIC SUR UNE COLONNE
 * 
 * Envoie le coup au serveur et met à jour l'affichage
 * avec la réponse reçue.
 * 
 * @param {number} col - Numéro de la colonne cliquée (0-6)
 */
async function handleColumnClick(col) {
  // Empêcher les clics multiples
  if (isWaiting || (gameData && gameData.gameOver)) {
    return;
  }

  isWaiting = true;

  try {
    // Envoyer le coup au serveur
    const response = await fetch('/api/make-move', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ column: col })
    });

    const moveResult = await response.json();

    if (moveResult.success) {
      // Mettre à jour l'état du jeu avec la réponse du serveur
      gameData = moveResult.game;
      updateDisplay();

      // Afficher un message si le jeu est terminé
      if (gameData.gameOver) {
        if (moveResult.tokens) {
          drawStars(moveResult.tokens);
        }
        setTimeout(() => {
          showGameOverMessage(moveResult.message);
        }, 300); // Petit délai pour voir le dernier coup
      }
    } else {
      // Afficher l'erreur si le coup n'est pas valide
      console.warn("Coup invalide:", moveResult.message);
    }
  } catch (error) {
    console.error("Erreur lors de l'envoi du coup:", error);
  } finally {
    isWaiting = false;
  }
}

/**
 * RÉCUPÉRER L'ÉTAT DU JEU DEPUIS LE SERVEUR
 * 
 * Synchronise l'état local avec le serveur.
 */
async function fetchGameState() {
  try {
    const response = await fetch('/api/game-state');
    gameData = await response.json();
  } catch (error) {
    console.error("Erreur lors de la récupération de l'état du jeu:", error);
  }
}

/**
 * METTRE À JOUR L'AFFICHAGE
 * 
 * Met à jour la grille visuelle et l'indicateur de joueur
 * basé sur l'état actuel du jeu.
 */
function updateDisplay() {
  if (!gameData) return;

  // Mettre à jour toutes les cellules de la grille
  const cells = grid.children;
  for (let row = 0; row < ROWS; row++) {
    for (let col = 0; col < COLS; col++) {
      const cellIndex = row * COLS + col;
      const cell = cells[cellIndex];
      const cellValue = gameData.board[row][col];

      // Supprimer les anciennes classes de couleur
      cell.classList.remove("red", "yellow");

      // Ajouter la nouvelle classe si la cellule n'est pas vide
      if (cellValue !== "empty") {
        cell.classList.add(cellValue);
      }
    }
  }

  // Mettre à jour l'indicateur de joueur actuel
  if (!gameData.gameOver) {
    const playerText = gameData.turn === "red" ? "Rouge" : "Jaune";
    playerDisplay.textContent = playerText;
    playerDisplay.className = gameData.turn;
  }
}

/**
 * AFFICHER LE MESSAGE DE FIN DE JEU
 * 
 * Montre qui a gagné ou si c'est un match nul,
 * et propose de recommencer.
 * 
 * @param {string} message - Message à afficher
 */
function showGameOverMessage(message) {
  const playAgain = confirm(message + "\n\nVoulez-vous jouer une nouvelle partie ?");
  
  if (playAgain) {
    resetGame();
  }
}

/**
 * RÉINITIALISER LE JEU
 * 
 * Demande au serveur de remettre le jeu à zéro
 * et met à jour l'affichage.
 */
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
