const grid = document.getElementById("grid");
const cells = Array.from(grid.children);
const COLS = 7;
const ROWS = 6;
let turn = "red";

for (let col = 0; col < COLS; col++) {
  for (let row = 0; row < ROWS; row++) {
    let index = row * COLS + col;
    cells[index].addEventListener("click", () => dropToken(col));
  }
}

function dropToken(col) {
  for (let row = ROWS - 1; row >= 0; row--) {
    let index = row * COLS + col;
    if (!cells[index].classList.contains("red") && !cells[index].classList.contains("yellow")) {
      cells[index].classList.add(turn);
      turn = turn === "red" ? "yellow" : "red";
      break;
    }
  }
}
