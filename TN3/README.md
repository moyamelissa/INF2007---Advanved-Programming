# Projet : Analyse de capteurs IoT (INF2007 – Travail 3)

![Tests](https://github.com/moyamelissa/Advanced-Programming/actions/workflows/tn3-coverage.yml/badge.svg)
[![codecov](https://codecov.io/gh/moyamelissa/Advanced-Programming/branch/main/graph/badge.svg?flag=tn3)](https://codecov.io/gh/moyamelissa/Advanced-Programming)

Ce projet implémente une fonction d'analyse de données binaires de capteurs IoT en Go.
Chaque entrée 32 bits encode un identifiant de capteur (7 bits), un bit de validation et une valeur sur 24 bits.

## Fonction implémentée

| Fonction | Description |
|----------|-------------|
| `Analyse(data, capteur)` | Extrait et valide les mesures d'un capteur donné par manipulation de bits |

### Opérations bit à bit utilisées

| Opération | Rôle |
|-----------|------|
| `(1<<7) - 1` | Masque pour isoler les 7 bits d'identifiant |
| `entry & (1<<7)` | Test du bit de validation (bit 7) |
| `entry >> 8` | Décalage pour extraire la valeur (bits 8–31) |
| `x & (x-1)` | Détection de plus d'un bit actif (sans boucle) |
| `bits.TrailingZeros32` | Position du bit actif en O(1) via instruction CPU `TZCNT` |

## Structure du projet

```
TN3/
├── go.mod                       # Module Go
├── analyse.go                   # Implémentation principale
├── analyse_test.go              # 6 tests unitaires
├── TN3-report.md                # Rapport d'analyse
├── TN3-AI-Prompts.md            # Prompts IA utilisés
├── TN3-Homework-Instruction.md  # Énoncé du travail
└── README.md                    # Ce fichier
```

## Prérequis

- Go 1.21+

## Tests unitaires

```bash
go test -v ./...
```

Avec couverture :

```bash
go test -v -cover ./...
```

## Tests disponibles

| Test | Ce qu'il vérifie |
|------|-----------------|
| `TestAnalyseDonneesValides` | Comptage correct des bits pour des entrées variées |
| `TestAnalyseCapteurInvalide` | Rejet d'un identifiant > 127 |
| `TestAnalyseBit7Invalide` | Rejet d'une entrée avec bit de validation à 1 |
| `TestAnalysePlusieursBitsValeur` | Détection de plusieurs bits actifs via `x & (x-1)` |
| `TestAnalyseExempleEnonce` | Reproduction exacte de l'exemple du professeur |
| `TestAnalyseTableauVide` | Cas limite : tableau vide retourne des zéros sans erreur |

## Liens

- [Rapport TN3](TN3-report.md)
- [Prompts IA](TN3-AI-Prompts.md)
- [Dépôt GitHub](https://github.com/moyamelissa/Advanced-Programming/tree/main/TN3)