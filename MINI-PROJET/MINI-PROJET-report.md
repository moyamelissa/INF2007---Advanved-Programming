# Mini-Projet — Rapport : Robot d'exploration Web concurrent en Go

**Cours** : INF2007 — Programmation avancée  
**Plateforme** : Windows/amd64, Intel Core i5-10300H @ 2.50 GHz, 8 threads

---

## 1. Description de l'implémentation

### Architecture générale

Le programme est structuré autour de la fonction `CrawlURLs(urls []string, maxGoroutines int)` qui orchestre l'exploration concurrente. J'ai combiné trois primitives de concurrence Go :

```
goroutine par URL ──► canal CrawlResult ──► goroutine principale (mutex + agrégation)
       ▲
   sémaphore (limite concurrence)
```

#### Le sémaphore : pourquoi et comment

Pour limiter le nombre de requêtes HTTP simultanées, j'utilise un **canal buffered comme sémaphore** :

```go
semaphore := make(chan struct{}, maxGoroutines)

go func(targetURL string) {
    defer wg.Done()
    semaphore <- struct{}{}        // acquérir (bloque si plein)
    defer func() { <-semaphore }() // libérer
    crawlURL(targetURL, client, ch)
}(u)
```

Ce pattern est plus idiomatique en Go qu'un worker pool classique. Toutes les goroutines sont lancées immédiatement, mais seules `maxGoroutines` s'exécutent réellement en parallèle. Les autres bloquent sur l'écriture dans le sémaphore jusqu'à ce qu'une place se libère. L'avantage par rapport à un pool de workers est la simplicité : pas de file d'attente explicite à gérer.

#### Le canal de résultats et le mutex

Chaque goroutine envoie un `CrawlResult` (URL, nombre de mots, erreur) sur un canal buffered de taille `len(urls)` :

```go
ch := make(chan CrawlResult, len(urls))
```

La goroutine principale lit les résultats et met à jour la map et le total avec un `sync.Mutex` :

```go
for result := range ch {
    mu.Lock()
    if result.Err != nil {
        errs = append(errs, result.Err)
    } else {
        results[result.URL] = result.WordCount
        totalWords += result.WordCount
    }
    mu.Unlock()
}
```

Le mutex est nécessaire ici car `results` (une map) n'est pas thread-safe en Go. Sans mutex, des écritures concurrentes provoqueraient un `fatal error: concurrent map writes`. En pratique, puisque la goroutine principale est la seule à lire le canal, le mutex pourrait être omis — mais je l'ai conservé pour suivre les exigences de l'énoncé et démontrer son utilisation.

#### Le WaitGroup pour fermer le canal

```go
var wg sync.WaitGroup
for _, u := range urls {
    wg.Add(1)
    go func(targetURL string) {
        defer wg.Done()
        // ...
    }(u)
}
go func() {
    wg.Wait()
    close(ch)  // permet à `range ch` de terminer
}()
```

Sans `wg.Wait()` avant `close(ch)`, le canal serait fermé prématurément, et certains résultats seraient perdus.

### Comptage des mots HTML

La fonction `countWordsHTML` utilise le tokenizer `golang.org/x/net/html` plutôt qu'un simple `strings.Fields` sur le HTML brut. La raison : un comptage naïf compterait les noms de balises et attributs comme des mots. Le tokenizer parcourt les tokens et ne collecte que le texte visible :

```go
case html.TextToken:
    if !skip {
        text := strings.TrimSpace(tokenizer.Token().Data)
        if text != "" {
            words := strings.Fields(text)
            count += len(words)
        }
    }
```

Les balises `<script>`, `<style>` et `<noscript>` sont ignorées via un flag `skip` qui s'active sur le `StartTagToken` et se désactive sur le `EndTagToken` correspondant. Par exemple, pour `<script>var x = "hello world";</script>`, les 5 « mots » du JavaScript ne sont pas comptés.

### Client HTTP

```go
func newHTTPClient() *http.Client {
    return &http.Client{Timeout: 10 * time.Second}
}
```

Le timeout de 10 secondes couvre l'intégralité de la requête (DNS, connexion TCP, TLS, envoi, réception). Si un serveur ne répond pas dans ce délai, l'erreur est capturée et renvoyée via le canal.

## 2. Conformité robots.txt

### Implémentation détaillée

Avant chaque exploration, la fonction `checkRobotsAllowed` effectue les étapes suivantes :

```go
func checkRobotsAllowed(targetURL string, client *http.Client) bool {
    parsed, _ := url.Parse(targetURL)
    robotsURL := fmt.Sprintf("%s://%s/robots.txt", parsed.Scheme, parsed.Host)
    resp, err := client.Get(robotsURL)
    // ...
    robots, _ := robotstxt.FromBytes(body)
    group := robots.FindGroup("*")
    return group.Test(parsed.Path)
}
```

Décisions clés :
- **User-agent `*`** : j'utilise le groupe générique car notre robot n'a pas d'identifiant spécifique enregistré.
- **robots.txt inaccessible = autorisé** : si la requête échoue (timeout, 404, erreur réseau), on autorise l'exploration par défaut. C'est le comportement standard (RFC 9309) : l'absence de robots.txt signifie « tout est permis ».
- **Vérification par URL** : `checkRobotsAllowed` est appelé dans chaque goroutine, avant `fetchPage`. Si le chemin est interdit (ex: `Disallow: /private`), la goroutine retourne immédiatement une erreur sans effectuer la requête sur la page.
- **Bibliothèque `github.com/temoto/robotstxt`** : parse correctement les directives `Allow`, `Disallow`, et les wildcards. Évite de réimplémenter un parser de robots.txt.

### Exemple concret testé

Le test `TestCheckRobotsAllowed` crée un serveur avec ce robots.txt :
```
User-agent: *
Disallow: /private/
Allow: /public/
```

Et vérifie que :
- `/public/page` → autorisé ✓
- `/private/secret` → interdit ✓
- `/` → autorisé ✓

## 3. Cas de test

J'ai écrit **14 tests** utilisant `httptest.NewServer` pour simuler des serveurs HTTP sans requêtes réseau réelles. Voici les plus significatifs :

| Test | Ce qu'il vérifie concrètement |
|------|-------------------------------|
| `TestCountWordsHTMLIgnoreScript` | `<p>Texte visible</p><script>var x = "ignoré";</script><p>Autre texte</p>` → 4 mots, pas 8. Prouve que le parsing HTML fonctionne. |
| `TestFetchPageTimeout` | Serveur qui dort 5s + client avec timeout 1s → erreur capturée. Vérifie que les serveurs lents ne bloquent pas le crawler indéfiniment. |
| `TestCrawlURLsIntegration` | 2 pages crawlées en parallèle → `page1: 3 mots`, `page2: 2 mots`, `total: 5`. Test de bout en bout. |
| `TestCrawlURLsRobotsBlocked` | URL interdite par `Disallow` → erreur, 0 mots comptés, map vide. Prouve que robots.txt est respecté avant le fetch. |

Les 6 tests de comptage HTML couvrent : HTML simple, balises multiples, script ignoré, style ignoré, HTML vide, et HTML sans texte.

## 4. Résultats des benchmarks

### Performance selon le nombre de goroutines (8 URLs, serveur local)

| Goroutines | ns/op     | B/op    | allocs/op | Speedup vs 1 goroutine |
|:----------:|:---------:|:-------:|:---------:|:----------------------:|
| 1          | 2 400 628 | 221 971 | 1 679     | 1.00×                  |
| 2          | 1 472 996 | 222 324 | 1 680     | 1.63×                  |
| 4          | 1 697 831 | 316 621 | 2 143     | 1.41×                  |
| 8          | 2 160 625 | 360 811 | 2 370     | 1.11×                  |

### Parsing HTML seul (1 900 mots)

- `BenchmarkCountWordsHTML` : **188 349 ns/op**, 49 144 B/op, 204 allocs/op
- Soit environ **99 ns par mot** pour le parsing HTML.

### Interprétation détaillée

Le **meilleur speedup (1.63×) est obtenu avec 2 goroutines**. Pourquoi la performance se dégrade au-delà ?

Le serveur `httptest.Server` est **local** (loopback `127.0.0.1`). La latence réseau est de ~0.01 ms au lieu de 50-200 ms pour un serveur distant. Dans ces conditions :
- Le goulot d'étranglement est le **listener unique** du serveur de test, qui sérialise les connexions entrantes.
- Plus de goroutines = plus de contention sur ce listener + plus d'allocations mémoire (stacks des goroutines, buffers HTTP).

**En conditions réelles** avec des serveurs distants (latence 100 ms), 8 goroutines traiteraient 8 URLs en ~100 ms au lieu de ~800 ms séquentiellement, soit un speedup théorique de **~8×**. La concurrence est donc critique pour un crawler réel, même si les benchmarks locaux ne le montrent pas.

### Allocations mémoire

On observe une corrélation entre goroutines et allocations :
- 1 goroutine : 1 679 allocs (baseline)
- 8 goroutines : 2 370 allocs (+41 %)

Le surcoût provient des stacks des goroutines (~2 Ko chacune), des buffers du sémaphore, et des structs `CrawlResult` envoyées sur le canal.

## 5. Défis rencontrés et optimisations

1. **Comptage de mots dans du HTML** : un `strings.Fields(html)` brut comptait les balises comme des mots. L'utilisation du tokenizer `x/net/html` a résolu le problème, au prix de 204 allocations par page (pour les tokens et sous-chaînes).

2. **Gestion des erreurs non-bloquante** : j'ai choisi de collecter les erreurs dans un slice plutôt que d'arrêter à la première erreur. Ainsi, si 1 URL sur 10 échoue, on obtient quand même les résultats des 9 autres. Le champ `Err` dans `CrawlResult` permet de distinguer succès et échec par URL.

3. **Dimensionnement du canal** : `make(chan CrawlResult, len(urls))` garantit qu'aucune goroutine ne bloque en écriture, même si la goroutine principale est lente à consommer. Sans le buffer, une goroutine rapide devrait attendre que la goroutine principale lise son résultat avant de libérer le sémaphore, créant un goulot d'étranglement.

4. **Sémaphore vs worker pool** : j'ai préféré le pattern sémaphore car il est plus simple (pas de file d'attente à gérer) et plus idiomatique en Go. Un worker pool serait justifié si les tâches avaient des durées très inégales, ce qui n'est pas le cas ici.

## 6. Conclusion

Le robot d'exploration démontre une utilisation combinée de quatre primitives de concurrence Go : goroutines, canaux, mutex et WaitGroup. Le sémaphore contrôle le parallélisme, le canal transmet les résultats de manière type-safe, le mutex protège la map partagée, et le WaitGroup assure la fermeture propre du canal. Le respect de robots.txt est intégré comme première étape de chaque goroutine, avant toute requête sur la page cible. Les 14 tests unitaires, tous basés sur des serveurs locaux (`httptest`), garantissent la correction sans dépendance réseau externe.
