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
| `CountWordsConcurrentN(content, n)` | Variante avec n goroutines explicites pour mesurer la linéarité |
| `CountWordsConcurrentPool(content, n)` | Worker pool avec n goroutines persistantes |
| `run(args)` | Logique principale testable, extraite de `main` |
| `exitFunc` | Variable injectable pour tester la branche d'erreur de `main` |
## Structure du projet

```
.
├── go.mod                          # Module Go (wordcount)
├── wordcount.go                    # Code principal + interface en ligne de commande
├── wordcount_test.go               # 16 tests unitaires + 25 sous-benchmarks
├── input.txt                       # Fichier texte de test
├── data/                           # Graphique et logs des benchmarks
│   ├── worker-count-chart.png       # Graphique 1 – goroutine-par-segment vs linéaire idéal
│   ├── worker-pool-chart.png        # Graphique 2 – worker pool vs goroutine-par-segment
│   ├── bench-raw.txt                # Données brutes des benchmarks
│   ├── workerpool-raw.txt           # Données brutes BenchmarkWorkerPool
│   └── benchstat.txt                # Analyse statistique (benchstat)
├── docs/                           # Documentation technique
│   ├── TN5-Calculs.md               # Calculs détaillés (speedup, Amdahl, allocations)
│   ├── TN5-AI-Prompts.md            # Prompts IA utilisés
│   └── TN5-Homework-Instructions.md # Énoncé du travail
├── chart/                          # Générateur de graphique (module Go séparé)
│   ├── go.mod                       # Module Go (chart)
│   └── main.go                      # Générateur du graphique
├── TN5-report.md                   # Rapport d'analyse
└── README.md                       # Ce fichier
```

## Prérequis

- Go 1.26.1

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

## Bancs d'essai

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
| `TestCountWordsConcurrentN` | `CountWordsConcurrentN` avec 1, 2, 4, 8, 16 workers + cas limites |
| `TestCountWordsConcurrentPool` | `CountWordsConcurrentPool` avec 1, 2, 4, 8, 16 workers + cas limites |

## Bancs d'essai disponibles

25 sous-benchmarks au total, soit 9 tailles de segment par `BenchmarkSegmentSize`,
4 comparaisons séquentiel contre concurrent par `BenchmarkSequentialVsConcurrent`,
6 nombres de goroutines par `BenchmarkWorkerCount` et 6 par `BenchmarkWorkerPool`.

| Résultat clé | Valeur |
|--------------|--------|
| Accélération optimale (taille segments) | 2.85× (segment 50 000 caract.) |
| Point de dégradation | < 1 000 caract. (surcharge d'ordonnancement) |
| Allocations au point optimal | 34 |
| Worker pool vs goroutine-par-segment | 4–5× plus rapide |
| Couverture de code | 100 % |

## Liens

- [Rapport TN5](TN5-report.md)
- [Calculs détaillés](docs/TN5-Calculs.md)
- [Prompts IA](docs/TN5-AI-Prompts.md)
- [Énoncé du travail](docs/TN5-Homework-Instructions.md)
- [Dépôt GitHub](https://github.com/moyamelissa/Advanced-Programming/tree/main/TN5)
