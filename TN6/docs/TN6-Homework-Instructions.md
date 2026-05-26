# TN6 – Énoncé du travail (Mini-Projet)

## Votre projet : Créez un robot d'exploration Web concurrent

### Informations Générales

**Date de remise :** Fin de la semaine 15.

**Objectif :** Vous allez mettre en pratique les concepts des chapitres 1, 6 et 8 pour développer une application Go robuste, en combinant tests, optimisation des performances et programmation concurrente.

### Description du Projet

Votre mission est de créer un robot d'exploration web concurrent qui récupère le contenu de plusieurs URL, compte les mots sur chaque page et calcule le total des mots pour toutes les pages. Vous utiliserez des goroutines pour gérer les requêtes en parallèle, des canaux pour transmettre les résultats et des mutex pour agréger les comptes de manière sécurisée. Ce projet vous permettra de démontrer vos compétences en programmation Go tout en respectant des pratiques éthiques.

### Consignes

Vous devez développer un programme Go qui :

- Accepte une liste d'URL en entrée (par exemple, `[]string{"https://www.google.com", "https://www.github.com"}`).
- Lance une goroutine pour chaque URL pour récupérer son contenu avec `http.Client`, en définissant un délai d'expiration de 10 secondes (voir Chapitre 8 du manuel).
- Analyse le contenu HTML pour compter les mots. Vous pouvez utiliser un comptage simple basé sur les espaces ou la bibliothèque `golang.org/x/net/html`.
- Utilise un canal pour collecter les comptes de mots et un mutex pour mettre à jour le total de manière sécurisée.
- Affiche le nombre de mots pour chaque URL et le total global.

Vous devez également écrire des tests unitaires pour vérifier :

- La logique de comptage des mots sur des extraits HTML.
- La gestion des erreurs pour les URL invalides ou les délais d'expiration.

Enfin, mesurez les performances :

- Comparez les temps d'exécution avec 1, 2, 4 et 8 goroutines, comme illustré dans le Chapitre 8.

### Considérations éthiques

En tant que développeurs, vous avez la responsabilité de respecter les règles et les normes des sites web que vous explorez. Cela inclut le respect de la norme `robots.txt`, un fichier standardisé utilisé par les sites pour indiquer quelles parties de leur contenu peuvent être explorées par des robots comme le vôtre.

**Qu'est-ce que robots.txt ?** Le fichier `robots.txt` est généralement situé à la racine d'un site (par exemple, `https://www.example.com/robots.txt`). Il contient des directives pour les robots d'exploration, comme des règles d'autorisation (`Allow`) ou d'interdiction (`Disallow`) pour certains chemins. Par exemple, une ligne comme `Disallow: /private` indique que le dossier `/private` ne doit pas être exploré.

**Comment respecter robots.txt ?** Avant de récupérer une URL, vous devez :

1. Vérifier l'existence du fichier `robots.txt` en effectuant une requête HTTP à l'URL racine du site (par exemple, `https://www.example.com/robots.txt`).
2. Analyser le contenu du fichier pour identifier les règles applicables à votre robot (recherchez les directives sous `User-agent: *` ou un `User-agent` spécifique si vous en définissez un).
3. Respecter les directives `Disallow` en évitant d'explorer les chemins interdits. Vous pouvez utiliser une bibliothèque comme `github.com/temoto/robotstxt` pour parser le fichier.
4. Limiter la fréquence de vos requêtes pour éviter de surcharger le serveur (par exemple, ajouter un délai entre les requêtes).

Votre programme doit inclure une vérification de `robots.txt` avant chaque exploration et documenter dans votre rapport comment vous avez implémenté cette conformité. Le non-respect de cette norme peut entraîner des blocages par les sites ou des conséquences juridiques, et sera pris en compte dans l'évaluation de votre projet.

### Ce que vous devez soumettre

Votre soumission doit inclure :

- Votre code source (`crawler.go`, `crawler_test.go`).
- Un rapport de 2 à 3 pages (PDF) décrivant votre implémentation, vos cas de test, les résultats de vos benchmarks, les défis rencontrés, les optimisations effectuées et votre approche pour respecter `robots.txt`.

### Critères d'évaluation

| Critère | Pondération |
|---------|:-----------:|
| Correction et robustesse du robot | 25 % |
| Qualité et couverture des tests unitaires | 20 % |
| Utilisation efficace de la concurrence | 20 % |
| Analyse et optimisation des performances | 15 % |
| Respect des considérations éthiques (robots.txt) | 10 % |
| Clarté et profondeur du rapport | 10 % |

### Directives

- Soumettez tous vos fichiers via la plateforme en ligne du cours à la fin de la semaine 15.
- Votre code doit être bien documenté et formaté avec `gofmt`.
- Votre rapport doit être en PDF, avec une police de 12 points, et inclure votre nom et le titre du projet.
- Vous pouvez remettre une vidéo avec votre travail.
- Les soumissions tardives peuvent recevoir une note de zéro.
