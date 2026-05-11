# Mesure et optimisation de la somme des sinus en Go

**INF2007 – Programmation Avancée · TN4 · Melissa Moya**

---

## 1. Approche et structure du programme

Le programme accepte un paramètre `--type` pour choisir entre entiers et flottants
et génère un tableau de 1 000 000 d'éléments avec `rand.NewSource(42)` pour garantir
la reproductibilité. Le léger biais de `rand.Intn(1001)` reste négligeable, l'objectif
étant des données représentatives plutôt qu'une distribution parfaitement uniforme.

Deux fonctions spécialisées, `computeSineSumInt` et `computeSineSumFloat`, isolent la
boucle de calcul par type afin que les benchmarks mesurent uniquement le coût de
`math.Sin`. Une fonction générique `computeSineSum` reste disponible pour la CLI, avec
un surcoût d'environ 1 ns, négligeable face aux 20 à 40 ns de `math.Sin`.

La logique est contenue dans `run()`, qui retourne un struct `RunResult` sans
affichage, et l'affichage est délégué entièrement à `main()`, ce qui facilite les
tests unitaires.

## 2. Résultats des benchmarks

Les 22 sous-benchmarks (11 paliers × 2 types) ont été exécutés avec
`go test -bench=. -benchmem -count=6` sur un Intel i5-10300H à 2,50 GHz
(Windows/amd64) et analysés avec `benchstat` pour obtenir les médianes et
intervalles de confiance à 95 %. Pour isoler le pur coût de calcul, les tableaux
sont pré-générés hors de `b.N`, `b.ResetTimer()` exclut le setup, et
`b.ReportAllocs()` confirme 0 B/op et 0 allocs/op pour les deux types.

**Tableau 1 – Temps de calcul par type et pourcentage du tableau**

| % du tableau | Éléments | Int (ms) | Float (ms) | Ratio |
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

Les flottants sont systématiquement plus rapides, avec un ratio variant de 1,4 à 1,9×
selon le palier. De 50 % à 100 %, le temps double presque exactement pour les entiers
(17,90 à 35,76 ms, facteur 2,0×), ce qui confirme la complexité O(n). Les paliers
Float présentent des intervalles de confiance plus élevés que les paliers Int, ce qui
est attendu puisqu'une même perturbation système a un impact relatif plus important
sur des mesures plus courtes.

**Graphique 1 – Temps de calcul selon le pourcentage du tableau (Int vs Float)**

![Graphique 1 – Int vs Float](docs/benchmark-chart.png)

*La zone mauve translucide représente l'écart Int − Float et l'annotation indique
le ratio à 100 %.*

## 3. Analyse des résultats

**Complexité et écart Int/Float**

Les flottants sont 1,4 à 1,9× plus rapides. L'écart provient de la conversion
`float64(v)` à chaque itération et du coût de réduction d'argument de `math.Sin`
pour les grands entiers : sin(500) nécessite environ 79 réductions modulo 2π, alors
que les flottants dans [0, 1) se situent déjà dans le domaine principal et n'en
requièrent aucune. De 50 % à 100 %, la complexité O(n) est confirmée par un
doublement quasi exact du temps pour les entiers.

**Hiérarchie mémoire**

L'i5-10300H dispose de 256 KB de L1, 1 MB de L2 et 8 MB de L3. Les petits paliers
(1 à 10 %, soit 80 à 800 KB) tiennent entièrement en L2. À partir de 30 %
(environ 2,4 MB), le L2 est dépassé et le L3 prend le relais. À partir de 90 %
(7,2 MB), le tableau approche la limite du L3 (8 MB) et les deux paliers 90 % et
100 % affichent des temps quasi identiques pour Float (22,52 vs 22,64 ms, confirmés
par rerun isolé), ce qui correspond à l'effet de saturation de la bande passante L3.
Une fois la taille du tableau proche du cache, ajouter 10 % d'éléments
supplémentaires ne crée plus de ralentissement proportionnel.

**Mécanismes micro-architecturaux**

`go build -gcflags=-m` confirme que `math.Sin` est inliné dans `sinesum.go`,
éliminant le surcoût d'appel de fonction. La boucle `sum += math.Sin(v)` crée une
dépendance de données sur l'accumulateur `sum` : chaque addition doit attendre le
résultat précédent, ce qui empêche le pipeline FPU de superposer les opérations,
comme dans l'exemple `prefixSum` du chapitre 6. Les appels `math.Sin(v)` sur des
valeurs indépendantes bénéficient de l'exécution dans le désordre, mais ce gain est
plafonné par la latence intrinsèque de `math.Sin`, estimée entre 20 et 40 ns.
Les benchmarks Float, étant plus courts, sont davantage sensibles au bruit système
lors d'une exécution séquentielle des 22 sous-benchmarks. Pour des mesures plus
stables, il est préférable de relancer les paliers suspects en isolation avec
`-bench="BenchmarkSineSumFloat/90pct" -count=6`. Utiliser plusieurs accumulateurs
fusionnés en fin de calcul réduirait la dépendance de données et constituerait la
principale piste d'optimisation.


## 4. Applications numériques

En prenant les médianes benchstat à 100 % du tableau, on obtient un temps par sinus
de 35,8 ns pour les entiers (35 760 000 ns ÷ 1 000 000 ≈ 35,8) et de 22,6 ns pour
les flottants (22 640 000 ns ÷ 1 000 000 ≈ 22,6).

**Q1 — Distance parcourue par la lumière pendant un sinus**

En prenant $c = 299\,792\,458$ m/s, la lumière parcourt la distance suivante pendant
le calcul d'un seul sinus.

| Type | Temps | Distance |
|:---:|:---:|:---:|
| Int | 35,8 ns | $299\,792\,458 \times \frac{35{,}8}{10^9} \approx$ **10,7 m** |
| Float | 22,6 ns | $299\,792\,458 \times \frac{22{,}6}{10^9} \approx$ **6,8 m** |


**Q2 — Nombre de sinus par tick à 120 fps**

Un tick à 120 images par seconde dure $\frac{1}{120}$ s, soit 8 333 333 ns. En
divisant cette durée par le temps d'un sinus, on obtient le nombre de calculs
possibles par frame.

| Type | Temps/sinus | Sinus/tick |
|:---:|:---:|:---:|
| Int | 35,8 ns | $8\,333\,333 \div 35{,}8 \approx$ **233 000** |
| Float | 22,6 ns | $8\,333\,333 \div 22{,}6 \approx$ **369 000** |

### Bibliographie

- Manuel INF2007, chapitres 5 et 6.
- Documentation Go gonum/plot v0.14.0 : https://pkg.go.dev/gonum.org/v1/plot (générateur du Graphique 1)
- Documentation Go Testing : https://pkg.go.dev/testing
- A Tour of Go : https://tour.golang.org/
- GitHub Copilot utilisé comme assistant avec vérification systématique des suggestions. Prompts documentés dans [TN4-AI-Prompts.md](TN4-AI-Prompts.md).
