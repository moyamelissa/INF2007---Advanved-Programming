# Robot d'exploration Web concurrent en Go

**INF2007 – Programmation Avancée · TN6 · Melissa Moya**

---

## 1. Implémentation

Le programme accepte une liste d'URL, lance une goroutine par URL et collecte
les résultats via un canal bufferisé (`chan CrawlResult`). Trois primitives de
concurrence coopèrent pour garantir une exécution correcte et efficace.

Le **sémaphore** (`make(chan struct{}, maxGoroutines)`) limite le nombre de
goroutines actives simultanément. Contrairement à un worker pool classique où
un nombre fixe de workers consomment une file de tâches, le sémaphore permet
de lancer toutes les goroutines immédiatement tout en bloquant celles qui
dépassent la capacité. L'acquisition se fait par `semaphore <- struct{}{}` et
la libération par `<-semaphore` dans un `defer`, ce qui garantit la libération
même en cas de panique.

Le **canal bufferisé** (`make(chan CrawlResult, len(urls))`) découple la
production des résultats de leur consommation. La capacité égale au nombre
d'URLs garantit qu'aucune goroutine ne bloque en envoyant son résultat, même
si la goroutine principale n'a pas encore commencé à lire. Une goroutine
dédiée ferme le canal après que `sync.WaitGroup` confirme la fin de toutes les
explorations, ce qui permet à la goroutine principale d'itérer proprement avec
`for range` sans savoir à l'avance combien de résultats attendre.

L'**agrégation des résultats** repose sur le motif consommateur unique (*single
consumer*), identique à l'exemple `countBytes` du chapitre 8 du manuel : les
goroutines productrices envoient leurs résultats sur le canal `ch`, et une seule
boucle `for range ch` — dans la goroutine principale — lit ces résultats et écrit
dans la carte `results` et dans `totalWords`. L'absence de mutex sur ces structures
n'est pas un oubli : c'est une conséquence directe de l'architecture. Le canal
constitue la frontière de synchronisation entre les producteurs et le consommateur
unique ; une fois un résultat reçu depuis le canal, une seule goroutine y accède,
rendant tout verrou structurellement superflu.

L'énoncé préconisait un mutex pour l'agrégation — approche correcte, mais qui
impose un point de contention : chaque goroutine productrice doit acquérir le
verrou avant d'écrire, ce qui sérialise les mises à jour. Avec le consommateur
unique, les goroutines déposent leurs résultats sur le canal sans se bloquer
mutuellement ; la sérialisation est naturellement assurée par l'ordre de réception.
Cette décision est délibérée : elle élimine toute condition de course sans
ajouter de primitive de synchronisation supplémentaire.

Un `sync.RWMutex` est néanmoins utilisé dans le code, mais pour une raison
distincte : protéger le cache `robotsCache`, qui est lu et écrit simultanément
par plusieurs goroutines productrices. Ce cas ne peut pas être traité par le
motif consommateur unique sans introduire une goroutine gardienne dédiée pour
le cache — d'où l'utilisation justifiée d'un vrai verrou.

Le comptage de mots utilise le tokeniseur `golang.org/x/net/html`. Il parcourt
les jetons HTML et ignore le contenu des balises `<script>`, `<style>` et
`<noscript>` via un drapeau `skip`. Cette approche est plus robuste qu'une
expression régulière car elle gère correctement les balises imbriquées, les
entités HTML et le HTML malformé. Le tokeniseur décode automatiquement les
entités telles que `&amp;` et `&nbsp;`, ce qui évite de les compter comme
des mots.

## 2. Respect de robots.txt

La fonction `checkRobotsAllowed` récupère `robots.txt` à la racine de chaque
domaine avant toute exploration. Elle utilise `github.com/temoto/robotstxt`
pour analyser les directives `User-agent: *` et appliquer les règles `Allow`
et `Disallow`. Cinq cas de robustesse sont couverts et testés, comme le détaille
le tableau suivant.

**Tableau 1 – Cas de robustesse de la vérification robots.txt**

| Situation | Comportement | Justification |
|:---|:---|:---|
| `robots.txt` absent (HTTP 404) | Autoriser | Standard RFC 9309 |
| Serveur injoignable (erreur réseau) | Autoriser | Indisponibilité temporaire |
| Corps tronqué (connexion interrompue) | Autoriser | Contenu partiel non fiable |
| URL non analysable (octet nul) | Interdire | URL malformée |
| Règle `Disallow` explicite | Interdire | Respect de la directive |

Un délai de politesse de 100 ms est appliqué après chaque vérification pour
limiter la fréquence des requêtes et éviter de surcharger les serveurs,
conformément à la norme RFC 9309.

## 3. Tests unitaires

Les 27 tests utilisent `httptest.NewServer` pour simuler des serveurs HTTP
locaux, sans aucun appel réseau réel. Cette approche garantit la
reproductibilité et l'isolation complète des tests, conformément aux principes
du chapitre 6. La technique `net.Listen` directe est utilisée dans
`TestCheckRobotsReadBodyError` pour simuler une connexion TCP interrompue avant
la fin du corps, un cas impossible à reproduire avec `httptest.NewServer`.

**Tableau 2 – Catégories de tests et couverture**

| Catégorie | Nombre | Exemples |
|:---|:---:|:---|
| Comptage HTML | 7 | Ignorer `<script>`, `<style>`, `<noscript>` |
| Récupération de pages | 4 | Succès, URL invalide, timeout, HTTP 404 |
| Vérification robots.txt | 6 | Allow/Disallow, absent, injoignable, corps tronqué |
| Exploration complète | 3 | Intégration, URL bloquée, maxGoroutines=0 |
| Fonctions run/main | 4 | Succès, erreurs, résultats mixtes, point d'entrée |
| Cas limites | 3 | Octet nul, connexion interrompue, HTTP 500 |

La couverture de code atteint 100 % sur les 8 fonctions, y compris `main()`
grâce à l'injection de la variable `mainURLs` dans `TestMainFunction`, et
`checkRobotsAllowed` grâce à `TestCheckRobotsInvalidURLParse` qui force
`url.Parse` à échouer avec un octet nul dans l'URL.

## 4. Résultats des bancs d'essai

Les mesures ont été effectuées sur un Intel i5-10300H à 2,50 GHz (4 cœurs
physiques et 8 threads logiques, Windows/amd64), avec 8 URLs par configuration
et 6 exécutions analysées avec `benchstat`. Le délai de politesse est désactivé
(`politenessDelay = 0`) pour isoler le coût réel de l'exploration.

**Tableau 3 – Temps selon le nombre de goroutines, serveur unique partagé**

| Goroutines | 1 | 2 | 4 | 8 |
|:---:|:---:|:---:|:---:|:---:|
| Temps (ms) | 2.87 | 1.67 | 5.91 | 11.53 |
| Accélération | 1.00× | 1.72× | 0.49× | 0.25× |

Le tableau 3 révèle une régression à partir de 4 goroutines. Quand plusieurs
goroutines accèdent simultanément au même serveur (8 URLs × 2 requêtes par URL
= 16 requêtes concurrentes), le handler de test devient le goulot
d'étranglement. L'intervalle de confiance de ±87 % à 4 goroutines confirme un
comportement bimodal. Les premières itérations s'exécutent en environ 1.9 ms,
puis la contention s'installe et les suivantes dépassent 10 ms.

**Tableau 4 – Temps selon le nombre de goroutines, serveurs distincts par URL**

| Goroutines | 1 | 2 | 4 | 8 |
|:---:|:---:|:---:|:---:|:---:|
| Temps (ms) | 3.02 | 1.52 | 0.92 | 0.86 |
| Accélération | 1.00× | 1.98× | 3.28× | 3.52× |

Le tableau 4 reflète le cas réel où chaque URL provient d'un serveur distinct.
Le gain atteint 3.52× à 8 goroutines, puis plafonne. Sur 4 cœurs physiques,
la loi d'Amdahl limite le gain car la récupération HTTP, l'analyse HTML et
l'agrégation finale comportent des parties séquentielles incompressibles.
En appliquant la formule d'Amdahl au speedup maximum observé (3.52×), on déduit
qu'environ 72 % du travail est parallélisable et 28 % est intrinsèquement
séquentiel.

**Tableau 5 – Performance de l'analyse HTML isolée**

| Métrique | Valeur |
|:---|:---:|
| Temps par analyse (~1 900 mots) | 173 µs |
| Mémoire allouée | 48 Ko |
| Allocations | 204 |
| Intervalle de confiance | ±8 % |

L'analyse HTML de ~1 900 mots prend 173 µs, soit une fraction négligeable du
temps total d'exploration (2.87 ms à 1 goroutine). Le vrai goulot
d'étranglement est le réseau et non le traitement local, ce qui justifie
l'utilité du parallélisme pour ce type de tâche.

## 5. Défis et optimisations

Le premier défi était de rendre les bancs d'essai représentatifs. Le délai de
politesse de 100 ms, nécessaire en production, rendait les mesures inutilisables
car `b.N` descendait à 1 ou 2 itérations et mesurait essentiellement
`time.Sleep`. La solution a été d'extraire le délai dans une variable
`politenessDelay` modifiable, identique au pattern `exitFunc` utilisé au TN5
et `mainURLs` dans ce même projet, ce qui permet de le désactiver dans les
bancs d'essai tout en le conservant actif en production.

Le second défi était l'interprétation du banc d'essai à serveur unique. La
régression observée aurait pu être confondue avec un défaut d'implémentation.
L'ajout de `BenchmarkCrawlGoroutinesMultiServer` avec 8 serveurs distincts a
permis d'isoler la cause, soit la contention côté serveur et non un problème
dans le code du robot d'exploration. Cette distinction est essentielle en
production, où chaque URL provient d'un domaine différent et où le parallélisme
apporte un gain réel de 3.52×.

### Bibliographie

- Manuel INF2007, chapitres 1, 6 et 8.
- Documentation Go : https://pkg.go.dev/net/http, https://pkg.go.dev/sync
- RFC 9309 — Protocole d'exclusion des robots : https://www.rfc-editor.org/rfc/rfc9309
- Bibliothèque robotstxt : https://github.com/temoto/robotstxt
- Bibliothèque HTML : https://pkg.go.dev/golang.org/x/net/html