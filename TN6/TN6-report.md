# Robot d'exploration Web concurrent en Go

**INF2007 – Programmation Avancée · TN6 · Melissa Moya**

---

## 1. Implémentation

Le programme accepte une liste d'URL, lance une goroutine par URL et collecte
les résultats via un canal bufferisé (`chan CrawlResult`). Un sémaphore
(`make(chan struct{}, maxGoroutines)`) limite le nombre de goroutines actives
simultanément à 1, 2, 4 ou 8. Un mutex (`sync.Mutex`) protège l'agrégation
des comptes dans la map de résultats et le total global, conformément aux
exigences de l'énoncé qui demande explicitement les trois primitives.

Chaque goroutine appelle `crawlURL`, qui vérifie d'abord `robots.txt`, applique
un délai de politesse de 100 ms, récupère la page avec `http.Client` (délai
d'expiration de 10 secondes), puis compte les mots visibles. La goroutine
principale itère sur le canal avec `for range` jusqu'à sa fermeture, qui est
déclenchée par une goroutine dédiée après que `sync.WaitGroup` confirme la fin
de toutes les explorations.

Le comptage de mots utilise le tokeniseur `golang.org/x/net/html`. Il parcourt
les jetons HTML et ignore le contenu des balises `<script>`, `<style>` et
`<noscript>` via un drapeau `skip`. Cette approche est plus robuste qu'une
expression régulière car elle gère correctement les balises imbriquées, les
entités HTML et le HTML malformé.

## 2. Respect de robots.txt

La fonction `checkRobotsAllowed` récupère `robots.txt` à la racine de chaque
domaine avant toute exploration. Elle utilise la bibliothèque
`github.com/temoto/robotstxt` pour analyser les directives `User-agent: *` et
appliquer les règles `Allow` et `Disallow`. Si `robots.txt` est absent (HTTP
404), inaccessible (erreur réseau) ou illisible (corps tronqué), la fonction
autorise l'exploration par défaut, conformément au comportement standard des
robots d'exploration. Un délai de politesse de 100 ms est ajouté après chaque
vérification pour limiter la fréquence des requêtes.

## 3. Tests unitaires

Les 27 tests utilisent `httptest.NewServer` pour simuler des serveurs HTTP
locaux, sans aucun appel réseau réel. Cette approche garantit la reproductibilité
et l'isolation des tests.

| Catégorie | Nombre | Exemples |
|:---|:---:|:---|
| Comptage HTML | 7 | Ignorer `<script>`, `<style>`, `<noscript>` |
| Récupération de pages | 4 | Succès, URL invalide, timeout, HTTP 404 |
| Vérification robots.txt | 6 | Allow/Disallow, absent, injoignable, corps tronqué |
| Exploration complète | 3 | Intégration, URL bloquée, maxGoroutines=0 |
| Fonctions run/main | 4 | Succès, erreurs, résultats mixtes, point d'entrée |
| Cas limites | 3 | Octet nul, connexion interrompue, HTTP 500 |

La couverture de code atteint 100 % sur les 8 fonctions, y compris `main()`
grâce à l'injection de la variable `mainURLs` dans `TestMainFunction`.

## 4. Résultats des bancs d'essai

Les mesures ont été effectuées sur un Intel i5-10300H à 2,50 GHz (4 cœurs
physiques et 8 threads logiques, Windows/amd64), 8 URLs par configuration,
6 exécutions analysées avec `benchstat`. Le délai de politesse est désactivé
(`politenessDelay = 0`) pour isoler le coût réel de l'exploration.

**Tableau 1 – Serveur unique partagé (contention visible)**

| Goroutines | 1 | 2 | 4 | 8 |
|:---:|:---:|:---:|:---:|:---:|
| Temps (ms) | 2.87 | 1.67 | 5.91 | 11.53 |
| Accélération | 1.00× | 1.72× | 0.49× | 0.25× |

**Tableau 2 – Serveurs distincts par URL (parallélisme réel)**

| Goroutines | 1 | 2 | 4 | 8 |
|:---:|:---:|:---:|:---:|:---:|
| Temps (ms) | 3.02 | 1.52 | 0.92 | 0.86 |
| Accélération | 1.00× | 1.98× | 3.28× | 3.52× |

Le tableau 1 révèle une régression à partir de 4 goroutines. Quand plusieurs
goroutines accèdent simultanément au même serveur (8 URLs × 2 requêtes par URL
= 16 requêtes concurrentes), le handler `httptest` devient le goulot
d'étranglement. L'intervalle de confiance de ±87 % à 4 goroutines confirme un
comportement bimodal : les premières itérations de `b.N` s'exécutent en ~1.9 ms,
puis la contention s'installe et les suivantes dépassent 10 ms.

Le tableau 2 reflète le cas réel où chaque URL provient d'un serveur distinct.
Le gain atteint 3.52× à 8 goroutines, puis plafonne. Sur 4 cœurs physiques,
au-delà de 8 goroutines, la loi d'Amdahl limite le gain : la récupération HTTP,
le parsing HTML et l'agrégation finale comportent des parties séquentielles
incompressibles.

## 5. Défis et optimisations

Le défi principal était de rendre les bancs d'essai représentatifs. Le délai de
politesse de 100 ms, nécessaire en production, rendait les mesures inutilisables
car `b.N` descendait à 1 ou 2 itérations et mesurait essentiellement
`time.Sleep`. La solution a été d'extraire le délai dans une variable
`politenessDelay` modifiable, identique au pattern `exitFunc` utilisé au TN5,
ce qui permet de le désactiver dans les bancs d'essai tout en le conservant
actif en production.

Un second défi était l'interprétation du benchmark à serveur unique. La
régression observée aurait pu être confondue avec un défaut d'implémentation.
L'ajout de `BenchmarkCrawlGoroutinesMultiServer` avec 8 serveurs distincts a
permis d'isoler la cause : la contention côté serveur, et non un problème dans
le code du robot.

### Bibliographie
- Manuel INF2007, chapitres 1, 6 et 8.
- Documentation Go : https://pkg.go.dev/net/http, https://pkg.go.dev/sync
- Bibliothèque robotstxt : https://github.com/temoto/robotstxt
- Bibliothèque HTML : https://pkg.go.dev/golang.org/x/net/html
