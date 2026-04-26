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
├── go.mod                          # Module Go
├── sinesum.go                      # Code principal + CLI (flag --type)
├── sinesum_test.go                 # 13 tests unitaires + 22 benchmarks
├── TN4-report.md                   # Rapport d'analyse avec tableau, graphique et calculs
├── TN4-AI-Prompts.md               # Prompts IA utilisés
├── TN4-Homework-Instructions.md    # Énoncé du travail
├── Results-and-Instructions/       # Résultats et guides
│   ├── Resultats-benchmarks-et-captures.md  # Données brutes et tableau complet
│   ├── Guide-applications-numeriques.md      # Comment calculer lumière + 120 fps
│   ├── Guide-creation-graphique-Mermaid.md   # Comment créer le graphique Mermaid
│   ├── bench_count6.txt             # Sortie brute des benchmarks (count=6)
│   ├── benchstat-output.txt         # Analyse benchstat (médianes ± IC 95 %)
│   ├── tests-output.txt             # Sortie de go test -v
│   └── coverage-output.txt          # Sortie de go test -cover (100 %)
└── README.md                       # Ce fichier
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
go test -bench="Benchmark" -benchmem -run="^$" -count=6 ./...
```

## Tests disponibles

| Test | Ce qu'il vérifie |
|------|-----------------|
| `TestComputeSineSumInt` | Correction du calcul pour `[1, 2, 3]` avec tolérance 1e-9 |
| `TestComputeSineSumFloat` | Correction du calcul pour `[0.1, 0.2, 0.3]` sans conversion |
| `TestComputeSineSumInvalidType` | Rejet d'un type non supporté (`"complex"`) |
| `TestComputeSineSumEmpty` | Somme = 0 pour un tableau vide (Int et Float) |
| `TestComputeSineSumNegative` | Entiers négatifs `[-1, 0, 1]`, propriété d'imparité de sin |
| `TestComputeSineSumLargeFloat` | Stabilité numérique avec `1e15` (réduction d'argument) |
| `TestComputeSineSumIntWrongData` | Dispatch "int" avec données `[]float64` retourne erreur |
| `TestComputeSineSumFloatWrongData` | Dispatch "float" avec données `[]int` retourne erreur |
| `TestRunInt` | `run("int")` retourne un résultat sans erreur |
| `TestRunFloat` | `run("float")` retourne un résultat sans erreur |
| `TestRunInvalidType` | `run("complex")` retourne une erreur |
| `TestMainFunction` | `main()` s'exécute sans panique (cas nominal) |
| `TestMainFunctionError` | `main()` gère un type invalide sans panique |

## Benchmarks disponibles

22 sous-benchmarks au total (11 paliers × 2 types) mesurant le temps de calcul de 1 % à 100 % du tableau de 1 000 000 éléments.

| Résultat clé | Valeur |
|--------------|--------|
| Ratio Int/Float (100 %) | 1.85× |
| Temps par sinus (Int) | ~39 ns |
| Temps par sinus (Float) | ~21 ns |
| Allocations mémoire | 0 B/op |

## Liens

- [Rapport TN4](TN4-report.md)
- [Résultats et captures](Results-and-Instructions/Resultats-benchmarks-et-captures.md)
- [Guide calculs](Results-and-Instructions/Guide-applications-numeriques.md)
- [Guide graphique Mermaid](Results-and-Instructions/Guide-creation-graphique-Mermaid.md)
- [Prompts IA](TN4-AI-Prompts.md)
- [Dépôt GitHub](https://github.com/moyamelissa/Advanced-Programming/tree/main/TN4)
