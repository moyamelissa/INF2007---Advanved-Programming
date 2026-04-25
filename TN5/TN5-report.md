# TN5 — Rapport : Comptez les mots avec la concurrence en Go

**Cours** : INF2007 — Programmation avancée  
**Plateforme** : Windows/amd64, Intel Core i5-10300H @ 2.50 GHz, 8 threads

---

## 1. Architecture concurrente : goroutines et canaux

Le programme suit le modèle **fan-out / fan-in** :

1. **Lecture** : le fichier est lu en mémoire (`os.ReadFile`).
2. **Division** (`splitIntoSegments`) : le contenu est découpé en segments d'environ N caractères.
3. **Fan-out** : une goroutine est lancée par segment via `go countWordsInSegment(seg, ch)`.
4. **Fan-in** : chaque goroutine envoie son résultat sur un canal buffered. La goroutine principale somme les résultats.

### Décision clé : ne pas couper les mots

Le problème principal de la division en segments est qu'un segment de N caractères peut couper un mot en deux. J'ai résolu cela dans `splitIntoSegments` en avançant la coupure jusqu'à la fin du mot courant :

```go
end := segmentSize
// Avancer jusqu'à la fin du mot courant
for end < len(content) && content[end] != ' ' && content[end] != '\n' &&
    content[end] != '\t' && content[end] != '\r' {
    end++
}
segments = append(segments, content[:end])
// Sauter les espaces entre les segments
for end < len(content) && (content[end] == ' ' || ...) {
    end++
}
content = content[end:]
```

Concrètement, pour le texte `"Hello world from Go"` avec `segmentSize=7` : au lieu de couper en `"Hello w"` et `"orld from Go"` (qui fausserait le comptage), on obtient `"Hello world"` et `"from Go"`. Le test `TestCountWordsConsistency` vérifie que le comptage est **identique** pour 7 tailles de segments différentes (1, 5, 10, 20, 50, 100, 500 caractères).

### Garanties de correction

- **Canal buffered** (`make(chan int, len(segments))`) : dimensionné exactement au nombre de goroutines, aucune ne bloque en écriture, évitant les deadlocks.
- **Pas de mémoire partagée** : chaque goroutine travaille sur sa propre sous-chaîne. Le canal est le seul point de communication, conformément au principe Go « ne communiquez pas en partageant la mémoire ; partagez la mémoire en communiquant ».
- **Synchronisation implicite** : la boucle `for range segments { total += <-ch }` attend automatiquement la fin de chaque goroutine. Pas besoin de `sync.WaitGroup` dans ce cas, car on connaît le nombre exact de résultats à recevoir.

## 2. Résultats des benchmarks (100 000 mots)

### Performance en fonction de la taille des segments

| Taille segment | Goroutines (≈) | ns/op       | Allocs/op | Speedup vs séquentiel |
|:--------------:|:---------------:|:-----------:|:---------:|:---------------------:|
| 10 chars       | ~100 000        | 37 312 282  | 100 037   | 0.12× (plus lent)     |
| 100 chars      | ~13 000         | 5 098 136   | 13 017    | 0.86×                 |
| 500 chars      | ~2 700          | 2 174 825   | 2 701     | 2.02×                 |
| 1 000 chars    | ~1 350          | 1 895 688   | 1 355     | 2.32×                 |
| **5 000 chars**| **~280**        | **1 480 232** | **279** | **2.97× (optimal)**   |
| 10 000 chars   | ~144            | 1 570 571   | 144       | 2.80×                 |
| 50 000 chars   | ~34             | 1 747 235   | 34        | 2.51×                 |
| 100 000 chars  | ~19             | 2 000 042   | 19        | 2.20×                 |
| Tout-en-un     | 1               | 4 374 865   | 4         | 1.00× (séquentiel)    |

### Séquentiel vs concurrent

| Mode               | ns/op       | Speedup |
|:------------------:|:-----------:|:-------:|
| Séquentiel         | 4 393 525   | 1.00×   |
| Concurrent (1 000) | 2 163 028   | 2.03×   |
| Concurrent (10 000)| 1 809 859   | 2.43×   |
| Concurrent (100 000)| 2 205 570  | 1.99×   |

## 3. Analyse : la performance croît-elle linéairement ?

**Non, la performance ne croît pas linéairement avec le nombre de goroutines.** Les résultats montrent une courbe en U inversé avec trois régimes distincts :

### Trop de goroutines (segment = 10 chars → ~100k goroutines)
Le programme est **8× plus lent** que le séquentiel. Chaque goroutine en Go a une pile initiale de ~2-8 Ko. Avec 100 000 goroutines, cela représente ~200 Mo–800 Mo rien qu'en piles. De plus, le runtime Go doit planifier 100 000 goroutines sur 8 threads OS, ce qui génère un surcoût massif de context switching. Les 100 037 allocations mémoire confirment ce problème.

### Point optimal (~5 000 chars → ~280 goroutines)
Le speedup atteint **~3×** sur 8 threads logiques (4 cœurs physiques). Pourquoi 3× et non 4× ? Parce que `strings.Fields` effectue un scan linéaire de la mémoire, et le goulot d'étranglement devient la **bande passante du bus mémoire** (lecture séquentielle depuis le cache L3). Avec 280 goroutines et seulement 279 allocations, le surcoût de coordination est minimal.

### Trop peu de goroutines (segment = 100k chars → ~19 goroutines)
Le speedup redescend à 2.2×. Avec 19 goroutines sur 8 threads, certains cœurs finissent avant les autres et restent inactifs. Le déséquilibre de charge (load imbalance) empêche d'exploiter pleinement le processeur.

### Corrélation allocations / performance

On observe une corrélation directe entre le nombre d'allocations et la dégradation :
- 279 allocs → 1.48 ms (optimal)
- 1 355 allocs → 1.90 ms
- 13 017 allocs → 5.10 ms
- 100 037 allocs → 37.31 ms

Chaque allocation supplémentaire ajoute environ 0.3 µs de surcoût (allocation + GC).

## 4. Tests unitaires

J'ai écrit 6 tests couvrant les cas critiques :

| Test | Pourquoi ce test est important |
|------|-------------------------------|
| `TestCountWordsEmpty` | Vérifie que `""` retourne 0, pas un panic sur un slice vide. |
| `TestCountWordsSingleWord` | Cas minimal : un seul mot sans espace. |
| `TestCountWordsMultipleLines` | Vérifie que `\n` est bien traité comme séparateur (3 lignes × 3 mots = 9). |
| `TestCountWordsMultipleSpaces` | `"  mot1   mot2   mot3  "` → 3, pas 7 (espaces multiples ignorés). |
| `TestCountWordsConsistency` | **Test le plus important** : vérifie que le résultat est identique pour 7 tailles de segment (1 à 500). Si le split coupait un mot, ce test échouerait. |
| `TestSplitIntoSegments` | Vérifie directement que les segments contiennent des mots complets. |

Le test `TestCountWordsConsistency` a été particulièrement utile pendant le développement : une version initiale de `splitIntoSegments` ne sautait pas les espaces entre segments, causant des mots comptés en double.

## 5. Conclusion

Les goroutines et canaux de Go permettent une implémentation concurrente correcte et élégante du comptage de mots. Le choix de la taille de segment est critique : un segment de ~5 000 caractères offre le meilleur compromis entre parallélisme et surcoût sur cette machine (speedup 3× avec 280 goroutines). Le modèle fan-out/fan-in garantit l'absence de race conditions sans nécessiter de mutex. La performance ne croît **pas** linéairement — elle suit une courbe en cloche dictée par l'équilibre entre parallélisme, surcoût de coordination et bande passante mémoire.
