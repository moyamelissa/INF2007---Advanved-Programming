# Projet : Robot d'exploration Web concurrent (INF2007 – Travail 6)

![Tests](https://github.com/moyamelissa/Advanced-Programming/actions/workflows/tn6-coverage.yml/badge.svg)
[![codecov](https://codecov.io/gh/moyamelissa/Advanced-Programming/branch/main/graph/badge.svg?flag=tn6)](https://codecov.io/gh/moyamelissa/Advanced-Programming)

Ce projet implémente un robot d'exploration Web concurrent en Go. Il récupère
des pages HTML, compte les mots visibles et respecte les règles de `robots.txt`.
La concurrence est gérée par des goroutines, un canal de résultats et un sémaphore
(canal bufferisé).

## Fonctions implémentées

| Fonction | Description |
|----------|-------------|
| `newHTTPClient()` | Crée un client HTTP avec délai d'expiration de 10 s |
| `checkRobotsAllowed(url, client)` | Vérifie si `robots.txt` autorise l'exploration |
| `fetchPage(url, client)` | Récupère le contenu HTML d'une URL |
| `countWordsHTML(html)` | Compte les mots visibles (ignore `<script>`, `<style>`, `<noscript>`) |
| `crawlURL(url, client, ch)` | Explore une URL et envoie le résultat sur un canal |
| `CrawlURLs(urls, maxGoroutines)` | Orchestre l'exploration concurrente avec sémaphore |
| `politenessDelay` | Délai de politesse configurable (100 ms par défaut) |
| `mainURLs` | Liste d'URLs injectable pour les tests |

## Structure du projet

```
.
├── go.mod                # Module Go
├── go.sum                # Dépendances vérifiées
├── crawler.go            # Code principal
├── crawler_test.go       # 28 tests unitaires + 3 bancs d'essai
├── README.md             # Ce fichier
├── TN6-report.md         # Rapport (copie racine)
├── chart/                        # Générateur de graphique de benchmarks
│   ├── go.mod                    # Module Go du sous-projet chart
│   └── main.go                   # Génère data/benchmark-chart.png
├── data/                         # Résultats de tests et benchmarks
│   ├── benchmark-chart.png       # Graphique des benchmarks
│   ├── bench-raw.txt             # Données brutes des benchmarks
│   ├── benchstat.txt             # Statistiques benchstat
│   ├── coverage.txt              # Rapport de couverture de code
│   └── tests.txt                 # Résultats des tests unitaires
└── docs/                         # Documentation du projet
    ├── TN6-report.md             # Rapport d'analyse
    ├── TN6-AI-Prompts.md         # Prompts IA utilisés
    └── TN6-Homework-Instructions.md  # Énoncé du travail
```

## Prérequis

- Go 1.21+

## Exécution

```bash
go run .
```

## Tests unitaires

```bash
go test -v -run="Test" ./...
```

## Bancs d'essai

```bash
go test -bench="Benchmark" -benchmem -run="^$" -count=1 ./...
```

## Tests disponibles

| Catégorie | Nombre | Tests |
|:---|:---:|:---|
| Comptage HTML | 7 | `Simple`, `MultipleTags`, `IgnoreScript`, `IgnoreStyle`, `Empty`, `OnlyTags`, `Noscript` |
| Récupération de pages | 5 | `Success`, `InvalidURL`, `Timeout`, `404`, `ReadError` |
| Vérification robots.txt | 8 | `Allowed`, `NoFile`, `InvalidURL`, `Unreachable`, `InvalidBody`, `ReadBodyError`, `InvalidURLParse`, `FetchRobotsParseError` |
| Exploration complète | 3 | `Integration`, `RobotsBlocked`, `ZeroGoroutines` |
| crawlURL | 1 | `FetchError` |
| Fonctions run/main | 4 | `RunFunction`, `RunFunctionWithErrors`, `RunFunctionMixedResults`, `MainFunction` |

### Détail des tests

| Test | Ce qu'il vérifie |
|------|-----------------|
| `TestCountWordsHTMLSimple` | Comptage de 3 mots dans un `<p>` |
| `TestCountWordsHTMLMultipleTags` | Comptage à travers `<h1>` et `<p>` multiples |
| `TestCountWordsHTMLIgnoreScript` | Contenu de `<script>` ignoré |
| `TestCountWordsHTMLIgnoreStyle` | Contenu de `<style>` ignoré |
| `TestCountWordsHTMLEmpty` | HTML vide retourne 0 |
| `TestCountWordsHTMLOnlyTags` | Balises sans texte retourne 0 |
| `TestCountWordsHTMLNoscript` | Contenu de `<noscript>` ignoré |
| `TestFetchPageSuccess` | Récupération d'une page via serveur de test |
| `TestFetchPageInvalidURL` | Gestion d'erreur pour URL invalide |
| `TestFetchPageTimeout` | Gestion du délai d'expiration |
| `TestFetchPage404` | Gestion d'un code HTTP 404 |
| `TestFetchPageReadError` | Gestion d'erreur si connexion interrompue |
| `TestCheckRobotsAllowed` | Respect des règles Allow/Disallow |
| `TestCheckRobotsNoFile` | Autorisation par défaut si robots.txt absent |
| `TestCheckRobotsInvalidURL` | Retourne false pour URL non parseable |
| `TestCheckRobotsUnreachable` | Autorise si le serveur est injoignable |
| `TestCheckRobotsInvalidBody` | Corps valide de robots.txt parsé sans erreur |
| `TestCheckRobotsReadBodyError` | Autorise si la lecture du corps robots.txt échoue |
| `TestCheckRobotsInvalidURLParse` | Retourne false pour URL avec octet nul |
| `TestFetchRobotsParseError` | Retourne nil quand `robotstxt.FromBytes` échoue (`Disallow` avant `User-agent`) |
| `TestCrawlURLsIntegration` | Exploration complète de 2 pages locales |
| `TestCrawlURLsRobotsBlocked` | URL bloquée par robots.txt retourne une erreur |
| `TestCrawlURLsZeroGoroutines` | maxGoroutines ≤ 0 traité comme 1 |
| `TestCrawlURLFetchError` | Erreur HTTP 500 capturée dans CrawlResult |
| `TestRunFunction` | `run` s'exécute sans panique avec serveur local |
| `TestRunFunctionWithErrors` | `run` gère les erreurs sans panique |
| `TestRunFunctionMixedResults` | `run` gère résultats et erreurs simultanément |
| `TestMainFunction` | `main()` couvre le point d'entrée sans appel réseau |

## Bancs d'essai disponibles

| Banc d'essai | Description |
|--------------|-------------|
| `BenchmarkCrawlGoroutines` | Compare 1, 2, 4, 8 goroutines sur 8 URLs via serveur unique |
| `BenchmarkCrawlGoroutinesMultiServer` | Compare 1, 2, 4, 8 goroutines sur 8 serveurs distincts |
| `BenchmarkCountWordsHTML` | Analyse HTML de ~1 900 mots avec balises ignorées |

## Résultats clés

| Mesure | Valeur |
|--------|--------|
| Couverture de code | 100 % |
| Accélération optimale (multi-serveurs) | 3.52× (8 goroutines) |
| Temps d'analyse HTML (~1 900 mots) | ~173 µs |
| Délai de politesse | 100 ms par URL |

## Liens

- [Rapport TN6](docs/TN6-report.md)
- [Prompts IA](docs/TN6-AI-Prompts.md)
- [Dépôt GitHub](https://github.com/moyamelissa/Advanced-Programming/tree/main/TN6)
