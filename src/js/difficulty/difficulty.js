document.querySelectorAll(".difficulty").forEach(btn => {
  btn.addEventListener("click", () => {
    const mode = btn.dataset.mode;
    window.location.href = `/game?mode=${mode}`; // Les backticks `` permettent d'insérer des variables avec ${...}
  });
});
