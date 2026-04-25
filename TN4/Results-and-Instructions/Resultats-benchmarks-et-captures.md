# TN4 – Résultats des benchmarks

## Comment reproduire les résultats

Depuis le dossier `TN4/`, lancer les commandes suivantes dans un terminal.

**Tests unitaires (6 tests)**

```bash
go test -v -run="Test" ./...
```

![Tests unitaires](Results-and-Instructions/Tests%20unitaires.PNG)

**Benchmarks complets (22 sous-benchmarks)**

```bash
go test -bench="Benchmark" -benchmem -run="^$" -count=1 ./...
```

Le flag `-run="^$"` exclut les tests unitaires, `-benchmem` active le reporting mémoire, et `-count=1` évite les répétitions inutiles.

![Benchmarks complets](Results-and-Instructions/Benchmarks%20complets.PNG)

**Couverture de code**

```bash
go test -v -cover ./...
```

![Couverture de code](Results-and-Instructions/Couverture%20de%20code.PNG)

## Tableau des résultats

| % du tableau | Éléments | Int (ms) | Float (ms) | Ratio Int/Float |
|:---:|:---:|:---:|:---:|:---:|
| 1 % | 10 000 | 0.39 | 0.21 | 1.86× |
| 10 % | 100 000 | 3.81 | 2.11 | 1.81× |
| 20 % | 200 000 | 7.55 | 4.24 | 1.78× |
| 30 % | 300 000 | 10.92 | 6.71 | 1.63× |
| 40 % | 400 000 | 15.48 | 8.38 | 1.85× |
| 50 % | 500 000 | 18.38 | 10.60 | 1.73× |
| 60 % | 600 000 | 21.71 | 12.79 | 1.70× |
| 70 % | 700 000 | 25.46 | 14.72 | 1.73× |
| 80 % | 800 000 | 29.09 | 17.07 | 1.70× |
| 90 % | 900 000 | 32.58 | 19.10 | 1.71× |
| 100 % | 1 000 000 | 36.03 | 20.94 | 1.72× |

Les valeurs en millisecondes sont converties depuis les ns/op affichés par `go test`. Par exemple, `386228 ns/op` donne `0.39 ms`. Aucune allocation mémoire n'a été mesurée (0 B/op, 0 allocs/op) pour les deux types.

## Graphique

Le graphique est généré avec Mermaid (syntaxe `xychart-beta`), qui est rendu automatiquement par GitHub dans les fichiers Markdown.

Les données proviennent directement de la colonne `ns/op` de la sortie des benchmarks. Pour convertir en millisecondes, on divise par 1 000 000. Par exemple, pour `BenchmarkSineSumInt/1pct-8` qui affiche `386228 ns/op`, on obtient `386228 / 1000000 = 0.39 ms`.

Les 11 valeurs Int et les 11 valeurs Float sont ensuite placées dans deux tableaux `line [...]` dans le bloc Mermaid, dans le même ordre que les pourcentages sur l'axe X.

```mermaid
%%{init: {'theme': 'default', 'themeVariables': {'xyChart': {'backgroundColor': '#ffffff'}}}}%%
xychart-beta
    title "Graphique 1 – Temps de calcul selon le pourcentage du tableau (Int vs Float)"
    x-axis "Pourcentage du tableau" ["1%", "10%", "20%", "30%", "40%", "50%", "60%", "70%", "80%", "90%", "100%"]
    y-axis "Temps d'exécution (ms)" 0 --> 40
    line [0.39, 3.81, 7.55, 10.92, 15.48, 18.38, 21.71, 25.46, 29.09, 32.58, 36.03]
    line [0.21, 2.11, 4.24, 6.71, 8.38, 10.60, 12.79, 14.72, 17.07, 19.10, 20.94]
```

La courbe du haut correspond aux entiers (Int), celle du bas aux flottants (Float). Les deux progressent linéairement, ce qui confirme la complexité O(n). Le ratio moyen Int/Float est de 1.75×, principalement dû à la conversion `float64(v)` exécutée à chaque itération pour les entiers.

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
