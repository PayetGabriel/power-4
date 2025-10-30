document.querySelectorAll(".difficulty").forEach(btn => { // sélectionne tous les boutons avec la classe "difficulty" et boucle dessus
  btn.addEventListener("click", () => { // ajoute un événement au clic pour chaque bouton
    const mode = btn.dataset.mode; // récupère le mode depuis l'attribut data-mode du bouton
    window.location.href = `/game?mode=${mode}`; // redirige le navigateur vers /game avec le mode choisi en query param
  });
});
