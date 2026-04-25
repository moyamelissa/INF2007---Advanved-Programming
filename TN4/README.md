# Projet : Somme des sinus (INF2007 – Travail 4)

![Tests](https://github.com/moyamelissa/Advanced-Programming/actions/workflows/tn4-coverage.yml/badge.svg)
[![codecov](https://codecov.io/gh/moyamelissa/Advanced-Programming/branch/main/graph/badge.svg?flag=tn4)](https://codecov.io/gh/moyamelissa/Advanced-Programming)

Ce projet mesure et compare les performances du calcul de la somme des sinus sur un tableau de 1 000 000 d'éléments, en entiers et en flottants. Il met en pratique le benchmarking avec `testing.B` et le parsing de flags en ligne de commande.

## Fonctions implémentées

| Fonction | Description |
|----------|-------------|
| `generateIntArray(n)` | Génère un tableau de `n` entiers aléatoires avec seed 42 |
| `generateFloatArray(n)` | Génère un tableau de `n` flottants aléatoires avec seed 42 |
| `computeSineSumInt(data)` | Somme des sinus avec conversion `float64(v)` à chaque itération |
| `computeSineSumFloat(data)` | Somme des sinus sans conversion de type |
| `computeSineSum(dataType, data)` | Dispatch dynamique via `interface{}` et assertion de type |

## Structure du projet

```
TN4/
├── go.mod                       # Module Go
├── sinesum.go                   # Code principal + CLI (flag --type)
├── sinesum_test.go              # 4 tests unitaires + 22 benchmarks
├── TN4-report.md                # Rapport d'analyse avec graphique
├── TN4-results.md               # Résultats détaillés et captures
├── Resultats/                   # Screenshots des sorties terminal
│   ├── Tests unitaires.PNG
│   ├── Benchmarks complets.PNG
│   └── Couverture de code.PNG
├── TN4-AI-Prompts.md            # Prompts IA utilisés
├── TN4-Homework-Instructions.md # Énoncé du travail
└── README.md                    # Ce fichier
```

## Prérequis

- Go 1.21+

## Exécution

```bash
go run sinesum.go --type=float
```

ou avec des entiers :

```bash
go run sinesum.go --type=int
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
| `TestComputeSineSumInt` | Correction du calcul pour `[1, 2, 3]` avec tolérance 1e-9 |
| `TestComputeSineSumFloat` | Correction du calcul pour `[0.1, 0.2, 0.3]` sans conversion |
| `TestComputeSineSumInvalidType` | Rejet d'un type non supporté (`"complex"`) |
| `TestComputeSineSumEmpty` | Somme = 0 pour un tableau vide (Int et Float) |

## Benchmarks disponibles

22 sous-benchmarks au total (11 paliers × 2 types) mesurant le temps de calcul de 1 % à 100 % du tableau de 1 000 000 éléments.

| Résultat clé | Valeur |
|--------------|--------|
| Ratio moyen Int/Float | 1.73× |
| Temps par sinus (Int) | ~36 ns |
| Temps par sinus (Float) | ~21 ns |
| Allocations mémoire | 0 B/op |

## Liens

- [Rapport TN4](TN4-report.md)
- [Résultats détaillés](TN4-results.md)
- [Prompts IA](TN4-AI-Prompts.md)
- [Dépôt GitHub](https://github.com/moyamelissa/Advanced-Programming/tree/main/TN4)
