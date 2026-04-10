# TN2 - Word Stats — Exercices Git & Go
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


> **Résultat attendu :**
> ```
> Initialized empty Git repository in C:/Users/.../word-stats/.git/
> go: creating new go.mod: module word-stats
> ```


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

> **Résultat attendu :**
> ```
> Nombre de lignes : 3
> ```

### Premier commit

Ajouter et valider le fichier avec Git :

```bash
git add main.go go.mod
git commit -m "Initial commit: ajout de countLines"
```

> **Résultat attendu :**
> ```
> [main (root-commit) xxxxxxx] Initial commit: ajout de countLines
>  2 files changed, ...
>  create mode 100644 go.mod
>  create mode 100644 main.go
> ```

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
- Optionnel : Vous pouvez forcer des fins de ligne `LF` pour les fichiers du projet avec un fichier `.gitattributes`.
  
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

> **Résultat attendu :**
> ```
> Switched to a new branch 'count-words'
> * count-words
>   main
> ```

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

> **Résultat attendu :**
> ```
> Nombre de lignes : 3
> Nombre de mots : 3
> ```

### Commit

```bash
git add main.go
git commit -m "Ajout de countWords"
```

> **Résultat attendu :**
> ```
> [count-words xxxxxxx] Ajout de countWords
>  1 file changed, ...
> ```

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


> **Résultat attendu :**
> ```
> Nombre de lignes : 3
> Nombre de caractères : 16
> ```


### Commit

```bash
git add main.go
git commit -m "Ajout de countChars"
```

> **Résultat attendu :**
> ```
> [count-chars xxxxxxx] Ajout de countChars
>  1 file changed, ...
> ```


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

> **Résultat attendu :**
> ```
> Updating xxxxxxx..xxxxxxx
> Fast-forward
>  main.go | ...
>  1 file changed, ...
> ```


### Tester le programme

```bash
go run main.go
```

> **Résultat attendu :**
> ```
> Nombre de lignes : 3
> Nombre de mots : 3
> ```

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

Enregistrer et fermer Notepad, puis tester :

```bash
go run main.go
```

> **Résultat attendu :**
> ```
> Nombre de caractères : 16
> ```

3) Commit :

```bash
git add main.go
git commit -m "Affichage caractères dans main()"
```

> **Résultat attendu :**
> ```
> [count-chars xxxxxxx] Affichage caractères dans main()
>  1 file changed, ...
> ```

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

> **Résultat attendu :**
> ```
> Nombre de mots : 3
> ```

5) Commit :

```bash
git add main.go
git commit -m "Affichage mots dans main()"
```

> **Résultat attendu :**
> ```
> [main xxxxxxx] Affichage mots dans main()
>  1 file changed, ...
> ```

### 5.3 Fusionner `count-chars` dans `main` (conflit attendu)

```bash
git merge count-chars
```

Git va signaler un conflit dans `main.go`.


> **Résultat attendu :**
> ```
> Auto-merging main.go
> CONFLICT (content): Merge conflict in main.go
> Automatic merge failed; fix conflicts and then commit the result.
> ```


### 5.4 Résoudre le conflit (version conforme au devoir)

> Objectif : supprimer les marqueurs de conflit et conserver **les deux fonctionnalités** (`countWords` et `countChars`), puis afficher **mots + caractères** dans `main()`.

1) Ouvrir `main.go` :

```bash
notepad main.go
```

2) Repérer les marqueurs de conflit Git (ils ressemblent à ceci) :

```text
<<<<<<< HEAD
... (version de la branche actuelle, ex: main)
=======
... (version de l’autre branche, ex: count-chars)
>>>>>>> count-chars
```

3) Résoudre le conflit en faisant 2 corrections :

#### A) Conserver les deux fonctions
Dans le fichier, Git a mis en conflit deux blocs (`countWords` d’un côté et `countChars` de l’autre).
Vous devez :
- **supprimer** toutes les lignes `<<<<<<<`, `=======`, `>>>>>>>`
- garder **les deux fonctions complètes** (une après l’autre), par exemple :

```go
// countWords retourne le nombre de mots dans une chaîne (séparés par des espaces/whitespace).
func countWords(text string) int {
	if text == "" {
		return 0
	}
	return len(strings.Fields(text))
}

// countChars retourne le nombre de caractères en excluant les espaces et les sauts de ligne.
func countChars(text string) int {
	count := 0
	for _, r := range text {
		if r != ' ' && r != '\n' && r != '\r' && r != '\t' {
			count++
		}
	}
	return count
}
```

#### B) Combiner les affichages dans `main()`
Dans `main()`, remplacer le bloc en conflit par **les deux affichages** :

```go
fmt.Printf("Nombre de mots : %d\n", countWords(text))
fmt.Printf("Nombre de caractères : %d\n", countChars(text))
```

Par exemple, `main()` doit ressembler à ceci :

```go
func main() {
	text := "Hello\nWorld\nGolang"
	fmt.Printf("Nombre de mots : %d\n", countWords(text))
	fmt.Printf("Nombre de caractères : %d\n", countChars(text))
}
```

4) Vérifier qu’il ne reste **aucun** marqueur (`<<<<<<<`, `=======`, `>>>>>>>`), puis enregistrer.

5) Marquer le conflit comme résolu et committer :

```bash
git add main.go
git commit -m "Résolution du conflit: affichage mots + caractères"
```

> **Résultat attendu :**
> ```
> [main xxxxxxx] Résolution du conflit: affichage mots + caractères
> ```

```bash
go run main.go
```

> **Résultat attendu :**
> ```
> Nombre de mots : 3
> Nombre de caractères : 16
> ```

---

## Étape 6 — Historique + README

### Objectif
Exporter l’historique et rédiger une documentation minimale.

### Générer l’historique

```bash
git log --oneline > history.txt
notepad history.txt
```

> **Résultat attendu (exemple, les hash seront différents) :**
> ```
> xxxxxxx Résolution du conflit: affichage mots + caractères
> xxxxxxx Affichage mots dans main()
> xxxxxxx Affichage caractères dans main()
> xxxxxxx Ajout de countWords
> xxxxxxx Ajout de countChars
> xxxxxxx Initial commit: ajout de countLines
> ```

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

    go run main.go

## Tests

    go test -cover
```

### Commit des fichiers de documentation

```bash
git add history.txt README.md
git commit -m "Ajout de history.txt et README.md"
```

---

## Étape 7 — Tests unitaires

### Objectif
Ajouter des tests unitaires pour chaque fonction (`countLines`, `countWords`, `countChars`) ainsi qu'un test pour `main()` afin d'atteindre une couverture de **100 %**.

### Créer `main_test.go`

```bash
notepad main_test.go
```

Copier ce code :

```go
package main

import (
	"os"
	"testing"
)

// TestMainOutput vérifie que la fonction main s'exécute sans erreur.
func TestMainOutput(t *testing.T) {
	// Redirige stdout pour éviter l'affichage pendant les tests
	old := os.Stdout
	os.Stdout, _ = os.Create(os.DevNull)
	defer func() { os.Stdout = old }()

	main()
}

func TestCountLines_Empty(t *testing.T) {
	if got := countLines(""); got != 0 {
		t.Errorf("countLines(\"\") = %d; want 0", got)
	}
}

func TestCountLines_TrailingNewline(t *testing.T) {
	// "a\n" split => ["a", ""] => 2 lines
	if got := countLines("a\n"); got != 2 {
		t.Errorf("countLines(\"a\\n\") = %d; want 2", got)
	}
}

func TestCountWords_Empty(t *testing.T) {
	if got := countWords(""); got != 0 {
		t.Errorf("countWords(\"\") = %d; want 0", got)
	}
}

func TestCountWords_MultipleSpacesAndNewlines(t *testing.T) {
	// strings.Fields splits on any whitespace
	if got := countWords(" Hello   World \n Golang "); got != 3 {
		t.Errorf("countWords(...) = %d; want 3", got)
	}
}

func TestCountChars_Empty(t *testing.T) {
	if got := countChars(""); got != 0 {
		t.Errorf("countChars(\"\") = %d; want 0", got)
	}
}

func TestCountChars_IgnoresSpacesAndNewlines(t *testing.T) {
	// Only letters count: 'H'(1) 'i'(1) => 2
	if got := countChars("H i\n"); got != 2 {
		t.Errorf("countChars(\"H i\\n\") = %d; want 2", got)
	}
}

func TestCountWords_OnlySpaces(t *testing.T) {
	// strings.Fields("   ") => [] so 0
	if got := countWords("   "); got != 0 {
		t.Errorf("countWords(\"   \") = %d; want 0", got)
	}
}

func TestCountChars_OnlyWhitespace(t *testing.T) {
	// All characters are whitespace, so count should be 0
	if got := countChars(" \n\r\t "); got != 0 {
		t.Errorf("countChars(\" \\\\n\\\\r\\\\t \") = %d; want 0", got)
	}
}

func TestCountLines_SingleLine(t *testing.T) {
	if got := countLines("Hello"); got != 1 {
		t.Errorf("countLines(\"Hello\") = %d; want 1", got)
	}
}
```

### Explication des tests

| Test | Fonction testée | Ce qu'il vérifie |
|------|----------------|------------------|
| `TestMainOutput` | `main()` | Exécution sans erreur (couvre `main()` pour atteindre 100 %) |
| `TestCountLines_Empty` | `countLines` | Chaîne vide → 0 lignes |
| `TestCountLines_TrailingNewline` | `countLines` | Saut de ligne en fin de chaîne |
| `TestCountLines_SingleLine` | `countLines` | Une seule ligne sans `\n` |
| `TestCountWords_Empty` | `countWords` | Chaîne vide → 0 mots |
| `TestCountWords_MultipleSpacesAndNewlines` | `countWords` | Espaces multiples et `\n` |
| `TestCountWords_OnlySpaces` | `countWords` | Uniquement des espaces → 0 mots |
| `TestCountChars_Empty` | `countChars` | Chaîne vide → 0 caractères |
| `TestCountChars_IgnoresSpacesAndNewlines` | `countChars` | Exclut espaces et `\n` |
| `TestCountChars_OnlyWhitespace` | `countChars` | Tous les types de whitespace → 0 |

### Lancer les tests + couverture

```bash
go test -cover
```

Résultat attendu :

```
PASS
coverage: 100.0% of statements
ok      word-stats
```

### Commit des tests

```bash
git add main_test.go
git commit -m "Ajout des tests unitaires (couverture 100%)"
```

---

## Résumé final

### Fichier `main.go` final

Après toutes les étapes (résolution du conflit incluse), le fichier `main.go` doit ressembler à ceci :

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

// countWords retourne le nombre de mots dans une chaîne (séparés par des espaces/whitespace).
func countWords(text string) int {
	if text == "" {
		return 0
	}
	return len(strings.Fields(text))
}

// countChars retourne le nombre de caractères en excluant les espaces et les sauts de ligne.
func countChars(text string) int {
	count := 0
	for _, r := range text {
		if r != ' ' && r != '\n' && r != '\r' && r != '\t' {
			count++
		}
	}
	return count
}

func main() {
	text := "Hello\nWorld\nGolang"
	fmt.Printf("Nombre de mots : %d\n", countWords(text))
	fmt.Printf("Nombre de caractères : %d\n", countChars(text))
}
```

### Commandes Git utilisées

| Commande | Utilisation |
|----------|-------------|
| `git init` | Initialiser le dépôt |
| `git branch -M main` | Renommer la branche principale |
| `git add` | Ajouter des fichiers au staging |
| `git commit -m "..."` | Valider les changements |
| `git branch` | Lister les branches |
| `git checkout -b <nom>` | Créer et basculer sur une branche |
| `git checkout <nom>` | Basculer sur une branche existante |
| `git merge <nom>` | Fusionner une branche dans la branche courante |
| `git log --oneline` | Afficher l'historique des commits |

### Structure finale du projet

```
word-stats/
├── go.mod
├── main.go
├── main_test.go
├── history.txt
└── README.md
```

### Résultat de couverture

```
PASS
coverage: 100.0% of statements
ok      word-stats
```
