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

## Structure du projet

```
TN5/
├── go.mod                          # Module Go
├── wordcount.go                    # Code principal + CLI
├── wordcount_test.go               # 6 tests unitaires + benchmarks
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
go run wordcount.go input.txt
```

Avec une taille de segment personnalisée :

```bash
go run wordcount.go input.txt 5000
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

## Benchmarks disponibles

9 sous-benchmarks par taille de segment (10 à tout-en-un) + 4 comparaisons séquentiel vs concurrent.

| Résultat clé | Valeur |
|--------------|--------|
| Speedup optimal | ~3× (segment 5 000 chars, ~280 goroutines) |
| Point de dégradation | < 100 chars (trop de goroutines) |
| Allocations au point optimal | 279 |

## Liens

- [Rapport TN5](TN5-report.md)
- [Prompts IA](TN5-AI-Prompts.md)
- [Dépôt GitHub](https://github.com/moyamelissa/Advanced-Programming/tree/main/TN5)
