# Travail 4 : Mesurez et optimisez la somme des sinus en Go

## Informations générales

**Semaine de remise** : Semaine 10

**Objectif** : Vous allez mesurer et optimiser les performances d'un programme Go, en appliquant les concepts de benchmarking du chapitre 6.

## Description du travail

Votre mission est de développer un programme Go qui génère un tableau aléatoire d'entiers ou de nombres à virgule flottante dans [0,1), puis calcule la somme des sinus des éléments de ce tableau. Vous analyserez les résultats pour identifier les facteurs qui influencent les performances, en comparant différents types de données et tailles de tableau.

## Consignes

Vous devez écrire un programme Go qui :

1. Génère un tableau de 1 000 000 d'éléments, soit d'entiers aléatoires (par exemple, entre 0 et 1000), soit de nombres à virgule flottante dans [0,1). Acceptez un paramètre via la ligne de commande pour choisir le type de données (par exemple, `go run sinesum.go --type=float` ou `--type=int`).
2. Implémente une fonction `computeSineSum` pour calculer la somme des sinus des éléments du tableau.
3. Écrit des tests de benchmarking avec `testing` pour mesurer le temps d'exécution en traitant 1 %, 10 %, 20 %, 30 %, 40 %, 50 %, 60 %, 70 %, 80 %, 90 % et 100 % du tableau, pour les deux types de données (entiers et flottants).

Vous devez également écrire des tests unitaires pour vérifier :

- La correction de `computeSineSum` pour un petit tableau (par exemple, `[1, 2, 3]` pour les entiers ou `[0.1, 0.2, 0.3]` pour les flottants).
- La gestion des erreurs pour un type de données invalide passé en paramètre.

### Conseils pour générer des tableaux aléatoires

- Utilisez le paquet `math/rand` pour générer des nombres aléatoires. Initialisez le générateur pour des résultats reproductibles (par exemple, `rand.Seed(42)`).
- Pour les entiers, utilisez `rand.Intn(max)` (par exemple, `rand.Intn(1001)` pour des entiers de 0 à 1000).
- Pour les flottants dans [0,1), utilisez `rand.Float64()`, qui génère directement des valeurs dans cet intervalle.
- Créez le tableau avec `make([]int, 1000000)` ou `make([]float64, 1000000)`, puis remplissez-le dans une boucle.

## Ce que vous devez soumettre

1. Votre code source (`sinesum.go`, `sinesum_test.go`).
2. Un rapport de 1 à 2 pages (PDF) analysant les résultats de vos benchmarks, avec des graphiques ou tableaux. Comparez les performances pour les entiers et les flottants, et discutez des facteurs influençant les résultats (par exemple, précision des flottants, taille du tableau).
3. À quelle distance voyage la lumière pendant le calcul d'un seul sinus ? La lumière voyage 299 792 458 mètres par seconde.
4. Vous voulez créer un jeu vidéo qui affiche une nouvelle image 120 fois par secondes. Dans un seul tick (1 seconde divisé par 120), combien de sinus pouvez-vous calculer ?

## Critères d'évaluation

- Correction des implémentations : 30 %
- Qualité des tests de benchmarking : 30 %
- Profondeur de l'analyse des performances dans le rapport : 40 %

## Directives de soumission

- Soumettez vos fichiers via la plateforme en ligne du cours à la fin de la semaine 10.
- Votre code doit être bien documenté et formaté avec `gofmt`.
- Votre rapport doit être en PDF, avec une police de 12 points, et inclure votre nom et le titre du travail.
