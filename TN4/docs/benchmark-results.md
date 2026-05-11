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

**Reruns isolés**

Les paliers Float 80, 90 et 100 % ont été relancés en isolation
([bench-float-high-rerun.txt](../logs/bench-float-high-rerun.txt)) pour confirmer la saturation de la
bande passante L3. Le palier 70 % a également été relancé séparément
([bench-float70-rerun.txt](../logs/bench-float70-rerun.txt)).

## Tableau des résultats

| % du tableau | Éléments | Int (ms) | Float (ms) | Ratio Int/Float |
|:---:|:---:|:---:|:---:|:---:|
| 1 % | 10 000 | 0.396 | 0.206 | 1.92× |
| 10 % | 100 000 | 3.57 | 2.06 | 1.73× |
| 20 % | 200 000 | 7.12 | 4.16 | 1.71× |
| 30 % | 300 000 | 10.74 | 6.31 | 1.70× |
| 40 % | 400 000 | 14.31 | 8.53 | 1.68× |
| 50 % | 500 000 | 17.90 | 11.43 | 1.57× |
| 60 % | 600 000 | 21.41 | 13.82 | 1.55× |
| 70 % | 700 000 | 25.11 | 15.28 | 1.64× |
| 80 % | 800 000 | 28.65 | 17.66 | 1.62× |
| 90 % | 900 000 | 32.29 | 22.52 | 1.43× |
| 100 % | 1 000 000 | 35.76 | 22.64 | 1.58× |

Les valeurs en millisecondes proviennent des médianes `benchstat` calculées sur 6 exécutions (rerun mai 2026). Par exemple, `benchstat` reporte `35.76m` pour Int/100pct, soit 35.76 ms. Les paliers Int restent stables (±0–7 %). Les paliers Float 30–90 % présentent des IC de ±9–32 % dus au bruit système lors de l'exécution séquentielle des 22 benchmarks (les Float étant plus rapides, une même perturbation a un impact relatif plus grand). Les paliers 80, 90 et 100 % Float proviennent d'un rerun isolé (6 exécutions) qui confirme la saturation de la bande passante L3 aux hauts paliers. Aucune allocation mémoire n'a été mesurée (0 B/op, 0 allocs/op) pour les deux types.

## Graphique

Le graphique est généré avec gonum/plot en Go (`chart/main.go`) et produit `docs/benchmark-chart.png`. Voir [chart-guide.md](chart-guide.md) pour régénérer.

![Graphique 1 – Int vs Float](benchmark-chart.png)

La courbe du haut correspond aux entiers (Int), celle du bas aux flottants (Float). La courbe Int progresse de façon quasi linéaire, tandis que la courbe Float présente de légères irrégularités (notamment aux paliers 30 % et 70 %). Les deux algorithmes sont O(n) : ces écarts proviennent du bruit de mesure. Les benchmarks Float étant plus rapides, une même perturbation (interruption système, throttling thermique du CPU) a un impact relatif plus important. Les intervalles de confiance `benchstat` le confirment : Float/30pct affiche ± 13 % contre ± 1 % pour Int/100pct. Le ratio moyen Int/Float est de 1.58× au palier 100 % (±9 %), principalement dû à la conversion `float64(v)` exécutée à chaque itération pour les entiers. Les paliers Float 80–90–100 % (rerun isolé) confirment la saturation de la bande passante L3 : Float/90 % et Float/100 % sont quasi identiques (22,52 vs 22,64 ms).

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
