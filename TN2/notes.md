# Word Stats — Exercices Git & Go (countLines / countWords / countChars)

Ce guide décrit, étape par étape, la création d’un petit projet Go permettant de compter :
- le nombre de lignes (`countLines`)
- le nombre de mots (`countWords`)
- le nombre de caractères (**hors espaces**) (`countChars`)

Il met également en pratique un workflow Git avec création de branches, fusions (avec et sans conflit), historique, README et tests unitaires.

---

## Prérequis

- Go **1.16+**
- Git installé et configuré
- **Windows + Command Prompt (cmd.exe)** (les commandes ci-dessous sont compatibles)

> Note : dans ce guide, on utilise la branche **main** (recommandé).  
> Si votre dépôt utilise `master`, vous pouvez soit garder `master`, soit renommer en `main`.

---

## Étape 1 — Créer le projet

### Objectif
Créer un nouveau dossier, initialiser Git, initialiser Go module, et ajouter un fichier `main.go` de base.

### Commandes (terminal)

```bash
cd %USERPROFILE%\Documents
mkdir word-stats
cd word-stats
git init
```

### (Important) Mettre la branche principale sur `main`

```bash
git branch -M main
```

Vérifier :

```bash
git branch
```

### Initialiser un module Go

```bash
go mod init word-stats
```

> `go.mod` sera créé automatiquement.

<img width="1340" height="135" alt="image" src="https://github.com/user-attachments/assets/d3ceab0c-942f-4844-af74-6d1d36178d1b" />


### Créer `main.go` dans le terminal

Créer le fichier :

```bash
notepad main.go
```

Puis copier ce code dans l’éditeur de texte :

```go
package main

import (
	"fmt"
	"strings"
)

// countLines retourne le nombre de lignes dans une chaîne (séparées par '\n').
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

Enregistrer le fichier et fermer Notepad.

### Vérifier que le programme fonctionne

```bash
go run main.go
```
<img width="959" height="55" alt="image" src="https://github.com/user-attachments/assets/cef04d7c-a5ee-4c49-9c0a-e99d0c726a89" />

### Premier commit

Ajouter et valider le fichier avec Git :

```bash
git add main.go go.mod
git commit -m "Initial commit: ajout de countLines"
```
<img width="1251" height="176" alt="image" src="https://github.com/user-attachments/assets/dcca8c4c-8727-4e5e-ab32-564a5f9cf2c2" />

#### ⚠️ Avertissement possible (normal sur Windows)

Après `git add`, vous pouvez voir ce message :

```text
warning: in the working copy of 'go.mod', LF will be replaced by CRLF the next time Git touches it
```

**Explication :**
- `LF` = fin de ligne format Linux/macOS (`\n`)
- `CRLF` = fin de ligne format Windows (`\r\n`)
- Git vous informe simplement qu’il pourrait convertir automatiquement les fins de ligne sur Windows.
- **Ce n’est pas une erreur** et votre projet Go fonctionne quand même.
- Optionnel :Vous pouvez forcer des fins de ligne `LF` pour les fichiers du projet avec un fichier `.gitattributes`.
  
---

## Étape 2 — Créer la branche `count-words`

### Objectif
Créer une branche pour ajouter une fonction `countWords`.

### Créer la branche

```bash
git checkout -b count-words
```

Vérification possible avec :

```bash
git branch
```
<img width="1003" height="129" alt="image" src="https://github.com/user-attachments/assets/e4b1a1c4-3894-4c93-8fd2-906a17ac237d" />

### Modifier `main.go` : ajouter `countWords`

Ouvrir le fichier `main.go` :

```bash
notepad main.go
```

Ajouter cette fonction **sous** `countLines` :

```go
// countWords retourne le nombre de mots dans une chaîne (séparés par des espaces/whitespace).
func countWords(text string) int {
	if text == "" {
		return 0
	}
	return len(strings.Fields(text))
}
```

### Ajouter un appel temporaire dans `main()`

Dans `main()`, ajouter (temporairement) l’affichage suivant :

```go
fmt.Printf("Nombre de mots : %d\n", countWords(text))
```

Enregistrer, fermer Notepad et tester :

```bash
go run main.go
```
<img width="892" height="73" alt="image" src="https://github.com/user-attachments/assets/c9ec3328-3ffb-4de2-8dec-047f3e773380" />

### Commit

```bash
git add main.go
git commit -m "Ajout de countWords"
```
<img width="976" height="110" alt="image" src="https://github.com/user-attachments/assets/022a47cb-1213-49de-999a-19b5c96d9ec0" />

---

## Étape 3 — Créer la branche `count-chars`

### Objectif
Créer une branche pour ajouter une fonction `countChars`.

### Revenir sur `main`

```bash
git checkout main
```

Vérifier :

```bash
git branch
```

### Créer la branche

```bash
git checkout -b count-chars
```

### Modifier `main.go` : ajouter `countChars`

Ouvrir `main.go` :

```bash
notepad main.go
```

Ajouter cette fonction **sous** `countLines` :

```go
// countChars retourne le nombre de caractères en excluant les espaces et les sauts de ligne.
func countChars(text string) int {
	count := 0
	for _, r := range text {
		// On exclut l'espace et les fins de ligne (plus précis pour Windows/Linux)
		if r != ' ' && r != '\n' && r != '\r' && r != '\t' {
			count++
		}
	}
	return count
}
```

### Ajouter un affichage temporaire dans `main()`

Ajouter (temporairement) :

```go
fmt.Printf("Nombre de caractères : %d\n", countChars(text))
```

Enregistrer, fermer Notepad, puis tester :

```bash
go run main.go
```

<img width="1149" height="74" alt="image" src="https://github.com/user-attachments/assets/941d4b89-3d6f-4114-b1a9-56a9330d7000" />


### Commit

```bash
git add main.go
git commit -m "Ajout de countChars"
```
<img width="1139" height="108" alt="image" src="https://github.com/user-attachments/assets/d7465a69-72b2-4f9f-ad6d-6c10bc394743" />


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
<img width="1262" height="112" alt="image" src="https://github.com/user-attachments/assets/f2d60ccd-7a5c-4bf1-a3ab-5f99452db7e1" />


### Tester le programme

```bash
go run main.go
```
<img width="1182" height="72" alt="image" src="https://github.com/user-attachments/assets/3597948d-5e35-456b-aaa9-dd38c1a3fa8f" />

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

Ouvrir :

```bash
notepad main.go
```

Dans `main()`, garder seulement :

```go
fmt.Printf("Nombre de caractères : %d\n", countChars(text))
```

Tester :

```bash
go run main.go
```
<img width="898" height="49" alt="image" src="https://github.com/user-attachments/assets/1faeaff9-7060-455e-8c32-5e7ca84d1efb" />

3) Commit :

```bash
git add main.go
git commit -m "Affichage caractères dans main()"
```
<img width="974" height="68" alt="image" src="https://github.com/user-attachments/assets/af62dd31-da39-43ec-b1d2-f5e8376828f1" />

### 5.2 Dans `main` : afficher uniquement les mots

1) Revenir sur `main` :

```bash
git checkout main
```

2) Ouvrir le fichier :

```bash
notepad main.go
```

3) Dans la fonction `main()`, **remplacer** les anciens affichages pour ne garder **qu’un seul affichage : les mots**.

La fonction `main()` doit maintenant être :

```go
func main() {
	text := "Hello\nWorld\nGolang"
	fmt.Printf("Nombre de mots : %d\n", countWords(text))
}
```

4) Enregistrer et fermer Notepad, puis tester :

```bash
go run main.go
```

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

1) Ouvrir `main.go` :

```bash
notepad main.go
```

2) Repérer les marqueurs de conflit Git :

```text
<<<<<<< HEAD
=======
>>>>>>> count-chars
```

3) Remplacer **toute** la section conflictuelle (y compris les marqueurs) par la version finale souhaitée :

```go
fmt.Printf("Nombre de mots : %d\n", countWords(text))
fmt.Printf("Nombre de caractères : %d\n", countChars(text))
```

4) Vérifier qu’il ne reste **aucun** marqueur (`<<<<<<<`, `=======`, `>>>>>>>`), puis enregistrer.

5) Marquer le conflit comme résolu et committer :

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
notepad history.txt
```

### Créer un `README.md`

Créer le fichier :

```bash
notepad README.md
```

Le `README.md` doit contenir au minimum :
- une description du projet
- les fonctions implémentées
- le workflow Git utilisé (branches, merges, résolution de conflit)
- comment exécuter le programme
- comment lancer les tests

Exemple de contenu :

```markdown
# Projet : Word Stats (INF2007 – Travail 2)

Ce projet est une simulation d’un **workflow Git collaboratif**.  
Il contient trois fonctions Go permettant de compter :
- le nombre de lignes (`countLines`)
- le nombre de mots (`countWords`)
- le nombre de caractères (`countChars`)

## Workflow Git utilisé
- Initialisation du dépôt (`git init`)
- Création de la branche `count-words` + commit
- Création de la branche `count-chars` + commit
- Fusion sans conflit (`count-words` → `main`)
- Fusion avec conflit (`count-chars` → `main`) + résolution manuelle

## Exécution
```bash
go run main.go
```

## Tests
```bash
go test -cover
```
```

---

## Tests unitaires (obligatoires)

### Objectif
Ajouter un test unitaire par fonction (`countLines`, `countWords`, `countChars`).

### Créer `main_test.go`

```bash
notepad main_test.go
```

Copier ce code :

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

### Lancer les tests + couverture

```bash
go test -cover
```

---
