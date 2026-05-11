# Mesure et optimisation de la somme des sinus en Go

**INF2007 – Programmation Avancée · TN4 · Melissa Moya**

---

## 1. Approche et structure du programme

Le programme accepte un paramètre `--type` pour choisir entre entiers et flottants
et génère un tableau de 1 000 000 d'éléments avec `rand.NewSource(42)` pour garantir
la reproductibilité. Le léger biais de `rand.Intn(1001)` reste négligeable, l'objectif
étant des données représentatives plutôt qu'une distribution parfaitement uniforme. Deux fonctions spécialisées, `computeSineSumInt` et `computeSineSumFloat`, isolent la
boucle de calcul par type afin que les benchmarks mesurent uniquement le coût de
`math.Sin`. Une fonction générique `computeSineSum` reste disponible pour la CLI, avec
un surcoût d'environ 1 ns, négligeable face aux 20 à 40 ns de `math.Sin`. La logique est contenue dans `run()`, qui retourne un struct `RunResult` sans
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
| 40 % | 400 000 | 14.31 | 8.56 | 1.67× |
| 50 % | 500 000 | 17.90 | 11.43 | 1.57× |
| 60 % | 600 000 | 21.41 | 13.82 | 1.55× |
| 70 % | 700 000 | 25.11 | 15.91 | 1.58× |
| 80 % | 800 000 | 28.65 | 17.28 | 1.66× |
| 90 % | 900 000 | 32.29 | 21.92 | 1.47× |
| 100 % | 1 000 000 | 35.76 | 21.40 | 1.67× |

Les flottants sont systématiquement plus rapides, avec un ratio variant de 1,5 à 1,9×
selon le palier. De 50 % à 100 %, le temps double presque exactement pour les entiers
(17,90 à 35,76 ms, facteur 2,0×), confirmant la complexité O(n). Les paliers Float 30–90 %
présentent des IC élevés (±9–32 %) lors du run complet en raison du bruit système lors de
l'exécution séquentielle des 22 benchmarks ; un rerun isolé de Float/70pct (6 exécutions)
donne 15,28 ms ±6 %, confirmant la stabilité du code.

**Graphique 1 – Temps de calcul selon le pourcentage du tableau (Int vs Float)**

![Graphique 1 – Int vs Float](docs/benchmark-chart.png)

*La zone mauve translucide représente l'écart Int − Float ; l'annotation « 1.67× » indique le ratio à 100 %.*

## 3. Analyse des résultats

**Complexité et écart Int/Float.**

Les flottants sont 1,5 à 1,9× plus rapides. L'écart provient de la conversion `float64(v)` à chaque itération et du coût de réduction d'argument de `math.Sin` pour les grands entiers, notamment sin(500) réduit ~79 tours modulo 2π, alors que les flottants dans [0, 1) se situent déjà dans le domaine principal et ne requièrent aucune réduction. La variabilité élevée des paliers Float 30–90 % (IC allant jusqu'à ±32 % pour 40 %) s'explique par le bruit système lors de l'exécution séquentielle des 22 benchmarks : les benchmarks Float étant plus rapides, une même perturbation a un impact relatif plus important. Un rerun isolé de Float/70pct (6 exécutions) donne 15,28 ms ±6 %, confirmant la stabilité du code.

**Hiérarchie mémoire — L1/L2/L3.**

L'i5-10300H dispose de L1 = 256 KB, L2 = 1 MB, L3 = 8 MB. Les petits paliers (1–10 %, 80–800 KB) tiennent entièrement en L2 ; le palier 30 % (~2,4 MB) dépasse le L2 et sollicite le L3 ; le tableau complet (1 000 000 × 8 octets = **8 MB**) atteint exactement la limite du L3, limitant les accès par la bande passante, ce qui explique les légères non-linéarités aux paliers 70 %+.

**Mécanismes micro-architecturaux (ch. 6).**

`go build -gcflags=-m` confirme que `math.Sin` est inliné aux lignes 40 et 49 de `sinesum.go`, éliminant le surcoût d'appel de fonction. La boucle `sum += math.Sin(v)` crée une **dépendance de données** sur l'accumulateur `sum` : chaque addition doit attendre le résultat précédent, ce qui empêche le pipeline FPU de superposer les opérations (pas de *pipelining* des additions), identique à l'exemple `prefixSum` du chapitre 6. Les appels `math.Sin(v)` sur des `v` indépendants sont eux exécutés en désordre (*out-of-order execution*), mais ce gain est plafonné par la **latence** intrinsèque de `math.Sin` (~20–40 ns), bien supérieure à son **débit** théorique. Les intervalles de confiance `benchstat` confirment la stabilité globale : Int/100pct affiche ±1 %, les paliers Int restant tous dans ±7 %. Float/70pct a été retesté en isolation (6 exécutions) et affiche ±6 %, confirmant que le ±17 % du run complet est dû au bruit système lors de l'exécution séquentielle. Les paliers Float 30–90 % présentent des IC de ±9–32 %, caractéristique de cette machine sous charge. Utiliser plusieurs accumulateurs fusionnés en fin de calcul réduirait la dépendance de données et constituerait la principale piste d'optimisation.

## 4. Applications numériques

En prenant les médianes benchstat à 100 % du tableau, on obtient un temps par sinus de **35,8 ns** pour les entiers (35 760 000 ns ÷ 1 000 000) et de **21,4 ns** pour les flottants (21 400 000 ns ÷ 1 000 000).

**Q1 — Distance parcourue par la lumière pendant un sinus** ($c = 299\,792\,458$ m/s)

| Type | Temps | Distance |
|:---:|:---:|:---:|
| Int | 35,8 ns | $299\,792\,458 \times \frac{35{,}8}{10^9} \approx$ **10,7 m** |
| Float | 21,4 ns | $299\,792\,458 \times \frac{21{,}4}{10^9} \approx$ **6,4 m** |

**Réponse :** Pendant le calcul d'un sinus, la lumière parcourt environ **10,7 m** avec des entiers et **6,4 m** avec des flottants.

**Q2 — Nombre de sinus par tick à 120 fps** (tick = $\frac{1}{120}$ s = 8 333 333 ns)

| Type | Temps/sinus | Sinus/tick |
|:---:|:---:|:---:|
| Int | 35,8 ns | $8\,333\,333 \div 35{,}8 \approx$ **233 000** |
| Float | 21,4 ns | $8\,333\,333 \div 21{,}4 \approx$ **389 000** |

**Réponse :** À 120 fps, on peut calculer environ **233 000 sinus** par frame avec des entiers et **389 000** avec des flottants.

*Détails : [calculs.md](docs/calculs.md)*


### Bibliographie

- Manuel INF2007, chapitres 5 et 6.
- Documentation Go gonum/plot v0.14.0 : https://pkg.go.dev/gonum.org/v1/plot (générateur du Graphique 1)
- Documentation Go Testing : https://pkg.go.dev/testing
- A Tour of Go : https://tour.golang.org/
- GitHub Copilot utilisé comme assistant avec vérification systématique des suggestions. Prompts documentés dans [TN4-AI-Prompts.md](TN4-AI-Prompts.md).
