# Comptage concurrent de mots en Go

**INF2007 – Programmation Avancée · TN5 · Melissa Moya**

---

## 1. Architecture concurrente

Le programme lit un fichier texte, divise son contenu en segments de taille
paramétrable, puis compte les mots de chaque segment en parallèle. La fonction
`splitIntoSegments` décale chaque coupure vers le prochain espace pour ne jamais
couper un mot. Une goroutine est lancée par segment via `go countWordsInSegment`,
et les résultats sont collectés sur un canal bufferisé (`chan int`) de capacité
égale au nombre de segments, ce qui évite tout blocage. La goroutine principale
itère sur le canal et somme les résultats. L'absence de variable partagée et
l'usage exclusif du canal garantissent une exécution correcte sans verrou. Une
variante `CountWordsConcurrentN` permet de fixer le nombre de goroutines pour
mesurer la linéarité du gain.

## 2. Résultats des benchmarks

Mesures effectuées sur un Intel i5-10300H à 2,50 GHz (4 cœurs physiques /
8 threads logiques, Windows/amd64), fichier de test d'environ 700 000 caractères
(100 000 mots), 6 exécutions par configuration analysées avec `benchstat`.

**Graphique 1 – Scaling réel vs linéaire idéal selon le nombre de goroutines**

![Graphique 1 – Réel vs linéaire idéal](data/worker-count-chart.png)

**Tableau 1 – Temps selon la taille des segments**

| Segment | 500 | 1 000 | 5 000 | 10 000 | 50 000 | 100 000 | Tout |
|:---:|:---:|:---:|:---:|:---:|:---:|:---:|:---:|
| Temps (ms) | 8.87 | 2.84 | 4.84 | 4.85 | **1.78** | 1.97 | 4.09 |

Le comptage séquentiel de référence prend 5.08 ms. Le concurrent atteint un
optimum à 50 000 caractères (1.78 ms), soit un gain de 2.85× sur le séquentiel.
Avec des segments trop petits (500 caractères, environ 1 200 goroutines),
l'overhead de scheduling dépasse le bénéfice du parallélisme et la performance
devient pire que le séquentiel.

**Tableau 2 – Temps selon le nombre de goroutines (`CountWordsConcurrentN`)**

| Workers | 1 | 2 | 4 | 8 | 16 | 32 |
|:---:|:---:|:---:|:---:|:---:|:---:|:---:|
| Temps (ms) | 10.39 | 9.02 | 8.47 | 6.93 | 6.31 | 6.21 |
| Speedup | 1.00× | 1.15× | 1.23× | 1.50× | 1.65× | 1.67× |

## 3. Analyse de la linéarité

La performance ne croît pas linéairement avec le nombre de goroutines. Sur
4 cœurs physiques disponibles, le speedup mesuré n'atteint que 1.67× au lieu
des 4× théoriques, et atteint un plateau dès 16 workers. Quatre facteurs
expliquent cet écart. Premièrement, `strings.Fields` est une opération
séquentielle interne rapide (un seul balayage du segment), donc le travail par
goroutine est trop court pour amortir le coût de création et d'orchestration.
Deuxièmement, le canal partagé crée un point de contention naturel à la collecte
des résultats. Troisièmement, la pression mémoire amplifie cet overhead :
l'allocation passe de 1 alloc/op pour la version séquentielle à 71 pour 32
workers, et jusqu'à 2 700 pour des segments de 500 caractères — chaque goroutine
entraîne des allocations supplémentaires pour sa pile et le canal.

Quatrièmement, la loi d'Amdahl s'applique : la lecture du
fichier, le découpage en segments et la sommation finale restent séquentiels. En
appliquant la formule d'Amdahl au speedup maximum observé (1.67×), on déduit
qu'environ 40 % du travail est parallélisable et 60 % est intrinsèquement
séquentiel. Pour atteindre un scaling linéaire, il faudrait soit une charge CPU
plus lourde par segment (ex. analyse syntaxique au lieu d'un simple comptage),
soit éliminer la contention sur le canal (ex. sommation locale par worker puis
fusion). Le choix optimal pour ce workload reste donc d'utiliser un segment de
50 000 caractères, qui équilibre parallélisme et overhead.

### Bibliographie
- Manuel INF2007, chapitre 8.
- Documentation Go : https://pkg.go.dev/sync, https://pkg.go.dev/testing
