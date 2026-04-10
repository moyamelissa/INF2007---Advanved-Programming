# Projet : Word Stats (INF2007 – Travail 2)

Ce projet est une simulation d'un **workflow Git collaboratif** dans le cadre du cours INF2007.
Il s'agit d'une application Go qui génère des statistiques sur un texte donné.

## Fonctions implémentées

| Fonction | Description |
|----------|-------------|
| `countLines(text)` | Compte le nombre de lignes (séparées par `\n`) |
| `countWords(text)` | Compte le nombre de mots (séparés par des espaces/whitespace) |
| `countChars(text)` | Compte le nombre de caractères (hors espaces, tabs et sauts de ligne) |

## Structure du projet

```
word-stats/
├── go.mod          # Module Go
├── main.go         # Code source principal (3 fonctions + main)
├── main_test.go    # Tests unitaires (9 tests, couverture 100 %)
├── history.txt     # Historique des commits (git log --oneline)
└── README.md       # Ce fichier
```

## Prérequis

- Go 1.16+
- Git

## Exécution

```bash
go run main.go
```

Résultat attendu :
```
Nombre de mots : 3
Nombre de caractères : 16
```

## Tests unitaires

```bash
go test -v
```

Avec couverture :
```bash
go test -cover
```

Résultat :
```
PASS
coverage: 100.0% of statements
```

## Workflow Git

### Branches utilisées

| Branche | Objectif |
|---------|----------|
| `main` | Branche principale, contient le code initial (`countLines`) |
| `count-words` | Ajout de la fonction `countWords` |
| `count-chars` | Ajout de la fonction `countChars` |

### Étapes du workflow

1. **Initialisation** — Création du dépôt, ajout de `countLines`, premier commit sur `main`
2. **Branche `count-words`** — Ajout de `countWords`, commit
3. **Branche `count-chars`** — Retour sur `main`, création de la branche, ajout de `countChars`, commit
4. **Fusion sans conflit** — Merge de `count-words` dans `main` (fast-forward)
5. **Fusion avec conflit** — Modification de `main()` dans les deux branches (affichages différents), merge de `count-chars` → conflit dans `main.go`
6. **Résolution du conflit** — Suppression des marqueurs de conflit, combinaison des deux affichages (mots + caractères)

### Commandes Git utilisées

| Commande | Utilisation |
|----------|-------------|
| `git init` | Initialiser le dépôt |
| `git branch -M main` | Renommer la branche principale |
| `git add <fichier>` | Ajouter des fichiers au staging |
| `git commit -m "..."` | Valider les changements |
| `git branch` | Lister les branches |
| `git checkout -b <nom>` | Créer et basculer sur une branche |
| `git checkout <nom>` | Basculer sur une branche existante |
| `git merge <nom>` | Fusionner une branche dans la branche courante |
| `git log --oneline` | Afficher l'historique des commits |

### Résolution du conflit

Lors du merge de `count-chars` dans `main`, Git a signalé un conflit dans `main.go` car la fonction `main()` avait été modifiée différemment dans les deux branches :
- `main` affichait uniquement les **mots**
- `count-chars` affichait uniquement les **caractères**

La résolution a consisté à combiner les deux affichages dans `main()` et à conserver les deux fonctions (`countWords` et `countChars`).
