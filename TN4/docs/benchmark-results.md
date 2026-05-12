# TN4 – Résultats des benchmarks

## Comment reproduire les résultats

Depuis le dossier `TN4/`, lancer les commandes suivantes dans un terminal.

**Tests unitaires (13 tests)**

```bash
go test -v -run="Test" ./...
```

Résultats : [tests.txt](../logs/tests.txt)

**Benchmarks complets (22 sous-benchmarks, 6 exécutions analysées par benchstat)**

```bash
go test -bench="Benchmark" -benchmem -run="^$" -count=6 ./...
benchstat bench-raw.txt
```

Le flag `-run="^$"` exclut les tests unitaires, `-benchmem` active le reporting mémoire, et `-count=6` fournit 6 échantillons pour que `benchstat` calcule les médianes et intervalles de confiance à 95 %.

Résultats bruts : [bench-raw.txt](../logs/bench-raw.txt)
Analyse benchstat : [benchstat.txt](../logs/benchstat.txt)

**Couverture de code**

```bash
go test -v -cover ./...
```

Résultats : [coverage.txt](../logs/coverage.txt)

## Tableau des résultats

| % du tableau | Éléments | Int (ms) | Float (ms) | Ratio Int/Float |
|:---:|:---:|:---:|:---:|:---:|
| 1 % | 10 000 | 0.373 | 0.223 | 1.67× |
| 10 % | 100 000 | 3.69 | 2.21 | 1.67× |
| 20 % | 200 000 | 7.38 | 4.46 | 1.66× |
| 30 % | 300 000 | 11.06 | 6.57 | 1.68× |
| 40 % | 400 000 | 14.71 | 8.81 | 1.67× |
| 50 % | 500 000 | 18.49 | 10.87 | 1.70× |
| 60 % | 600 000 | 22.04 | 13.35 | 1.65× |
| 70 % | 700 000 | 25.86 | 15.54 | 1.66× |
| 80 % | 800 000 | 29.47 | 17.82 | 1.65× |
| 90 % | 900 000 | 33.44 | 19.71 | 1.70× |
| 100 % | 1 000 000 | 36.91 | 21.72 | 1.70× |

Les valeurs en millisecondes proviennent des médianes `benchstat` calculées sur 6 exécutions (mai 2026). Par exemple, `benchstat` reporte `36.91m` pour Int/100pct, soit 36.91 ms. La plupart des paliers affichent des IC de ±1–5 %. Quelques paliers Float (20 %, 60 %, 90 %, 100 %) et Int/30pct affichent des IC de ±8–11 %, dus au bruit système lors de l'exécution séquentielle des 22 benchmarks (les Float étant plus rapides, une même perturbation a un impact relatif plus grand). Aucune allocation mémoire n'a été mesurée (0 B/op, 0 allocs/op) pour les deux types.

## Graphique

Le graphique est généré avec gonum/plot en Go (`chart/main.go`) et produit `docs/benchmark-chart.png`. Voir [chart-guide.md](chart-guide.md) pour régénérer.

![Graphique 1 – Int vs Float](benchmark-chart.png)

La courbe du haut correspond aux entiers (Int), celle du bas aux flottants (Float). Les deux courbes progressent de façon quasi linéaire, avec de légères irrégularités dues au bruit de mesure. Les deux algorithmes sont O(n) : les écarts proviennent du bruit système lors de l'exécution séquentielle. Les intervalles de confiance `benchstat` le confirment : Float/60pct affiche ±…11 % contre ±…1 % pour Int/70pct. Le ratio Int/Float est de 1.70× au palier 100 %, principalement dû à la conversion `float64(v)` exécutée à chaque itération pour les entiers.

## Lecture des résultats

Chaque ligne de la sortie `go test` se lit comme suit :

```
BenchmarkSineSumInt/40pct-8     85     17215340 ns/op     0 B/op     0 allocs/op
│                        │       │     │                  │          │
│                        │       │     │                  │          └─ Allocations par opération
│                        │       │     │                  └─ Mémoire allouée par opération
│                        │       │     └─ Nanosecondes par opération
│                        │       └─ Nombre d'itérations exécutées
│                        └─ Nombre de threads (GOMAXPROCS)
└─ Nom du benchmark / sous-benchmark
```

Le framework `testing.B` ajuste automatiquement le nombre d'itérations (`b.N`) pour obtenir une mesure stable. Plus le benchmark est lent, moins il y a d'itérations.
