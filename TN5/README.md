# Projet : Comptage de mots concurrent (INF2007 – Travail 5)

![Tests](https://github.com/moyamelissa/Advanced-Programming/actions/workflows/tn5-coverage.yml/badge.svg)
[![codecov](https://codecov.io/gh/moyamelissa/Advanced-Programming/branch/main/graph/badge.svg?flag=tn5)](https://codecov.io/gh/moyamelissa/Advanced-Programming)

Ce projet compte les mots d'un fichier texte en utilisant des goroutines et des canaux (modèle fan-out / fan-in). Il met en pratique la concurrence du Chapitre 8 en Go.

## Fonctions implémentées

| Fonction | Description |
|----------|-------------|
| `countWords(text)` | Compte les mots dans une chaîne via `strings.Fields` |
| `splitIntoSegments(content, segmentSize)` | Découpe le texte en segments sans couper les mots |
| `countWordsInSegment(segment, ch)` | Goroutine qui envoie le nombre de mots sur un canal |
| `CountWordsConcurrent(content, segmentSize)` | Orchestre le fan-out/fan-in et somme les résultats |
| `run(args)` | Logique principale testable, extraite de `main` |
| `exitFunc` | Variable injectable pour tester la branche d'erreur de `main` |
## Structure du projet

```
.
├── go.mod                          # Module Go (wordcount)
├── wordcount.go                    # Code principal + CLI
├── wordcount_test.go               # 14 tests unitaires + 13 benchmarks
├── input.txt                       # Fichier texte de test
├── TN5-report.md                   # Rapport d'analyse
├── TN5-AI-Prompts.md               # Prompts IA utilisés
├── TN5-Homework-Instructions.md    # Énoncé du travail
└── README.md                       # Ce fichier
```

## Prérequis

- Go 1.21+

## Exécution

```bash
go run . input.txt
```

Avec une taille de segment personnalisée :

```bash
go run . input.txt 5000
```

## Tests unitaires

```bash
go test -v -run="Test" ./...
```

## Benchmarks

```bash
go test -bench="Benchmark" -benchmem -run="^$" -count=1 ./...
```

## Tests disponibles

| Test | Ce qu'il vérifie |
|------|-----------------|
| `TestCountWordsEmpty` | Contenu vide retourne 0 |
| `TestCountWordsSingleWord` | Un seul mot `"Bonjour"` retourne 1 |
| `TestCountWordsMultipleLines` | 3 lignes × 3 mots = 9 avec `\n` comme séparateur |
| `TestCountWordsMultipleSpaces` | Espaces multiples ignorés, 3 mots comptés |
| `TestCountWordsConsistency` | Résultat identique pour 7 tailles de segment (1 à 500) |
| `TestSplitIntoSegments` | Les segments contiennent des mots complets |
| `TestSplitIntoSegmentsNegativeSize` | Taille ≤ 0 retourne tout le contenu en un seul segment |
| `TestRunValidFile` | `run` avec fichier réel retourne le bon compte |
| `TestRunWithSegmentSize` | `run` accepte une taille de segment en argument |
| `TestRunNoArgs` | `run` sans arguments retourne une erreur |
| `TestRunInvalidSegmentSize` | `run` avec taille invalide (`"abc"`, `"-5"`) retourne une erreur |
| `TestRunMissingFile` | `run` avec fichier inexistant retourne une erreur |
| `TestMainFunction` | `main()` sans panique pour un fichier valide |
| `TestMainFunctionError` | `main()` appelle `exitFunc(1)` en cas d'erreur |
## Benchmarks disponibles

9 sous-benchmarks par taille de segment (10 à tout-en-un) + 4 comparaisons séquentiel vs concurrent (13 au total).

| Résultat clé | Valeur |
|--------------|--------|
| Speedup optimal | ~3× (segment 5 000 chars, ~280 goroutines) |
| Point de dégradation | < 100 chars (trop de goroutines) |
| Allocations au point optimal | 279 |

## Liens

- [Rapport TN5](TN5-report.md)
- [Prompts IA](TN5-AI-Prompts.md)
- [Dépôt GitHub](https://github.com/moyamelissa/Advanced-Programming/tree/main/TN5)
