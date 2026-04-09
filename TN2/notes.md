# Word Stats — Exercices Git & Go (countLines / countWords / countChars)

Ce guide décrit, étape par étape, la création d’un petit projet Go permettant de compter :
- le nombre de lignes (`countLines`)
- le nombre de mots (`countWords`)
- le nombre de caractères (hors espaces et sauts de ligne) (`countChars`)

Il met également en pratique un workflow Git avec création de branches, fusions (avec et sans conflit), historique, README et tests unitaires.

---

## Prérequis

- Go **1.16+**
- Git installé et configuré
- Un terminal (macOS/Linux) ou PowerShell/WSL (Windows)

---

## Étape 1 — Créer le projet

### Objectif
Créer un nouveau dossier, initialiser Git, et ajouter un fichier `main.go` de base.

### Commandes (terminal)

```bash
cd Documents
mkdir word-stats
cd word-stats
git init
```

### Créer `main.go`

Crée ensuite le fichier `main.go` avec le code suivant :

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

### Premier commit

```bash
git add main.go
git commit -m "Initial commit: ajout de countLines"
```

---

## Étape 2 — Créer la branche `count-words`

### Objectif
Créer une branche pour ajouter une fonction `countWords`.

### Créer la branche

```bash
git checkout -b count-words
```

### Modifier `main.go` : ajouter `countWords`

Ajoute cette fonction :

```go
func countWords(text string) int {
	if text == "" {
		return 0
	}
	return len(strings.Fields(text))
}
```

### Ajouter un appel temporaire dans `main()`

Ajoute (temporairement) l’affichage suivant dans `main()` :

```go
fmt.Printf("Nombre de mots : %d\n", countWords(text))
```

### Commit

```bash
git add main.go
git commit -m "Ajout de countWords"
```


---

## Étape 3 — Créer la branche `count-chars`

### Objectif
Créer une branche pour ajouter une fonction `countChars`.

### Revenir sur `main`

```bash
git checkout main
```

### Créer la branche

```bash
git checkout -b count-chars
```

### Modifier `main.go` : ajouter `countChars`

Ajoute la fonction suivante :

```go
func countChars(text string) int {
	count := 0
	for _, r := range text {
		if r != ' ' && r != '\n' {
			count++
		}
	}
	return count
}
```

### Ajouter un affichage temporaire dans `main()`

Ajoute (temporairement) :

```go
fmt.Printf("Nombre de caractères : %d\n", countChars(text))
```

### Commit

```bash
git add main.go
git commit -m "Ajout de countChars"
```

---

## Étape 4 — Fusion sans conflit

### Objectif
Fusionner la branche `count-words` dans `main` (sans conflit).

### Revenir sur `main`

```bash
git checkout main
```

### Fusionner `count-words`

```bash
git merge count-words
```

### Tester le programme

```bash
go run main.go
```

---

## Étape 5 — Fusion avec conflit (résolution manuelle)

### Objectif
Provoquer un conflit en modifiant `main()` différemment dans deux branches, puis résoudre le conflit.

### 5.1 Dans `count-chars` : afficher uniquement les caractères

1) Basculer sur `count-chars` :

```bash
git checkout count-chars
```

2) Modifier `main()` pour afficher **uniquement** les caractères (et pas les mots).

3) Commit :

```bash
git add main.go
git commit -m "Affichage caractères dans main()"
```

### 5.2 Dans `main` : afficher uniquement les mots

1) Revenir sur `main` :

```bash
git checkout main
```

2) Modifier `main()` pour afficher **uniquement** les mots (et pas les caractères).

3) Commit :

```bash
git add main.go
git commit -m "Affichage mots dans main()"
```

### 5.3 Fusionner `count-chars` dans `main` (conflit attendu)

```bash
git merge count-chars
```

Git va signaler un conflit dans `main.go`.

### 5.4 Résoudre le conflit

1) Ouvrir `main.go` et remplacer la section conflictuelle par :

```go
fmt.Printf("Nombre de mots : %d\n", countWords(text))
fmt.Printf("Nombre de caractères : %d\n", countChars(text))
```

2) Marquer le conflit comme résolu et committer :

```bash
git add main.go
git commit -m "Résolution du conflit: affichage mots + caractères"
```

---

## Étape 6 — Historique + README

### Objectif
Exporter l’historique et rédiger une documentation minimale.

### Générer l’historique

```bash
git log --oneline > history.txt
```

### Créer un `README.md`

Le `README.md` doit contenir au minimum :
- une description du projet
- les fonctions implémentées
- le workflow Git utilisé (branches, merges, résolution de conflit)
- comment exécuter le programme


---

## Tests unitaires (obligatoires)

### Objectif
Ajouter des tests unitaires pour valider les fonctions.

### Créer `main_test.go`

```go
package main

import "testing"

func TestCountLines(t *testing.T) {
	if got := countLines("Hello\nWorld"); got != 2 {
		t.Errorf("countLines() = %d; want 2", got)
	}
}

func TestCountWords(t *testing.T) {
	if got := countWords("Hello World Golang"); got != 3 {
		t.Errorf("countWords() = %d; want 3", got)
	}
}

func TestCountChars(t *testing.T) {
	if got := countChars("Hi!"); got != 3 {
		t.Errorf("countChars() = %d; want 3", got)
	}
}
```

### Lancer les tests

```bash
go test -cover
```
