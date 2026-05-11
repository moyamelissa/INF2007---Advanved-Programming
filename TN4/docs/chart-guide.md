# Comment créer le graphique avec gonum/plot (Go)

Le graphique est généré en Go avec la bibliothèque `gonum.org/v1/plot`. Il produit
un PNG avec zone translucide entre les deux courbes, points avec anneau blanc,
annotation du ratio et labels de valeurs finales.

## Étape 1. Lancer les benchmarks et noter les médianes

```bash
cd TN4
go test -bench="Benchmark" -benchmem -run="^$" -count=6 ./... | tee bench-raw.txt
benchstat bench-raw.txt
```

Relever les médianes `ns/op` de `benchstat` et les diviser par 1 000 000 pour obtenir les millisecondes. Mettre à jour les variables `intMS` et `floatMS` dans `chart/main.go`.

## Étape 2. Structure du module chart

```
TN4/
└── chart/
    ├── go.mod    ← module chart, require gonum.org/v1/plot v0.14.0
    └── main.go   ← programme générateur
```

Le module est **séparé** du module principal `sinesum` pour éviter d'alourdir les dépendances de production.

## Étape 3. Points clés de main.go

| Élément | API gonum/plot |
|---------|----------------|
| Deux courbes colorées | `plotter.NewLine()` + `plotter.NewScatter()` |
| Zone translucide entre les courbes | `plotter.NewPolygon()` avec `color.NRGBA{A: 75}` |
| Points avec anneau blanc | deux `NewScatter` superposés (blanc rayon 4, couleur rayon 3) |
| Légende | `p.Legend.Add()` |
| Labels de valeurs finales | `plotter.NewLabels()` décalés à droite de la dernière colonne |
| Annotation ratio | `plotter.NewLabels()` positionné dans le milieu de la zone fill |
| Ticks X personnalisés | `plot.ConstantTicks([]plot.Tick{...})` |
| Sortie PNG | `p.Save(28*vg.Centimeter, 15*vg.Centimeter, path)` |

La zone translucide est construite en concaténant la courbe Int (gauche → droite) et la courbe Float (droite → gauche), puis en fermant le polygone sur le premier point.

## Étape 4. Générer le PNG

```bash
cd TN4/chart
go mod tidy          # première fois seulement
go run .
```

Sortie : `TN4/docs/benchmark-chart.png`

Le rapport `TN4-report.md` l'intègre avec :

```markdown
![Graphique 1 – Int vs Float](docs/benchmark-chart.png)
```

## Référence

- gonum/plot : https://pkg.go.dev/gonum.org/v1/plot
- Code source : [`chart/main.go`](../chart/main.go)
