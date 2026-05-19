# Travail 5 : Comptez les mots dans un fichier avec la concurrence en Go

## Informations générales

**Semaine de remise** : Semaine 12

**Objectif** : Appliquer les concepts de programmation concurrente du Chapitre 8 pour créer un programme Go qui compte les mots dans un fichier texte en utilisant des goroutines et des canaux.

## Description du travail

Votre tâche est de développer un programme Go qui lit un fichier texte passé en paramètre, divise son contenu en segments, et compte de manière concurrente le nombre total de mots. Vous utiliserez des goroutines pour traiter chaque segment en parallèle et des canaux pour collecter les résultats, en veillant à une exécution correcte et efficace. Un mot est une séquence de caractères séparée par des espaces.

## Consignes

Vous devez écrire un programme Go qui :

1. Accepte un chemin de fichier texte en entrée via les arguments de la ligne de commande (par exemple, `go run wordcount.go input.txt`).
2. Lit le contenu du fichier et le divise en segments. La taille des segments doit être un paramètre en ligne de commande. Prenez en compte qu'un segment comprenant un nombre donné de caractères peut se terminer au milieu d'un mot. Vous devrez en tenir compte.
3. Lance une goroutine pour chaque segment afin de compter ses mots.
4. Utilise un canal pour collecter le nombre de mots calculé par chaque goroutine.
5. Somme les résultats dans la goroutine principale et affiche le nombre total de mots.

Vous devez également écrire des tests unitaires pour vérifier la logique de comptage des mots, en couvrant au moins trois cas :

- Fichier vide.
- Fichier avec une seule ligne contenant un mot (par exemple, `"Bonjour"`).
- Fichier avec plusieurs lignes contenant plusieurs mots.

Vous devez inclure des mesures de performance. Comment est-ce que la performance évolue en fonction de la taille des segments ? Écrivez des benchmarks. Est-ce que la performance croît linéairement avec le nombre de goroutines ? Essayez de faire en sorte que ça soit le cas. Est-ce que vous y arrivez ?

## Ce que vous devez soumettre

1. Votre code source (`wordcount.go`, `wordcount_test.go`).
2. Un rapport d'une page (PDF) expliquant comment les goroutines et les canaux garantissent une exécution concurrente correcte.

## Critères d'évaluation

- Correction de l'implémentation concurrente : 40 %
- Qualité des tests unitaires : 20 %
- Utilisation appropriée des goroutines et des canaux : 20 %
- Clarté du rapport : 20 %

## Directives de soumission

- Soumettez vos fichiers via la plateforme en ligne avant la fin de la semaine 12.
- Votre code doit être bien documenté et formaté avec `gofmt`.
- Votre rapport doit être en PDF, avec une police de 12 points, et inclure votre nom et le titre du travail.
- Vous pouvez remettre une vidéo avec votre travail.
