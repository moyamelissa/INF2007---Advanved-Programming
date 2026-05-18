# Robot d'exploration Web concurrent en Go

**INF2007 – Programmation Avancée  TN6  Melissa Moya**

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

Le **mutex** (`sync.Mutex`) protège la mise à jour de la carte de résultats et
du total global. Sans mutex, deux goroutines pourraient lire puis écrire
simultanément dans la carte, provoquant une condition de course. Un canal seul
aurait suffi pour collecter les résultats, mais le mutex est requis explicitement
par l'énoncé et permet d'agréger directement dans la goroutine principale sans
goroutine intermédiaire.

Le comptage de mots utilise le tokeniseur `golang.org/x/net/html`. Il parcourt
les jetons HTML et ignore le contenu des balises `<script>`, `<style>` et
`<noscript>` via un drapeau `skip`. Cette approche est plus robuste qu'une
expression régulière car elle gère correctement les balises imbriquées, les
entités HTML et le HTML malformé. Le tokeniseur décote automatiquement les
entités (`&amp;` → `&`, `&nbsp;` → espace), ce qui évite de les compter comme
des mots.

## 2. Respect de robots.txt

La fonction `checkRobotsAllowed` récupère `robots.txt` à la racine de chaque
domaine avant toute exploration. Elle utilise `github.com/temoto/robotstxt`
pour analyser les directives `User-agent: *` et appliquer les règles `Allow`
et `Disallow`. Cinq cas de robustesse sont couverts et testés :

**Tableau 1 – Cas de robustesse de la vérification robots.txt**

| Situation | Comportement | Justification |
|:---|:---|:---|
| `robots.txt` absent (HTTP 404) | Autoriser | Standard RFC 9309 |
| Serveur injoignable (erreur réseau) | Autoriser | Indisponibilité temporaire |
| Corps tronqué (connexion interrompue) | Autoriser | Contenu partiel non fiable |
| URL non analysable (octet nul) | Interdire | URL malformée |
| Règle `Disallow` explicite | Interdire | Respect de la directive |

Un délai de politesse de 100 ms est appliqué après chaque vérification pour
limiter la fréquence des requêtes et éviter de surcharger les serveurs, comme
recommandé par la norme RFC 9309.

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
physiques et 8 threads logiques, Windows/amd64), 8 URLs par configuration,
6 exécutions analysées avec `benchstat`. Le délai de politesse est désactivé
(`politenessDelay = 0`) pour isoler le coût réel de l'exploration.

**Tableau 3 – Temps selon le nombre de goroutines, serveur unique partagé**

| Goroutines | 1 | 2 | 4 | 8 |
|:---:|:---:|:---:|:---:|:---:|
| Temps (ms) | 2.87 | 1.67 | 5.91 | 11.53 |
| Accélération | 1.00 | 1.72 | 0.49 | 0.25 |

Le tableau 3 révèle une régression à partir de 4 goroutines. Quand plusieurs
goroutines accèdent simultanément au même serveur (8 URLs  2 requêtes par URL
= 16 requêtes concurrentes), le handler de test devient le goulot
d'étranglement. L'intervalle de confiance de 87 % à 4 goroutines confirme un
comportement bimodal : les premières itérations s'exécutent en environ 1.9 ms,
puis la contention s'installe et les suivantes dépassent 10 ms.

**Tableau 4 – Temps selon le nombre de goroutines, serveurs distincts par URL**

| Goroutines | 1 | 2 | 4 | 8 |
|:---:|:---:|:---:|:---:|:---:|
| Temps (ms) | 3.02 | 1.52 | 0.92 | 0.86 |
| Accélération | 1.00 | 1.98 | 3.28 | 3.52 |

Le tableau 4 reflète le cas réel où chaque URL provient d'un serveur distinct.
Le gain atteint 3.52 à 8 goroutines, puis plafonne. Sur 4 cœurs physiques,
la loi d'Amdahl limite le gain : la récupération HTTP, l'analyse HTML et
l'agrégation finale comportent des parties séquentielles incompressibles.
En appliquant la formule d'Amdahl au speedup maximum observé (3.52), on déduit
qu'environ 72 % du travail est parallélisable et 28 % intrinsèquement séquentiel.

**Tableau 5 – Performance de l'analyse HTML isolée**

| Métrique | Valeur |
|:---|:---:|
| Temps par analyse (~1 900 mots) | 173 µs |
| Mémoire allouée | 48 Ko |
| Allocations | 204 |
| Intervalle de confiance | 8 % |

L'analyse HTML de ~1 900 mots prend 173 µs, soit une fraction négligeable du
temps total d'exploration (2.87 ms à 1 goroutine). Le vrai goulot
d'étranglement est le réseau, pas le traitement local — ce qui justifie
l'utilité du parallélisme pour ce type de tâche.

## 5. Défis et optimisations

Le premier défi était de rendre les bancs d'essai représentatifs. Le délai de
politesse de 100 ms, nécessaire en production, rendait les mesures inutilisables
car `b.N` descendait à 1 ou 2 itérations et mesurait essentiellement
`time.Sleep`. La solution a été d'extraire le délai dans une variable
`politenessDelay` modifiable — identique au pattern `exitFunc` utilisé au TN5
et `mainURLs` dans ce même projet — ce qui permet de le désactiver dans les
bancs d'essai tout en le conservant actif en production.

Le second défi était l'interprétation du banc d'essai à serveur unique. La
régression observée aurait pu être confondue avec un défaut d'implémentation.
L'ajout de `BenchmarkCrawlGoroutinesMultiServer` avec 8 serveurs distincts a
permis d'isoler la cause : la contention côté serveur, et non un problème dans
le code du robot d'exploration. Cette distinction est essentielle pour
l'optimisation : en production, où chaque URL provient d'un domaine différent,
le parallélisme apporte un gain réel de 3.52.

### Bibliographie
- Manuel INF2007, chapitres 1, 6 et 8.
- Documentation Go : https://pkg.go.dev/net/http, https://pkg.go.dev/sync
- RFC 9309 — Protocole d'exclusion des robots : https://www.rfc-editor.org/rfc/rfc9309
- Bibliothèque robotstxt : https://github.com/temoto/robotstxt
- Bibliothèque HTML : https://pkg.go.dev/golang.org/x/net/html
