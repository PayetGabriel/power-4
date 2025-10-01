/*
 * POWER 4 (CONNECT FOUR) GAME - JAVASCRIPT LOGIC
 * 
 * This file contains all the game logic for a Connect Four game.
 * The game is played on a 7x6 grid where players drop tokens from the top.
 * Players alternate between red and yellow tokens.
 * The goal is to connect 4 tokens in a row (horizontally, vertically, or diagonally).
 */

// ===== GAME CONFIGURATION =====
const COLS = 7;  // Number of columns in the game grid (standard Connect Four)
const ROWS = 6;  // Number of rows in the game grid (standard Connect Four)

// ===== GAME STATE =====
let currentPlayer = "red";  // Track whose turn it is ("red" or "yellow")

// ===== DOM ELEMENTS =====
// Get references to HTML elements we'll need to manipulate
const grid = document.getElementById("grid");                    // The game board container
const playerDisplay = document.querySelector("#current-player span");  // Text showing current player

/**
 * CREATE THE GAME GRID
 * 
 * This function dynamically creates all 42 cells (7 columns × 6 rows)
 * and adds click event listeners to each cell.
 * 
 * How it works:
 * - Creates 42 div elements (one for each cell)
 * - Each cell gets a click listener that determines which column was clicked
 * - The modulo operator (%) converts the cell index to a column number
 *   Example: cell 8 → 8 % 7 = column 1
 */
function createGrid() {
  // Loop through all 42 cells (ROWS * COLS = 6 * 7 = 42)
  for (let i = 0; i < ROWS * COLS; i++) {
    const cell = document.createElement("div");  // Create a new div element
    cell.className = "cell";                    // Add CSS class for styling
    
    // Add click listener - when clicked, call handleClick with the column number
    // i % COLS converts the cell index to column number (0-6)
    cell.addEventListener("click", () => handleClick(i % COLS));
    
    grid.appendChild(cell);  // Add the cell to the grid container
  }
}

/**
 * HANDLE COLUMN CLICK
 * 
 * When a player clicks on any cell in a column, this function:
 * 1. Finds the lowest empty cell in that column
 * 2. Drops the current player's token there
 * 3. Switches to the next player
 * 
 * @param {number} col - The column number (0-6) where the player clicked
 */
function handleClick(col) {
  const cells = grid.children;  // Get all the cell elements
  
  // Start from the bottom row and work upward to find an empty cell
  // This simulates gravity - tokens fall to the lowest available position
  for (let row = ROWS - 1; row >= 0; row--) {
    // Convert row/column coordinates to array index
    // Formula: index = row * COLS + col
    // Example: row 5, col 3 → 5 * 7 + 3 = 38
    const index = row * COLS + col;
    const cell = cells[index];
    
    // Check if this cell is empty (doesn't have red or yellow class)
    if (!cell.classList.contains("red") && !cell.classList.contains("yellow")) {
      // Drop the token here
      cell.classList.add(currentPlayer);  // Add color class ("red" or "yellow")
      switchPlayer();                     // Switch to the other player
      break;                             // Stop looking - we found our spot
    }
  }
}

/**
 * SWITCH TO THE NEXT PLAYER
 * 
 * This function:
 * 1. Changes the current player from red to yellow (or vice versa)
 * 2. Updates the display text to show whose turn it is
 * 3. Updates the CSS class for proper text coloring
 */
function switchPlayer() {
  // Toggle between "red" and "yellow" using ternary operator
  currentPlayer = currentPlayer === "red" ? "yellow" : "red";
  
  // Update the display text ("Rouge" = Red in French, "Jaune" = Yellow in French)
  playerDisplay.textContent = currentPlayer === "red" ? "Rouge" : "Jaune";
  
  // Update the CSS class so the text color matches the current player
  playerDisplay.className = currentPlayer;
}

// ===== GAME INITIALIZATION =====
// Start the game by creating the grid when the script loads
createGrid();
