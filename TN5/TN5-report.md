# TN5 — Rapport : Comptez les mots avec la concurrence en Go

**Cours** : INF2007 — Programmation avancée  
**Plateforme** : Windows/amd64, Intel Core i5-10300H @ 2.50 GHz, 8 threads

---

## 1. Architecture concurrente : goroutines et canaux

Le programme suit le modèle **fan-out / fan-in** :

1. **Lecture** : le fichier est lu en mémoire (`os.ReadFile`).
2. **Division** (`splitIntoSegments`) : le contenu est découpé en segments d'environ N caractères, en avançant la coupure jusqu'à la fin du mot courant pour ne jamais couper un mot en deux.
3. **Fan-out** : une goroutine est lancée par segment via `go countWordsInSegment(seg, ch)`.
4. **Fan-in** : chaque goroutine envoie son résultat sur un canal buffered (`chan int`). La goroutine principale lit exactement `len(segments)` valeurs et les somme.

### Garanties de correction

- **Canal buffered** (`make(chan int, len(segments))`) : aucune goroutine ne bloque en écriture, ce qui évite les deadlocks.
- **Pas de mémoire partagée** : chaque goroutine travaille sur sa propre sous-chaîne. Le seul point de communication est le canal, conformément au principe Go « ne communiquez pas en partageant la mémoire ; partagez la mémoire en communiquant ».
- **Synchronisation implicite** : la boucle `for range segments { total += <-ch }` garantit que la goroutine principale attend la fin de toutes les goroutines avant d'afficher le résultat.
- **Pas de coupure de mot** : `splitIntoSegments` avance la position de coupure au-delà du mot en cours, puis saute les espaces, pour que chaque mot apparaisse dans exactement un segment.

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

**Non, la performance ne croît pas linéairement avec le nombre de goroutines.** Les résultats montrent une courbe en U inversé :

- **Trop de goroutines** (segment 10 chars → ~100k goroutines) : le surcoût de création/planification des goroutines et d'envoi/réception sur le canal domine. Le programme est **8× plus lent** que le séquentiel.
- **Point optimal** (~5 000 chars → ~280 goroutines) : le parallélisme compense le surcoût. Speedup de **~3×** sur 8 threads logiques, ce qui est cohérent avec les 4 cœurs physiques.
- **Trop peu de goroutines** : pas assez de parallélisme pour exploiter les cœurs disponibles.

Le speedup maximal (~3×) est inférieur au nombre de cœurs (4 physiques / 8 logiques) car `strings.Fields` est limité par la bande passante mémoire (lecture séquentielle du cache L3), et le surcoût de coordination (canal, allocation des sous-chaînes) consomme du temps CPU.

## 4. Conclusion

Les goroutines et canaux de Go permettent une implémentation concurrente correcte et élégante du comptage de mots. Le choix de la taille de segment est critique : un segment de ~5 000 caractères offre le meilleur compromis entre parallélisme et surcoût sur cette machine. Le modèle fan-out/fan-in garantit l'absence de race conditions sans nécessiter de mutex.
