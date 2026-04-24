# Mini-Projet — Rapport : Robot d'exploration Web concurrent en Go

**Cours** : INF2007 — Programmation avancée  
**Plateforme** : Windows/amd64, Intel Core i5-10300H @ 2.50 GHz, 8 threads

---

## 1. Description de l'implémentation

### Architecture générale

Le programme est structuré autour de la fonction `CrawlURLs(urls []string, maxGoroutines int)` qui orchestre l'exploration concurrente. L'architecture suit le modèle **fan-out / fan-in** avec contrôle de concurrence par sémaphore :

1. **Fan-out** : pour chaque URL, une goroutine est lancée. Un **sémaphore** (canal buffered de taille `maxGoroutines`) limite le nombre de goroutines actives simultanément.
2. **Canal de résultats** : chaque goroutine envoie un `CrawlResult` (URL, nombre de mots, erreur) sur un canal buffered.
3. **Fan-in avec mutex** : la goroutine principale lit les résultats depuis le canal et utilise un `sync.Mutex` pour mettre à jour la map des résultats et le total global de manière sécurisée.
4. **WaitGroup** : un `sync.WaitGroup` assure que le canal est fermé uniquement après la terminaison de toutes les goroutines.

### Comptage des mots HTML

La fonction `countWordsHTML` utilise le tokenizer `golang.org/x/net/html` pour un parsing robuste. Elle parcourt les tokens HTML et ne compte que le texte visible, en ignorant le contenu des balises `<script>`, `<style>` et `<noscript>`. Les mots sont extraits avec `strings.Fields` qui gère correctement les espaces multiples, tabulations et retours à la ligne.

### Client HTTP

Un `http.Client` avec un **timeout de 10 secondes** est utilisé pour toutes les requêtes, conformément aux exigences. Ce timeout couvre la connexion, les redirections et la lecture du corps de la réponse.

## 2. Conformité robots.txt

### Implémentation

Avant chaque exploration, la fonction `checkRobotsAllowed` :

1. Parse l'URL cible pour extraire le schéma et l'hôte.
2. Récupère `<scheme>://<host>/robots.txt` via une requête HTTP.
3. Parse le contenu avec la bibliothèque `github.com/temoto/robotstxt`.
4. Vérifie les directives pour `User-agent: *` en appelant `group.Test(path)`.
5. Si robots.txt n'existe pas (HTTP 404) ou est inaccessible, l'exploration est autorisée par défaut (comportement standard).

Si le chemin est interdit par `Disallow`, la goroutine retourne une erreur et la page n'est pas récupérée.

### Tests de conformité

- `TestCheckRobotsAllowed` : vérifie que `/public/` est autorisé et `/private/` est interdit.
- `TestCheckRobotsNoFile` : vérifie l'autorisation par défaut quand robots.txt n'existe pas.
- `TestCrawlURLsRobotsBlocked` : vérifie le crawl complet avec une URL bloquée.

## 3. Cas de test

| Test | Type | Description |
|------|------|-------------|
| `TestCountWordsHTMLSimple` | Comptage | HTML simple avec 3 mots |
| `TestCountWordsHTMLMultipleTags` | Comptage | Plusieurs balises `<h1>`, `<p>` → 10 mots |
| `TestCountWordsHTMLIgnoreScript` | Comptage | Contenu `<script>` ignoré |
| `TestCountWordsHTMLIgnoreStyle` | Comptage | Contenu `<style>` ignoré |
| `TestCountWordsHTMLEmpty` | Comptage | HTML vide → 0 mots |
| `TestCountWordsHTMLOnlyTags` | Comptage | Balises sans texte → 0 mots |
| `TestFetchPageSuccess` | Réseau | Récupération réussie via serveur de test |
| `TestFetchPageInvalidURL` | Erreur | URL inexistante → erreur |
| `TestFetchPageTimeout` | Erreur | Serveur lent → timeout après 1s |
| `TestFetchPage404` | Erreur | HTTP 404 → erreur |
| `TestCheckRobotsAllowed` | Éthique | Vérification Allow/Disallow |
| `TestCheckRobotsNoFile` | Éthique | Pas de robots.txt → autorisé |
| `TestCrawlURLsIntegration` | Intégration | Crawl complet de 2 pages → total correct |
| `TestCrawlURLsRobotsBlocked` | Intégration | URL bloquée → erreur, 0 mots |

Tous les tests utilisent `httptest.NewServer` pour éviter les requêtes réseau réelles.

## 4. Résultats des benchmarks

### Performance selon le nombre de goroutines (8 URLs, serveur local)

| Goroutines | ns/op     | B/op    | allocs/op | Speedup vs 1 goroutine |
|:----------:|:---------:|:-------:|:---------:|:----------------------:|
| 1          | 2 400 628 | 221 971 | 1 679     | 1.00×                  |
| 2          | 1 472 996 | 222 324 | 1 680     | 1.63×                  |
| 4          | 1 697 831 | 316 621 | 2 143     | 1.41×                  |
| 8          | 2 160 625 | 360 811 | 2 370     | 1.11×                  |

### Parsing HTML (1 900 mots)

- `BenchmarkCountWordsHTML` : **188 349 ns/op**, 49 144 B/op, 204 allocs/op

### Analyse

Le **meilleur speedup (1.63×) est obtenu avec 2 goroutines**. Au-delà, la performance se dégrade car :

- Le serveur de test (`httptest.Server`) est local : la latence réseau est quasi nulle, donc le goulot d'étranglement est le serveur lui-même (un seul listener).
- Plus de goroutines = plus de contention sur le sémaphore, le canal et le mutex.
- Les allocations mémoire augmentent avec le nombre de goroutines (stacks, buffers).

**En conditions réelles** (URLs distantes avec latence réseau de 50-200 ms), le speedup serait beaucoup plus significatif avec 4-8 goroutines, car le temps d'attente réseau (I/O-bound) domine le temps CPU.

## 5. Défis rencontrés et optimisations

1. **Coupure de mots dans le HTML** : le tokenizer `x/net/html` résout ce problème en parsant la structure du document plutôt que le texte brut.
2. **Gestion des erreurs sans arrêt** : les erreurs sont collectées dans un slice, permettant au crawl de continuer pour les autres URLs.
3. **Canal buffered** : dimensionné à `len(urls)` pour éviter les blocages des goroutines en écriture.
4. **Sémaphore vs pool** : le pattern sémaphore (canal buffered) est plus idiomatique en Go qu'un worker pool pour ce cas d'utilisation.

## 6. Conclusion

Le robot d'exploration démontre une utilisation correcte des primitives de concurrence Go (goroutines, canaux, mutex, WaitGroup) avec un respect strict de robots.txt. Les benchmarks confirment que la concurrence apporte un gain réel pour les tâches I/O-bound, avec un point optimal dépendant de la latence réseau et du nombre de cœurs disponibles.
