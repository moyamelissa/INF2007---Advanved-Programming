# Projet : Robot d'exploration Web concurrent (INF2007 – Travail 6)

![Tests](https://github.com/moyamelissa/Advanced-Programming/actions/workflows/tn6-coverage.yml/badge.svg)
[![codecov](https://codecov.io/gh/moyamelissa/Advanced-Programming/branch/main/graph/badge.svg?flag=tn6)](https://codecov.io/gh/moyamelissa/Advanced-Programming)

Ce projet implémente un robot d'exploration Web concurrent en Go. Il récupère des pages HTML, compte les mots visibles et respecte les règles de `robots.txt`. La concurrence est gérée par des goroutines, un canal de résultats et un sémaphore (canal buffered).

## Fonctions implémentées

| Fonction | Description |
|----------|-------------|
| `newHTTPClient()` | Crée un client HTTP avec timeout de 10 s |
| `checkRobotsAllowed(url, client)` | Vérifie si `robots.txt` autorise l'exploration |
| `fetchPage(url, client)` | Récupère le contenu HTML d'une URL |
| `countWordsHTML(html)` | Compte les mots visibles (ignore `<script>`, `<style>`, `<noscript>`) |
| `crawlURL(url, client, ch)` | Explore une URL et envoie le résultat sur un canal |
| `CrawlURLs(urls, maxGoroutines)` | Orchestre le crawl concurrent avec sémaphore |

## Structure du projet

```
TN6/
├── go.mod                # Module Go
├── go.sum                # Dépendances vérifiées
├── crawler.go            # Code principal + CLI
├── crawler_test.go       # 14 tests unitaires + benchmarks
├── TN6-report.md         # Rapport d'analyse
├── TN6-AI-Prompts.md     # Prompts IA utilisés
└── README.md             # Ce fichier
```

## Prérequis

- Go 1.21+

## Exécution

```bash
go run crawler.go
```

## Tests unitaires

```bash
go test -v -run="Test" ./...
```

## Benchmarks

```bash
go test -bench="Benchmark" -benchmem -run="^$" -count=1 ./...
```

## Tests disponibles

| Test | Ce qu'il vérifie |
|------|-----------------|
| `TestCountWordsHTMLSimple` | Comptage de 3 mots dans un `<p>` |
| `TestCountWordsHTMLMultipleTags` | Comptage à travers `<h1>` et `<p>` multiples |
| `TestCountWordsHTMLIgnoreScript` | Contenu de `<script>` ignoré |
| `TestCountWordsHTMLIgnoreStyle` | Contenu de `<style>` ignoré |
| `TestCountWordsHTMLEmpty` | HTML vide retourne 0 |
| `TestCountWordsHTMLOnlyTags` | Balises sans texte retourne 0 |
| `TestFetchPageSuccess` | Récupération d'une page via serveur de test |
| `TestFetchPageInvalidURL` | Gestion d'erreur pour URL invalide |
| `TestFetchPageTimeout` | Gestion du délai d'expiration |
| `TestFetchPage404` | Gestion d'un code HTTP 404 |
| `TestCheckRobotsAllowed` | Respect des règles Allow/Disallow |
| `TestCheckRobotsNoFile` | Autorisation par défaut si robots.txt absent |
| `TestCrawlURLsIntegration` | Crawl complet de 2 pages locales |
| `TestCrawlURLsRobotsBlocked` | URL bloquée par robots.txt retourne erreur |

## Benchmarks disponibles

| Benchmark | Description |
|-----------|-------------|
| `BenchmarkCrawlGoroutines` | Compare 1, 2, 4, 8 goroutines sur 8 URLs |
| `BenchmarkCountWordsHTML` | Parsing HTML de ~1 900 mots |

## Liens

- [Rapport TN6](TN6-report.md)
- [Prompts IA](TN6-AI-Prompts.md)
- [Dépôt GitHub](https://github.com/moyamelissa/Advanced-Programming/tree/main/TN6)
