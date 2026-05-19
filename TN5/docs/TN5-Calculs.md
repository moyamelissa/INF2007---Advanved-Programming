# Calculs détaillés — TN5

Les valeurs utilisées ci-dessous proviennent des médianes `benchstat` sur 6 exécutions
(`-count=6`), mesurées sur un Intel i5-10300H à 2,50 GHz (4 cœurs physiques, 8 threads
logiques, Windows/amd64), contenu de test de 100 000 mots (~700 000 caractères).

---

## 1. Calcul de l'accélération (speedup)

### Formule générale

```
Accélération = Temps de référence / Temps mesuré
```

### Accélération du segment optimal (CountWordsConcurrent)

Référence séquentielle (benchstat) : 5.08 ms  
Optimum concurrent — segment 50 000 chars (benchstat) : 1.78 ms

```
Accélération = 5.08 / 1.78 = 2.85×
```

### Accélération du worker pool par rapport à goroutine-par-segment

Chaque palier compare `BenchmarkWorkerPool` (pool) à `BenchmarkWorkerCount`
(goroutine-par-segment via `CountWordsConcurrentN`).

| Workers | WorkerCount (ms) | WorkerPool (ms) | Calcul          | Accélération |
|:-------:|:----------------:|:---------------:|:---------------:|:------------:|
| 1       | 10.39            | 2.864           | 10.39 / 2.864   | 3.63×        |
| 2       | 9.021            | 2.078           | 9.021 / 2.078   | 4.34×        |
| 4       | 8.473            | 1.705           | 8.473 / 1.705   | 4.97×        |
| 8       | 6.926            | 1.590           | 6.926 / 1.590   | 4.36×        |
| 16      | 6.307            | 1.463           | 6.307 / 1.463   | 4.31×        |
| 32      | 6.208            | 1.533           | 6.208 / 1.533   | 4.05×        |

Fourchette : 3.63× à 4.97×, d'ou l'affirmation "4 a 5× plus rapide" dans le rapport.

---

## 2. Pourquoi le worker pool est plus rapide a nombre de workers égal

La différence vient de la taille des segments utilisés par chaque implémentation.

### CountWordsConcurrentN avec 1 worker

```
Taille contenu       : ~700 000 chars
segmentSize          = len(content) / numWorkers = 700 000 / 1 = 700 000 chars
Nombre de segments   = 1
Nombre de goroutines = 1  (comportement quasi-séquentiel)
Temps mesuré         = 10.39 ms
```

### CountWordsConcurrentPool avec 1 worker

```
Taille contenu       : ~700 000 chars
segmentSize fixe     = 50 000 chars  (optimum déterminé par BenchmarkSegmentSize)
Nombre de segments   = 700 000 / 50 000 = ~14 segments
Nombre de goroutines = 1 worker, mais traite 14 segments en séquence
Temps mesuré         = 2.864 ms
```

Le pool produit ~14 segments de taille optimale même avec 1 seul worker,
ce qui exploite mieux le découpage interne de `strings.Fields`. `CountWordsConcurrentN`
ajuste la taille des segments au nombre de workers, ce qui donne un seul segment
de 700 000 chars quand `numWorkers = 1`.

---

## 3. Application de la loi d'Amdahl

### Formule

```
Speedup(N) = 1 / (S + (1 - S) / N)

où  S = fraction séquentielle du programme
    N = nombre de goroutines
```

Pour trouver S, on résout l'équation a partir du speedup maximum observé.

### Fraction séquentielle pour CountWordsConcurrentN

Speedup maximum observé : passage de 10.39 ms (1 worker) a 6.21 ms (32 workers).

```
Speedup = 10.39 / 6.21 = 1.67×

1.67 = 1 / (S + (1 - S) / 32)
1.67 × (S + (1 - S) / 32) = 1
1.67S + 1.67(1 - S) / 32 = 1
1.67S + 0.0522 - 0.0522S  = 1
1.618S                     = 0.948
S                         ≈ 0.586 ≈ 60 %
```

60 % du travail est séquentiel, seulement 40 % est parallélisable.

### Fraction séquentielle pour CountWordsConcurrentPool

Speedup maximum observé : passage de 2.864 ms (1 worker) a 1.463 ms (16 workers).

```
Speedup = 2.864 / 1.463 = 1.96×

1.96 = 1 / (S + (1 - S) / 16)
1.96 × (S + (1 - S) / 16) = 1
1.96S + 0.1225 - 0.1225S  = 1
1.8375S                    = 0.8775
S                         ≈ 0.477 ≈ 48 %
```

48 % du travail est séquentiel, 52 % est parallélisable.

Le worker pool parallélise davantage (52 % contre 40 %) car il supprime la
recréation de goroutines a chaque appel, réduisant la part séquentielle de
l'orchestration.

---

## 4. Calcul du point de dégradation

### Pourquoi 500 chars est pire que le séquentiel

```
Taille contenu           : ~700 000 chars
Taille segment           :     500 chars
Nombre de goroutines     : 700 000 / 500 = ~1 400 goroutines
Coût création goroutine  : ~3 µs (estimation courante pour le runtime Go)
Coût total création      : 1 400 × 3 µs = 4 200 µs = 4.2 ms

Temps mesuré (500 chars) : 8.87 ms
Temps séquentiel         : 5.08 ms
Surcoût observé          : 8.87 - 5.08 = 3.79 ms  ≈  overhead de création estimé
```

### Seuil de rentabilité d'un segment

Un segment est rentable lorsque le temps de calcul qu'il génère dépasse
le coût de création de la goroutine associée (~3 µs).

Segment de 50 000 chars (optimum) :

```
Temps total mesuré    : 1.78 ms
Nombre de segments    : ~14
Temps par segment     : 1 780 µs / 14 ≈ 127 µs
Overhead création     :   ~3 µs
Rapport signal/bruit  : 127 / 3 ≈ 42×   → très rentable
```

Segment de 500 chars (dégradation) :

```
Temps séquentiel      : 5 080 µs
Nombre de segments    : ~1 400
Temps par segment     : 5 080 / 1 400 ≈ 3.6 µs
Overhead création     :   ~3 µs
Rapport signal/bruit  : 3.6 / 3 ≈ 1.2×  → pas rentable
```

---

## 5. Évolution des allocations mémoire

### CountWordsConcurrentN

| Workers | allocs/op | Delta par rapport a 1 worker |
|:-------:|:---------:|:-----------------------------:|
| 1       | 4         | —                             |
| 2       | 7         | +3                            |
| 4       | 12        | +8                            |
| 8       | 21        | +17                           |
| 16      | 38        | +34                           |
| 32      | 71        | +67                           |

Chaque goroutine supplémentaire alloue environ 2 objets (pile initiale + entrée
dans le canal de résultats). La croissance est quasi-linéaire avec le nombre de
workers, ce qui confirme que chaque goroutine est recréée a chaque appel.

### CountWordsConcurrentPool

| Workers | allocs/op | Delta par rapport a 1 worker |
|:-------:|:---------:|:-----------------------------:|
| 1       | 23        | —                             |
| 2       | 24        | +1                            |
| 4       | 26        | +3                            |
| 8       | 30        | +7                            |
| 16      | 38        | +15                           |
| 32      | 55        | +32                           |

Le pool alloue davantage au départ (23 contre 4) car il crée les canaux
`jobs` et `results` ainsi que les ~14 entrées de segments des le lancement.
En revanche, la croissance avec le nombre de workers est beaucoup plus faible
(+32 pour passer de 1 a 32 workers, contre +67 pour `CountWordsConcurrentN`)
car le nombre de segments reste fixe.
