# INF2007 – TN4 – Melissa Moya

## Approche et structure du programme

Le programme calcule la somme des sinus d'un tableau de 1 000 000 d'éléments, en entiers ou en flottants selon le flag `--type`. J'ai séparé le code en trois couches. `generateIntArray` et `generateFloatArray` créent les tableaux avec `rand.NewSource(42)` pour que chaque exécution produise exactement les mêmes données, ce qui est essentiel pour la reproductibilité des benchmarks. J'aurais pu utiliser `crypto/rand`, mais celui-ci fait des appels système à chaque tirage, ce qui fausserait les mesures en mélangeant le coût du calcul avec celui de la génération. `computeSineSumInt` et `computeSineSumFloat` contiennent la boucle de calcul spécialisée pour chaque type, et `computeSineSum` sert de dispatch via un `switch` sur le type reçu en `interface{}`. Les benchmarks passent par `computeSineSum` pour mesurer le programme tel qu'il est réellement exécuté, dispatch inclus. Le surcoût du `switch` et de l'assertion de type est négligeable face à `math.Sin` (environ 1 ns contre 30 ns), mais c'est plus fidèle à l'utilisation réelle.

## Résultats des benchmarks

Les mesures reposent exclusivement sur `testing.B`, le seul outil fiable pour du benchmarking en Go. Les `time.Since` dans `main` ne servent qu'à donner une intuition à l'utilisateur et ne sont pas utilisées dans l'analyse. `testing.B` ajuste automatiquement le nombre d'itérations (`b.N`) pour stabiliser la mesure, et `b.ResetTimer()` est appelé avant chaque boucle pour exclure le setup. Les 22 sous-benchmarks (11 paliers par type) ont été exécutés avec `go test -bench=. -benchmem -count=1` sur un Intel i5-10300H à 2.50 GHz sous Windows/amd64 avec 8 threads. Aucune allocation mémoire n'a été détectée (0 B/op), ce qui confirme que `sum` et les itérateurs restent sur la pile. En complément des 4 tests unitaires demandés, j'ai ajouté deux tests sur des valeurs négatives et des flottants extrêmes (`1e15`) pour valider la robustesse numérique de `math.Sin`, notamment sa réduction d'argument sur de grands multiples de pi.

**Tableau 1 – Temps de calcul par type et pourcentage du tableau (1 000 000 éléments)**

| % du tableau | Éléments | Int (ms) | Float (ms) | Ratio |
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

Les flottants sont systématiquement plus rapides avec un ratio moyen de 1.75×. En passant de 50 % à 100 %, le temps double presque exactement (18.38 vers 36.03 ms pour Int, 10.60 vers 20.94 ms), ce qui confirme la complexité O(n). Aucune allocation mémoire n'a été détectée (0 B/op), ce qui signifie que `sum` et les itérateurs restent sur la pile.

**Graphique 1 – Temps de calcul selon le pourcentage du tableau (Int vs Float)**

```mermaid
%%{init: {'theme': 'default', 'themeVariables': {'xyChart': {'backgroundColor': '#ffffff'}}}}%%
xychart-beta
    x-axis "Pourcentage du tableau" ["1%", "10%", "20%", "30%", "40%", "50%", "60%", "70%", "80%", "90%", "100%"]
    y-axis "Temps d'exécution (ms)" 0 --> 40
    line [0.39, 3.81, 7.55, 10.92, 15.48, 18.38, 21.71, 25.46, 29.09, 32.58, 36.03]
    line [0.21, 2.11, 4.24, 6.71, 8.38, 10.60, 12.79, 14.72, 17.07, 19.10, 20.94]
```

La courbe du haut correspond aux entiers, celle du bas aux flottants. La progression est quasi linéaire pour les deux types, et les courbes restent parallèles sur toute la plage. L'écart entre Int et Float s'explique par la conversion `float64(v)` que la version Int exécute à chaque itération. Sur x86-64, cette conversion se traduit par l'instruction `CVTSI2SD` qui ajoute 4 à 5 cycles par élément. Sur 1 million d'éléments à 2.5 GHz, ça représente environ 2 ms de surcoût pur, mais l'écart observé d'environ 15 ms suggère que la conversion perturbe aussi le pipeline du processeur en cassant la chaîne de dépendances de données. `math.Sin` elle-même utilise une réduction de l'argument suivie d'une approximation polynomiale (Chebyshev) et c'est l'opération qui domine le temps de calcul.

## Calculs pour les questions du professeur

Pour obtenir le temps par sinus, on divise le `ns/op` du benchmark à 100 % par le nombre d'éléments. Le benchmark Int à 100 % affiche 36 027 981 ns/op, donc $36\,027\,981 \div 1\,000\,000 = 36.0$ ns par sinus. Pour Float, 20 936 602 ns/op donne $20\,936\,602 \div 1\,000\,000 = 20.9$ ns par sinus.

**Distance parcourue par la lumière pendant un sinus.** La lumière voyage à $c = 299\,792\,458$ m/s. On multiplie par le temps en secondes.

$$d_{int} = 299\,792\,458 \times \frac{36.0}{1\,000\,000\,000} = 10.8 \text{ mètres}$$

$$d_{float} = 299\,792\,458 \times \frac{20.9}{1\,000\,000\,000} = 6.3 \text{ mètres}$$

La lumière parcourt entre 6 et 11 mètres pendant un seul calcul de sinus. C'est la taille d'une pièce d'appartement. On pense que `math.Sin` est instantané, mais la lumière a le temps de traverser un salon pendant ce calcul.

**Nombre de sinus par tick à 120 fps.** Un tick dure $\frac{1}{120} = 8\,333\,333$ ns. On divise par le temps d'un sinus.

$$n_{int} = \frac{8\,333\,333}{36.0} \approx 231\,481 \text{ sinus par tick}$$

$$n_{float} = \frac{8\,333\,333}{20.9} \approx 398\,726 \text{ sinus par tick}$$

En pratique, si on réserve 10 % du budget de frame au calcul de sinus, ça laisse environ 23 100 (Int) ou 39 900 (Float) sinus par tick, ce qui est largement suffisant pour animer un millier d'objets avec des rotations et des oscillations.

*Pour les détails de chaque calcul, voir [Guide-calculs-questions-professeur.md](Results-and-Instructions/Guide-calculs-questions-professeur.md)

### Liens

- Dépôt GitHub [github.com/moyamelissa/Advanced-Programming/tree/main/TN4](https://github.com/moyamelissa/Advanced-Programming/tree/main/TN4)
- Implémentation [sinesum.go](sinesum.go), tests et benchmarks [sinesum_test.go](sinesum_test.go)

### Bibliographie

- Documentation Go `math/rand`, `testing`, `flag` sur https://pkg.go.dev
- Documentation Mermaid XY Chart https://mermaid.js.org/syntax/xyChart.html
- GitHub Copilot, utilisé comme assistant avec vérification systématique des suggestions
