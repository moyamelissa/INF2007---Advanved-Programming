# Projet : Word Stats (INF2007 – Travail 2)

Ce projet est une simulation d’un **workflow Git collaboratif**.  
Il contient trois fonctions Go permettant de compter :
- le nombre de lignes (`countLines`)
- le nombre de mots (`countWords`)
- le nombre de caractères (`countChars`)

## Workflow Git utilisé
- Initialisation du dépôt (`git init`)
- Création de la branche `count-words` + commit
- Création de la branche `count-chars` + commit
- Fusion sans conflit (`count-words` → `main`)
- Fusion avec conflit (`count-chars` → `main`) + résolution manuelle

## Exécution
```bash
go run main.go
```

## Tests
```bash
go test -cover
```
