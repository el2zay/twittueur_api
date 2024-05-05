# API Twittueur
>  Ceci est un projet de NSI.

Ce repo contient le code source de l'API de Twittueur.

## Installation

1. Clonez le repo

`git clone https://github.com/el2zay/twittueur_api.git`

2. Créez un dossier "db" à la racine du projet
3. Lancez le serveur avec [Go](https://go.dev/doc/install)

`go build . && ./twittueur_api`

4. Vous pouvez maintenant accéder à l'API à l'adresse http://localhost:1323

## Arborescence
- Le dossier `db` contient les fichiers de la base de données.
- Le dossier `models` contient les modèles de données.
- Le dossier `routes` (le plus intéressant) contient les routes de l'API, avec les fonctions associées.
- Le dossier `src/utils` contient des fonctions utilitaires.
- Le fichier `server.go` contient le code pour lancer le serveur, et d'associer les fonctions du dossier `routes` aux requêtes HTTP.
- Le fichier `words.json` contient un dictioannaire de mots pour générer une passphrase aléatoirement.

