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
| 1 % | 10 000 | 0.44 | 0.24 | 1.86× |
| 10 % | 100 000 | 4.09 | 2.11 | 1.94× |
| 20 % | 200 000 | 8.11 | 4.24 | 1.91× |
| 30 % | 300 000 | 11.83 | 7.79 | 1.52× |
| 40 % | 400 000 | 15.52 | 8.99 | 1.73× |
| 50 % | 500 000 | 19.28 | 11.98 | 1.61× |
| 60 % | 600 000 | 22.98 | 13.61 | 1.69× |
| 70 % | 700 000 | 26.58 | 14.69 | 1.81× |
| 80 % | 800 000 | 30.94 | 16.82 | 1.84× |
| 90 % | 900 000 | 34.78 | 18.96 | 1.83× |
| 100 % | 1 000 000 | 38.71 | 20.93 | 1.85× |

Les valeurs en millisecondes proviennent des médianes `benchstat` calculées sur 6 exécutions. Par exemple, `benchstat` reporte `38.71m` pour Int/100pct, soit 38.71 ms. Les paliers 90–100 % affichent une variation de ± 1 %, confirmant la stabilité des mesures. Aucune allocation mémoire n'a été mesurée (0 B/op, 0 allocs/op) pour les deux types.

## Graphique

Le graphique est généré avec gonum/plot en Go (`chart/main.go`) et produit `docs/benchmark-chart.png`. Voir [chart-guide.md](chart-guide.md) pour régénérer.

![Graphique 1 – Int vs Float](benchmark-chart.png)

La courbe du haut correspond aux entiers (Int), celle du bas aux flottants (Float). La courbe Int progresse de façon quasi linéaire, tandis que la courbe Float présente de légères irrégularités (notamment aux paliers 30 % et 70 %). Les deux algorithmes sont O(n) : ces écarts proviennent du bruit de mesure. Les benchmarks Float étant plus rapides, une même perturbation (interruption système, throttling thermique du CPU) a un impact relatif plus important. Les intervalles de confiance `benchstat` le confirment : Float/30pct affiche ± 7 % contre ± 1 % pour Int/100pct. Le ratio moyen Int/Float est de 1.85× au palier 100 % (± 1 %), principalement dû à la conversion `float64(v)` exécutée à chaque itération pour les entiers.

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
