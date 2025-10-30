/*
 * POWER 4 (CONNECT FOUR) GAME - JAVASCRIPT FRONTEND
 * Version : grille parfaitement espacée et animation ajustée
 */

const grid = document.getElementById("grid"); // sélectionne l'élément de la grille
const playerDisplay = document.querySelector("#current-player span"); // affichage du joueur actuel

let gameData = null; // stocke l'état du jeu côté frontend
let isWaiting = false; // flag pour empêcher double clic pendant animation

function getRows() {
  return gameData.board.length; // retourne le nombre de lignes du plateau
}

function getCols() {
  return gameData.board[0].length; // retourne le nombre de colonnes du plateau
}

// Retourne la taille (en px) d'une case selon breakpoint (responsive)
function getCellSize() {
  return window.innerWidth <= 600 ? 45 : 70; // 45px si mobile, sinon 70px
}

// Initialise le jeu : récupère l'état, crée la grille et met à jour l'affichage
async function initGame() {
  await fetchGameState(); // récupère l'état du jeu depuis le backend
  createGrid();            // construit la grille HTML
  updateDisplay();         // met à jour couleurs et joueur actuel
}

// Création de la grille : on positionne explicitement les colonnes
function createGrid() {
  if (!gameData || !gameData.board) return;

  const rows = getRows();       // nb lignes
  const cols = getCols();       // nb colonnes
  const cellSize = getCellSize(); // taille cellule

  grid.innerHTML = ''; // vide la grille avant reconstruction
  grid.style.gridTemplateColumns = `repeat(${cols}, ${cellSize}px)`; // définit colonnes à taille fixe

  for (let r = 0; r < rows; r++) { // parcours lignes
    for (let c = 0; c < cols; c++) { // parcours colonnes
      const cell = document.createElement("div"); // crée div pour cellule
      cell.className = "cell";                    // classe CSS
      cell.dataset.row = r;                        // stocke ligne
      cell.dataset.col = c;                        // stocke colonne
      cell.style.width = `${cellSize}px`;         // taille width
      cell.style.height = `${cellSize}px`;        // taille height
      cell.addEventListener("click", () => handleColumnClick(c)); // click colonne
      grid.appendChild(cell);                     // ajoute à la grille
    }
  }
}

// Gestion du clic sur une colonne
async function handleColumnClick(col) {
  if (isWaiting || (gameData && gameData.gameOver)) return; // ignore si animation ou partie terminée
  isWaiting = true; // bloque autres clics

  try {
    const emptyRow = getEmptyRow(col); // trouve la ligne vide la plus basse

    const response = await fetch('/api/make-move', { // envoie coup au backend
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ column: col }) // envoie colonne choisie
    });

    const moveResult = await response.json(); // récupère la réponse JSON

    if (moveResult.success) { // si coup valide
      const previousPlayer = gameData.turn; // joueur avant mise à jour
      gameData = moveResult.game;           // met à jour état du jeu

      await animateTokenFall(col, emptyRow, previousPlayer); // anime chute du jeton

      updateDisplay(); // met à jour couleurs et joueur

      if (gameData.gameOver) {
        setTimeout(() => {
          if (gameData.winner === "draw") {
            showModal("Égalité !", "Personne n'a gagné cette partie");
          } else {
            const playerName = gameData.winner === "red" ? "Rouge" : "Jaune";
            showModal("Victoire !", `Le joueur ${playerName} a gagné`);
          }
        }, 300);
      }
    } else {
      console.warn("Coup invalide:", moveResult.message); // colonne pleine ou erreur
    }
  } catch (error) {
    console.error("Erreur lors du coup:", error); // problème fetch
  } finally {
    isWaiting = false; // autorise à nouveau clic
  }
}

// Renvoie la première ligne vide d’une colonne
function getEmptyRow(col) {
  for (let row = getRows() - 1; row >= 0; row--) { // de bas en haut
    if (gameData.board[row][col] === "empty") return row; // retourne ligne vide
  }
  return null; // colonne pleine
}

// Animation de chute du jeton
async function animateTokenFall(col, targetRow, player) {
  return new Promise(resolve => {
    const token = document.createElement("div"); // crée div pour jeton animé
    token.classList.add("falling-token", player); // classe animation + couleur
    document.body.appendChild(token);             // ajoute au DOM

    const cols = getCols();
    const cellIndex = targetRow * cols + col; // index de la cellule cible
    const targetCell = grid.children[cellIndex];
    if (!targetCell) { // sécurité si cellule introuvable
      token.remove();
      resolve();
      return;
    }

    const cellRect = targetCell.getBoundingClientRect(); // position & taille cellule
    const gridRect = grid.getBoundingClientRect();       // position grille

    const width = Math.round(cellRect.width);  // largeur jeton
    const height = Math.round(cellRect.height); // hauteur jeton
    token.style.width = `${width}px`;
    token.style.height = `${height}px`;
    token.style.borderRadius = `${Math.min(width, height) / 2}px`; // cercle parfait

    const startLeft = cellRect.left + window.scrollX; // position gauche initiale
    const startTop = gridRect.top + window.scrollY - cellRect.height; // position top initiale
    token.style.left = `${startLeft}px`;
    token.style.top = `${startTop}px`;

    requestAnimationFrame(() => { // démarre transition
      const endTop = cellRect.top + window.scrollY;
      const translateY = endTop - startTop;
      token.style.transform = `translateY(${translateY}px)`; // animation chute
    });

    token.addEventListener("transitionend", () => { // fin animation
      token.remove(); // supprime jeton temporaire
      resolve();      // résout la promesse
    }, { once: true });
  });
}

// Récupère l'état du jeu depuis le backend
async function fetchGameState() {
  try {
    const response = await fetch('/api/game-state'); // requête GET
    gameData = await response.json();                 // met à jour gameData
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
      const idx = row * getCols() + col; // index cellule dans le DOM
      const cell = cells[idx];
      if (!cell) continue;
      const value = gameData.board[row][col]; // état cellule
      cell.classList.remove("red", "yellow"); // reset couleur
      if (value !== "empty") cell.classList.add(value); // applique couleur si jeton
    }
  }

  if (!gameData.gameOver) { // si partie en cours
    const playerText = gameData.turn === "red" ? "Rouge" : "Jaune";
    playerDisplay.textContent = playerText;  // affiche joueur actuel
    playerDisplay.className = gameData.turn; // applique classe couleur
  }
}

// Réinitialisation : conserve le mode actuel
async function resetGame() {
  try {
    const params = new URLSearchParams(window.location.search);
    const mode = params.get('mode') || (gameData && gameData.mode) || 'normal';
    const response = await fetch(`/api/reset-game?mode=${encodeURIComponent(mode)}`, { method: 'POST' }); // POST reset
    gameData = await response.json(); // met à jour état
    createGrid();                     // recrée la grille
    updateDisplay();                  // met à jour affichage
  } catch (error) {
    console.error(error);
  }
}

document.addEventListener("DOMContentLoaded", initGame); // lance initGame au chargement page

// ==========================
// MODAL VICTOIRE / ÉGALITÉ
// ==========================

// Affiche le modal avec le titre et message
function showModal(title, message) {
  document.getElementById("modal-title").textContent = title;
  document.getElementById("modal-message").textContent = message;
  document.getElementById("modal").classList.remove("hidden");
}

document.addEventListener("DOMContentLoaded", () => {
  const modalClose = document.getElementById("modal-close");
  const modalReplay = document.getElementById("modal-replay");
  const modalMenu = document.getElementById("modal-menu");

  if (modalClose) {
    modalClose.addEventListener("click", () => {
      document.getElementById("modal").classList.add("hidden");
    });
  }

  if (modalReplay) {
    modalReplay.addEventListener("click", async () => {
      document.getElementById("modal").classList.add("hidden");
      await resetGame(); // conserve le mode actuel
    });
  }

  if (modalMenu) {
    modalMenu.addEventListener("click", () => {
      window.location.href = '/menu';
    });
  }
});
