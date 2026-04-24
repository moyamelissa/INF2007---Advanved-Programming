# TN4 — Rapport : Mesurez et optimisez la somme des sinus en Go

**Cours** : INF2007 — Programmation avancée  
**Plateforme** : Windows/amd64, Intel Core i5-10300H @ 2.50 GHz, 8 threads

---

## 1. Résultats des benchmarks

Les benchmarks ont été exécutés avec `go test -bench=. -benchmem` sur un tableau de 1 000 000 éléments. Aucune allocation mémoire n'a été observée (0 B/op, 0 allocs/op) pour les deux types.

### Tableau comparatif (ns/op)

| % du tableau | Éléments  | Int (ns/op)  | Float (ns/op) | Ratio Int/Float |
|:------------:|:---------:|:------------:|:--------------:|:---------------:|
| 1 %          | 10 000    | 437 494      | 342 851        | 1.28×           |
| 10 %         | 100 000   | 4 410 294    | 3 613 594      | 1.22×           |
| 20 %         | 200 000   | 9 470 669    | 6 711 544      | 1.41×           |
| 30 %         | 300 000   | 14 053 830   | 11 333 780     | 1.24×           |
| 40 %         | 400 000   | 18 323 905   | 16 646 386     | 1.10×           |
| 50 %         | 500 000   | 21 718 398   | 16 902 883     | 1.28×           |
| 60 %         | 600 000   | 25 966 251   | 20 248 249     | 1.28×           |
| 70 %         | 700 000   | 33 715 978   | 23 806 892     | 1.42×           |
| 80 %         | 800 000   | 38 586 420   | 33 308 424     | 1.16×           |
| 90 %         | 900 000   | 43 307 119   | 36 601 474     | 1.18×           |
| 100 %        | 1 000 000 | 48 182 164   | 38 946 579     | 1.24×           |

### Observations clés

- **Scalabilité linéaire** : le temps d'exécution croît linéairement avec la taille du tableau pour les deux types, ce qui confirme la complexité O(n) de `computeSineSum`.
- **Les flottants sont ~20-30 % plus rapides que les entiers**. Cela s'explique par le coût de la conversion `float64(v)` nécessaire pour les entiers avant d'appeler `math.Sin`, alors que les flottants sont déjà au bon type.
- **Zéro allocation mémoire** : les deux implémentations n'allouent aucune mémoire sur le tas, ce qui est optimal.

## 2. Analyse des facteurs de performance

### Coût de la conversion de type (int → float64)
Pour les entiers, chaque itération nécessite `math.Sin(float64(v))`, ajoutant une instruction de conversion entier vers flottant. Pour les `float64`, `math.Sin(v)` s'exécute directement. Ce surcoût explique la différence constante de ~20-30 %.

### Coût de `math.Sin`
La fonction `math.Sin` de Go utilise une approximation polynomiale (série de Taylor / Chebyshev) et domine largement le temps de calcul. L'addition (`sum +=`) est négligeable en comparaison.

### Temps par sinus
- **Int** : 48 182 164 ns / 1 000 000 = **~48.2 ns/sin**
- **Float** : 38 946 579 ns / 1 000 000 = **~38.9 ns/sin**

## 3. Questions spéciales

### Distance parcourue par la lumière pendant un seul calcul de sinus

La lumière voyage à $c = 299\,792\,458$ m/s.

- Pour un sinus sur un **entier** (~48.2 ns) :  
  $d = c \times t = 299\,792\,458 \times 48.2 \times 10^{-9} \approx 14.45$ mètres

- Pour un sinus sur un **flottant** (~38.9 ns) :  
  $d = 299\,792\,458 \times 38.9 \times 10^{-9} \approx 11.66$ mètres

**La lumière parcourt environ 12 à 14 mètres pendant le calcul d'un seul sinus.**

### Nombre de sinus calculables dans un tick de jeu vidéo à 120 FPS

Un tick = $\frac{1}{120} \approx 8\,333\,333$ ns.

- **Int** : $\frac{8\,333\,333}{48.2} \approx 172\,890$ sinus par tick
- **Float** : $\frac{8\,333\,333}{38.9} \approx 214\,225$ sinus par tick

**On peut calculer environ 170 000 à 215 000 sinus par tick à 120 FPS** sur cette machine. Cela signifie qu'un jeu vidéo devrait limiter les appels à `math.Sin` ou utiliser des tables de lookup pour les scénarios exigeant un grand nombre de calculs trigonométriques par image.

## 4. Conclusion

Le benchmarking confirme que :
1. La performance est **linéaire en O(n)**, dominée par le coût de `math.Sin`.
2. Les **flottants sont plus performants** que les entiers grâce à l'absence de conversion de type.
3. Le calcul de sinus est relativement coûteux (~40-50 ns), ce qui limite son utilisation intensive dans des applications temps réel.
4. L'utilisation de `testing.B` avec des sous-benchmarks permet une analyse granulaire et reproductible des performances.
