# Travail 2 : Maîtrise de Git pour la collaboration

## Informations générales

- **Date de remise :** Fin de la semaine 6
- **Objectif :** Comprendre les systèmes de contrôle de version distribués, en se concentrant sur Git, comme discuté dans le chapitre 2. Vous appliquerez les concepts de gestion de branches, de fusion et de résolution de conflits dans un projet collaboratif simulé.

## Contexte pratique

Vous êtes développeur dans une petite équipe travaillant sur un projet open-source : une application Go qui génère des statistiques sur un fichier texte (par exemple, comptage de mots ou de lignes). L'équipe utilise Git pour gérer le code source, et chaque membre travaille sur une fonctionnalité spécifique. Votre tâche est de simuler ce scénario en créant un dépôt Git local, en gérant des branches pour deux fonctionnalités, en effectuant des fusions, et en résolvant un conflit intentionnel.

## Tâche à réaliser

Vous allez créer un dépôt Git local et effectuer une série d'opérations pour simuler un workflow collaboratif. Suivez les étapes ci-dessous :

### 1. Initialisation du dépôt

- Créez un répertoire nommé `word-stats`.
- Initialisez un dépôt Git (`git init`).
- Créez un fichier `main.go` avec une fonction simple qui compte le nombre de lignes dans une chaîne (exemple fourni ci-dessous).
- Ajoutez et validez le fichier dans la branche `main` (`git add`, `git commit`).

### 2. Branche pour la fonctionnalité 1

- Créez une branche nommée `count-words` (`git branch` ou `git checkout -b`).
- Modifiez `main.go` pour ajouter une fonction qui compte le nombre de mots (séparés par des espaces).
- Validez les changements (au moins un commit).

### 3. Branche pour la fonctionnalité 2

- Revenez à la branche `main` (`git checkout main`).
- Créez une branche nommée `count-chars`.
- Modifiez `main.go` pour ajouter une fonction qui compte le nombre de caractères (hors espaces).
- Validez les changements (au moins un commit).

### 4. Fusion sans conflit

- Fusionnez la branche `count-words` dans `main` (`git merge count-words`).
- Vérifiez que le code fonctionne (`go run main.go`).

### 5. Fusion avec conflit

- Dans la branche `count-chars`, modifiez la fonction `main` pour afficher le résultat du comptage des caractères.
- Dans la branche `main`, modifiez la même partie de la fonction `main` pour afficher le résultat du comptage des mots (créant un conflit).
- Validez les changements dans les deux branches.
- Tentez de fusionner `count-chars` dans `main`. Résolvez le conflit en combinant les deux affichages (mots et caractères).
- Validez la résolution du conflit.

### 6. Historique et documentation

- Générez un fichier `history.txt` avec l'historique des commits (`git log --oneline`).
- Créez un fichier `README.md` expliquant le projet et les commandes Git utilisées.

### Exemple de `main.go` initial

```go
package main

import (
    "fmt"
    "strings"
)

func countLines(text string) int {
    if text == "" {
        return 0
    }
    return len(strings.Split(text, "\n"))
}

func main() {
    text := "Hello\nWorld\nGolang"
    fmt.Printf("Nombre de lignes : %d\n", countLines(text))
}
```

## Exigences techniques

- **Git :** Démontrer la maîtrise des commandes `git init`, `git add`, `git commit`, `git branch`, `git checkout`, `git merge`, et la résolution de conflits.
- **Code Go :** Les fonctions ajoutées doivent être simples, fonctionnelles, et accompagnées de commentaires expliquant leur rôle.
- **Tests :** Incluez au moins un test unitaire dans `main_test.go` pour chaque fonction (`countLines`, `countWords`, `countChars`). Exemple :

```go
func TestCountLines(t *testing.T) {
    if got := countLines("Hello\nWorld"); got != 2 {
        t.Errorf("countLines() = %d; want 2", got)
    }
}
```

- **Rapport :** Soumettez un rapport (1 page, PDF) décrivant :
  - Votre workflow Git (comment vous avez créé les branches, géré les fusions).
  - Comment vous avez résolu le conflit de fusion.
  - Les défis rencontrés et leur résolution.

## Critères d'évaluation

| Critère | Points | Description |
|---|---|---|
| Correction Git | 40 | Dépôt correctement initialisé, branches créées, fusions effectuées, conflit résolu. |
| Code et tests | 30 | Fonctions Go fonctionnelles, tests unitaires pertinents (couverture ≥ 80 %). |
| Documentation | 15 | README clair, historique des commits, commentaires dans le code. |
| Rapport | 15 | Rapport clair, décrivant le workflow et les défis. |

## Ressources recommandées

- Manuel INF2007, chapitre 2 (pages 26–39 : Systèmes de contrôle de version distribués).
- [Documentation Git](https://git-scm.com/doc)
- [Tutoriel Git](https://git-scm.com/book/fr/v2)
- [A Tour of Go](https://tour.golang.org/)

## Conseils

- Testez vos commandes Git dans un dépôt temporaire pour éviter les erreurs.
- Utilisez `git status` fréquemment pour vérifier l'état du dépôt.
- Pour résoudre les conflits, ouvrez `main.go` dans un éditeur et combinez les modifications manuellement.
- Vérifiez la couverture des tests avec `go test -cover`.
- Rédigez des messages de commit clairs et descriptifs.
- Vous pouvez remettre une vidéo avec votre travail.

**Bonne chance et amusez-vous à collaborer avec Git !**