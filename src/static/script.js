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
