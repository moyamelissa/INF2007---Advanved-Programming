# Comptage concurrent de mots en Go

**INF2007 – Programmation Avancée · TN5 · Melissa Moya**

---

## 1. Architecture concurrente

Le programme divise le contenu du fichier en segments de taille paramétrable
via `splitIntoSegments`, qui décale chaque coupure vers le prochain espace pour
ne jamais couper un mot. Une goroutine est lancée par segment avec
`go countWordsInSegment`, et chaque goroutine envoie son résultat partiel sur
un canal partagé (`chan int`).

La correction de l'exécution concurrente repose sur trois mécanismes. D'abord,
aucune variable n'est partagée entre les goroutines. Chaque goroutine calcule
son total local dans une variable de pile, puis l'envoie sur le canal, ce qui
élimine toute condition de course sans mutex. Ensuite, le canal est en mémoire
tampon avec une capacité égale au nombre de segments, donc tous les envois
aboutissent sans blocage. Enfin, la boucle `for range segments` dans la goroutine
principale consomme exactement N résultats avant de retourner, ce qui garantit
que toutes les goroutines ont terminé avant la sommation finale. L'addition
étant commutative, l'ordre d'arrivée n'affecte pas le résultat.

## 2. Résultats des benchmarks

Mesures sur un Intel i5-10300H à 2,50 GHz (4 cœurs physiques et 8 threads
logiques, Windows/amd64), contenu de 100 000 mots (~700 000 caractères),
6 exécutions analysées avec `benchstat`.

**Tableau 1 – Temps selon la taille des segments**

| Segment | 500 | 1 000 | 5 000 | 10 000 | 50 000 | 100 000 | Tout |
|:---:|:---:|:---:|:---:|:---:|:---:|:---:|:---:|
| Temps (ms) | 8.87 | 2.84 | 4.84 | 4.85 | **1.78** | 1.97 | 4.09 |

Le séquentiel prend 5.08 ms. L'optimum concurrent est à 50 000 caractères
(1.78 ms, soit 2.85× plus rapide). Avec des segments de 500 caractères
(~1 200 goroutines), la surcharge d'ordonnancement dépasse le bénéfice du
parallélisme.

**Tableau 2 – Temps selon le nombre de goroutines (CountWordsConcurrentN)**

| Workers | 1 | 2 | 4 | 8 | 16 | 32 |
|:---:|:---:|:---:|:---:|:---:|:---:|:---:|
| Temps (ms) | 10.39 | 9.02 | 8.47 | 6.93 | 6.31 | 6.21 |
| Accélération | 1.00× | 1.15× | 1.23× | 1.50× | 1.65× | 1.67× |

## 3. Analyse de la linéarité et worker pool

La performance ne croît pas linéairement. Sur 4 cœurs physiques, l'accélération
atteint 1.67× au lieu des 4× théoriques et plafonne dès 16 workers. D'abord,
`strings.Fields` est trop rapide par goroutine pour amortir le coût de création
et d'orchestration. Ensuite, le canal partagé crée un point de contention.
Enfin, la loi d'Amdahl s'applique car la lecture, le découpage et la sommation
finale restent séquentiels. Appliquée au speedup maximum (1.67×), elle indique
que 40 % du travail est parallélisable et 60 % est intrinsèquement séquentiel.

Pour améliorer la linéarité, `CountWordsConcurrentPool` remplace la création
d'une goroutine par segment par un pool de `numWorkers` goroutines persistantes
qui traitent les segments via un canal de jobs.

**Tableau 3 – Worker pool vs goroutine-par-segment**

| Workers | 1 | 2 | 4 | 8 | 16 | 32 |
|:---:|:---:|:---:|:---:|:---:|:---:|:---:|
| Pool (ms) | 2.86 | 2.08 | 1.71 | 1.59 | **1.46** | 1.53 |
| Par-segment (ms) | 10.39 | 9.02 | 8.47 | 6.93 | 6.31 | 6.21 |
| Gain pool | 3.6× | 4.3× | 5.0× | 4.4× | 4.3× | 4.1× |

Le pool est 4 à 5× plus rapide car il fixe la taille des segments à 50 000
caractères, évitant la surcharge de l'approche naïve. Le plateau dès 16 workers
persiste car la bande passante mémoire reste le vrai goulot d'étranglement.
En appliquant Amdahl au pool (speedup max 1.96×), 52 % du travail devient
parallélisable contre 40 % pour l'approche naïve.

### Bibliographie
- Manuel INF2007, chapitre 8.
- Documentation Go : https://pkg.go.dev/sync, https://pkg.go.dev/testing
- Calculs détaillés : docs/TN5-Calculs.md