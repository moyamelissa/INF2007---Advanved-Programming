# Somme des sinus — INF2007 Travail 4

![Tests](https://github.com/moyamelissa/Advanced-Programming/actions/workflows/tn4-coverage.yml/badge.svg)
[![codecov](https://codecov.io/gh/moyamelissa/Advanced-Programming/branch/main/graph/badge.svg?flag=tn4)](https://codecov.io/gh/moyamelissa/Advanced-Programming)

Mesure et comparaison des performances du calcul de la somme des sinus sur un tableau de 1 000 000 d'éléments, en entiers et en nombres à virgule flottante. Le projet met en pratique le benchmarking avec `testing.B`, le parsing de flags CLI et la couverture de code à 100 % en Go.

## Fonctions implémentées

| Fonction | Description |
|----------|-------------|
| `generateIntArray(n)` | Génère `n` entiers aléatoires dans \[0, 1000\] avec graine 42 |
| `generateFloatArray(n)` | Génère `n` flottants aléatoires dans \[0, 1) avec graine 42 |
| `computeSineSumInt(data)` | Somme des sinus d'un `[]int` avec conversion `float64(v)` à chaque itération |
| `computeSineSumFloat(data)` | Somme des sinus d'un `[]float64` sans conversion de type |
| `computeSineSum(dataType, data)` | Dispatch dynamique via `interface{}` et assertion de type |
| `run(dataType)` | Logique principale testable, extraite de `main` |

## Structure du projet

Le programme calcule **Σ sin(xᵢ)** pour i = 1..1 000 000 éléments générés aléatoirement (graine 42). Le résultat numérique est proche de zéro (les valeurs aléatoires couvrent plusieurs cycles de 2π, les contributions positives et négatives s'annulent). L'intérêt du projet est la **performance** : combien de temps prend ce calcul selon que les données sont des `int` ou des `float64` ?

```
.
├── go.mod                          # Module Go (sinesum)
├── sinesum.go                      # Code principal + CLI (flag --type)
├── sinesum_test.go                 # 13 tests unitaires + 22 benchmarks — couverture 100 %
├── TN4-report.md                   # Rapport d'analyse (tableau, graphique, calculs)
├── TN4-AI-Prompts.md               # Prompts IA utilisés
├── TN4-Homework-Instructions.md    # Énoncé du travail
├── chart/                          # Générateur de graphique (module Go séparé)
│   ├── go.mod                      # Module chart — dépendance gonum/plot isolée
│   └── main.go                     # Produit docs/benchmark-chart.png
├── docs/                           # Guides et visuels de référence
│   ├── benchmark-chart.png              # Graphique Int vs Float (généré par chart/)
│   ├── benchmark-results.md            # Données brutes et tableau complet
│   ├── calculs.md                       # Calculs détaillés Q1 (lumière) et Q2 (120 fps)
│   └── chart-guide.md                   # Comment régénérer le graphique avec gonum/plot
└── logs/                           # Sorties brutes des commandes Go
    ├── bench-raw.txt                # go test -bench … -count=6 (sortie brute)
    ├── benchstat.txt                # benchstat — médianes ± IC 95 %
    ├── tests.txt                    # go test -v (13 tests PASS)
    └── coverage.txt                 # go test -cover → 100.0 %
```

## Prérequis

- Go 1.21+

## Démarrage rapide

```bash
# Depuis le dossier TN4/
go run . --type=float   # tableau de nombres à virgule flottante
go run . --type=int     # tableau d'entiers
```

Exemple de sortie (`--type=float`) :

```
=== Somme des sinus — type=float, taille=1000000 ===

Génération du tableau : 14.2ms
Calcul de la somme    : 20.9ms

Résultat              : -103.655941
Temps total           : 35.1ms
```

Le **résultat** est la somme des sin(x) pour x ∈ [0, 1) — proche de zéro car sin oscille symétriquement. Pour les entiers (`--type=int`), les valeurs dans [0, 1000] couvrent ~159 cycles complets de 2π, ce qui produit également une somme proche de zéro mais avec un calcul ~1,85× plus lent.

## Tests

**Tests unitaires (13 tests)**

```bash
go test -v -run="Test" ./...
```

**Couverture de code**

```bash
go test -cover ./...
# output : coverage: 100.0% of statements
```

**Benchmarks (22 sous-benchmarks, 11 paliers × 2 types)**

```bash
go test -bench="Benchmark" -benchmem -run="^$" -count=6 ./...
```

## Tests disponibles

| Test | Ce qu'il vérifie |
|------|-----------------|
| `TestComputeSineSumInt` | Correction du calcul pour `[1, 2, 3]`, tolérance 1e-9 |
| `TestComputeSineSumFloat` | Correction du calcul pour `[0.1, 0.2, 0.3]` sans conversion |
| `TestComputeSineSumInvalidType` | Rejet du type non supporté `"complex"` |
| `TestComputeSineSumEmpty` | Somme = 0 pour un tableau vide (int et float) |
| `TestComputeSineSumNegative` | Entiers négatifs `[-1, 0, 1]`, propriété d'imparité de sin |
| `TestComputeSineSumLargeFloat` | Stabilité numérique avec `1e15` (réduction d'argument) |
| `TestComputeSineSumIntWrongData` | Dispatch `"int"` + données `[]float64` → erreur |
| `TestComputeSineSumFloatWrongData` | Dispatch `"float"` + données `[]int` → erreur |
| `TestRunInt` | `run("int")` retourne un résultat sans erreur |
| `TestRunFloat` | `run("float")` retourne un résultat sans erreur |
| `TestRunInvalidType` | `run("complex")` retourne une erreur |
| `TestMainFunction` | `main()` sans panique avec `--type=float` |
| `TestMainFunctionError` | `main()` sans panique avec type invalide |

## Résultats clés

22 sous-benchmarks mesurés sur un Intel i5-10300H @ 2.50 GHz (amd64, Windows), médianes calculées sur 6 exécutions avec `benchstat`.

| Métrique | Valeur |
|----------|--------|
| Temps par sinus — Int | ~38.7 ns |
| Temps par sinus — Float | ~20.9 ns |
| Ratio Int/Float à 100 % | 1.85× |
| Allocations mémoire | 0 B/op |
| Couverture de code | 100 % |

## Qualité du code

Le code respecte `gofmt` et passe `go vet` sans avertissement. La graine `rand.NewSource(42)` garantit des tableaux identiques d'une exécution à l'autre pour assurer la reproductibilité des benchmarks.

## Liens

- [Rapport TN4](TN4-report.md)
- [Résultats et captures](docs/benchmark-results.md)
- [Calculs détaillés (Q1 & Q2)](docs/calculs.md)
- [Guide graphique gonum/plot](docs/chart-guide.md)
- [Données brutes benchmarks](logs/bench-raw.txt)
- [Prompts IA](TN4-AI-Prompts.md)
- [Dépôt GitHub](https://github.com/moyamelissa/Advanced-Programming/tree/main/TN4)
