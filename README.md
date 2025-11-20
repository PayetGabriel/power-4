# power-4

## Présentation

power-4 est une petite application web écrite en Go. Elle sert une interface basée sur des templates HTML et utilise des ressources statiques (CSS/JS) pour l'affichage.

## Sommaire

- Présentation
- Fonctionnalités
- Technologies
- Structure du dépôt
- Prérequis
- Installation et exécution
- Développement
- Organisation du code
- Remarques et dépannage
- FAQ
- Prochaines améliorations
- Auteurs & contact

## Fonctionnalités principales

- Serveur HTTP en Go
- Templates HTML dans `templates/html/`
- Assets statiques (CSS, JS) dans `assets/` et `src/js/`
- Pages : menu, dessin, difficulté, connexion/inscription, écran de victoire

## Technologies

- Go (module `go.mod`)
- Templates Go (`html/template`)
- HTML, CSS, JavaScript côté client

## Structure du dépôt

- `go.mod` : modules Go
- `main.go` : point d'entrée
- `assets/` : styles et autres ressources statiques
- `src/go/` : code Go additionnel
	- `db/db.go` : accès DB (si présent)
	- `game/game.go` : logique de jeu
- `src/js/` : scripts JavaScript
- `templates/html/` : pages HTML

## Prérequis

- Go (version compatible avec `go.mod`)
- Un terminal (exemples fournis pour PowerShell)

## Installation et exécution

1) Cloner le dépôt :

```powershell
git clone https://github.com/PayetGabriel/power-4.git
cd power-4
```

2) Récupérer les dépendances Go :

```powershell
go mod download
```

3) Compiler et lancer :

```powershell
go build ./...
go run main.go

```

Ouvrez ensuite votre navigateur sur l'adresse indiquée par le serveur (habituellement http://localhost:8080) ou appuyer sur ctrl + Click sur l'Url.

## Développement

- Modifiez les fichiers Go dans `main.go` ou sous `src/go/`.
- Les templates sont dans `templates/html/`.
- Les fichiers CSS/JS sont dans `assets/` et `src/js/`.
- Pour recharger automatiquement le serveur lors d'un développement actif, utilisez des outils comme `air` ou `reflex`.

Exemple rapide :

```powershell
go install github.com/cosmtrek/air@latest
air
```

## Organisation du code

- `main.go` : configure les routes et démarre le serveur
- `src/go/db/db.go` : code lié à la base de données
- `src/go/game/game.go` : règles et logique du jeu

Commencez par lire `main.go` pour voir les routes et handlers exposés.

## Remarques et dépannage

- Si le port est déjà utilisé : changez le port dans `main.go` ou libérez le port.
- Erreurs de compilation : lancez `go vet` et `go build` pour repérer les problèmes.
- Ressources statiques manquantes : vérifiez les chemins dans les templates et que le serveur sert bien `assets/`.



## Prochaines améliorations possibles

- Ajouter des tests unitaires pour la logique du jeu
- Mettre en place une CI pour build + tests
- Fournir un Dockerfile si nécessaire

## Auteurs & contact

- Voir l'historique Git pour la liste des contributeurs.

---


