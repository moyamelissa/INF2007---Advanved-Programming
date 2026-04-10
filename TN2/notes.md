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
<img width="715" height="184" alt="image" src="https://github.com/user-attachments/assets/78b79efa-db8c-4989-956d-c62ddc0a4879" />

### Créer `main.go` dans le terminal

Crée ensuite le fichier `main.go` avec le code suivant :
```bash
notepad main.go
```
Puis copier ce code dans lediteur de text:
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
Enregistrer le fichier et fermer notepad.

Verifier que le programme fonctionne:
```bash
go run main.go
```
<img width="533" height="85" alt="image" src="https://github.com/user-attachments/assets/0ceef55c-9c71-4fb5-8f64-927f39ebb44f" />

### Premier commit
Ajouter etvalider le fichier avec Git

```bash
git add main.go
git commit -m "Initial commit: ajout de countLines"
```
<img width="751" height="79" alt="image" src="https://github.com/user-attachments/assets/85911f31-8eb7-483c-8d36-0365d124a27b" />

---

## Étape 2 — Créer la branche `count-words`

### Objectif
Créer une branche pour ajouter une fonction `countWords`.

### Créer la branche

```bash
git checkout -b count-words
```

Verififaction peut etre fait avec
```bash
git branch
```
<img width="648" height="112" alt="image" src="https://github.com/user-attachments/assets/e6287e17-b031-4644-b532-8609914df7f9" />


### Modifier `main.go` : ajouter `countWords`

Ouvrir lefichier main.go depuis le terminal
```bash
notepad main.go
```
Ajoute cette fonction sous la fonction "countLines" :

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
Enregistrer lesmodifications, fermer le notepad et tester le programme
```bash
go run main.go
```
<img width="488" height="63" alt="image" src="https://github.com/user-attachments/assets/616d64ae-913c-4ca2-b604-37d0a3857657" />

### Commit

```bash
git add main.go
git commit -m "Ajout de countWords"
```
<img width="637" height="92" alt="image" src="https://github.com/user-attachments/assets/2feb9d51-e585-4baa-bc04-d221bcc90d0a" />


---

## Étape 3 — Créer la branche `count-chars`

### Objectif
Créer une branche pour ajouter une fonction `countChars`.

### Revenir sur `main` qui est nomme `master` dans le terminal

```bash
git checkout master
```

Verifier:
```bash
git branch
```
<img width="600" height="119" alt="image" src="https://github.com/user-attachments/assets/31e8e214-3ea6-4fc2-884a-3e3a1fcb509d" />

### Créer la branche

```bash
git checkout -b count-chars
git branch
```
<img width="547" height="128" alt="image" src="https://github.com/user-attachments/assets/716dce83-0b4b-4d7f-8ec7-57d059f9c352" />

### Modifier `main.go` : ajouter `countChars`

ouvrir Notepad
```bash
notepad main.go
```

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
Enregister et ferme le notepad

<img width="480" height="132" alt="image" src="https://github.com/user-attachments/assets/d0306d66-6f4d-4a63-a927-77f2e264a754" />

### Commit

```bash
git add main.go
git commit -m "Ajout de countChars"
```
verification
```bash
git run main.go
```
<img width="617" height="118" alt="image" src="https://github.com/user-attachments/assets/91edd238-d526-46e1-b09c-3317dd15dcb0" />

---

## Étape 4 — Fusion sans conflit

### Objectif
Fusionner la branche `count-words` dans `main` (sans conflit).

### Revenir sur `master`

```bash
git checkout master
```

### Fusionner `count-words`

```bash
git merge count-words
```

### Tester le programme

```bash
go run main.go
```

<img width="724" height="212" alt="image" src="https://github.com/user-attachments/assets/4583a0a9-1dd7-4f93-861e-94faba1c737c" />

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

```bash
notepad main.go
```
Dans la fonction main(), supprimer les autres affichages et garder seulement :
fmt.Printf("Nombre de caractères : %d\n", countChars(text))

Verifier
```bash
go run main.go
```
<img width="559" height="80" alt="image" src="https://github.com/user-attachments/assets/829e1b77-116b-4116-b42c-850556e9e964" />

3) Commit :

```bash
git add main.go
git commit -m "Affichage caractères dans main()"
```
<img width="717" height="90" alt="image" src="https://github.com/user-attachments/assets/0eda1aa0-ea57-47f0-a040-4a3dc0ace0fc" />

### 5.2 Dans `main` : afficher uniquement les mots

1) Revenir sur `master` :

```bash
git checkout master
```

2) Modifier `main()` pour afficher **uniquement** les mots (et pas les caractères).

```bash
notepad main.go
```

Dans main(), supprimer: 
```
fmt.Printf("Nombre de lignes : %d\n", countLines(text))
```

verifier 
```bash
go run main.go
```
<img width="513" height="44" alt="image" src="https://github.com/user-attachments/assets/71d1b7f1-6c83-4830-9bb6-fc72a4f87ebc" />

3) Commit :

```bash
git add main.go
git commit -m "Affichage mots dans main()"
```
<img width="690" height="100" alt="image" src="https://github.com/user-attachments/assets/93480c3e-d025-473f-b87c-34c9a2247abe" />

### 5.3 Fusionner `count-chars` dans `main` (conflit attendu)

```bash
git merge count-chars
```

Git va signaler un conflit dans `main.go`.

<img width="781" height="142" alt="image" src="https://github.com/user-attachments/assets/5e14bf45-360b-4ab5-8a6f-65d0c3ba3294" />

### 5.4 Résoudre le conflit

1) Ouvrir `main.go` et remplacer la section conflictuelle par :
```bash
notepad main.go
```
1.1) Repérer les marqueurs de conflit Git tel que:
```
<<<<<<< HEAD
======= 
>>>>>>> count-chars
```

1.2) Remplacer toute la section conflictuelle
Dans main(), tu dois remplacer tout le bloc, y compris les marqueurs, par la version finale désirée :
fmt.Printf("Nombre de mots : %d\n", countWords(text))
fmt.Printf("Nombre de caractères : %d\n", countChars(text))

1.3)Vérifier qu’il ne reste AUCUN marqueur Git puis sauvarger

```go
fmt.Printf("Nombre de mots : %d\n", countWords(text))
fmt.Printf("Nombre de caractères : %d\n", countChars(text))
```
ca devrasitressembler a ca:
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

func countWords(text string) int {
    if text == "" {
        return 0
    }
    return len(strings.Fields(text))
}

func countChars(text string) int {
    count := 0
    for _, r := range text {
        if r != ' ' && r != '\n' {
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
<img width="861" height="127" alt="image" src="https://github.com/user-attachments/assets/dfb48667-9b66-4097-8bd9-2268e97b2edf" />

2) Marquer le conflit comme résolu et committer :

```bash
git add main.go
git commit -m "Résolution du conflit: affichage mots + caractères"
```
<img width="870" height="90" alt="image" src="https://github.com/user-attachments/assets/3b55b527-07bf-42de-aa4f-8ceb083cefda" />

---

## Étape 6 — Historique + README

### Objectif
Exporter l’historique et rédiger une documentation minimale.

### Générer l’historique

```bash
git log --oneline > history.txt
notepad history.txt
```
<img width="1057" height="151" alt="image" src="https://github.com/user-attachments/assets/bc1fce9c-e4e6-4ae2-bbe2-e484fcd8dfe6" />

### Créer un `README.md`

Le `README.md` doit contenir au minimum :
- une description du projet
- les fonctions implémentées
- le workflow Git utilisé (branches, merges, résolution de conflit)
- comment exécuter le programme

Text du README.md:
``` 
# Projet : Word Stats (INF2007 – Travail 2)

Ce projet est une simulation d’un **workflow Git collaboratif**.  
Il contient trois fonctions Go permettant de compter :

- le nombre de lignes (`countLines`)
- le nombre de mots (`countWords`)
- le nombre de caractères (`countChars`)

---

## Workflow Git utilisé

- Initialisation du dépôt (`git init`)
- Création de la branche `count-words`
- Ajout de la fonction `countWords`
- Création de la branche `count-chars`
- Ajout de la fonction `countChars`
- Fusion sans conflit (`count-words` → `main`)
- Fusion avec conflit (`count-chars` → `main`)
- Résolution manuelle du conflit

---

## Exécution du programme

```bash
go run main.go
```

---

## Tests unitaires

Les tests se trouvent dans `main_test.go` :

```bash
go test -cover
```
```
---

## Tests unitaires (obligatoires)

### Objectif
Ajouter des tests unitaires pour valider les fonctions.

### Créer `main_test.go`

```bash
notepad main_test.go
```

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
### Initialiser un module Go
```bash
go mod init word-stats
```
<img width="496" height="85" alt="image" src="https://github.com/user-attachments/assets/2ad0bf22-f117-4f07-94d4-a4bad0d5b5ba" />

### Lancer les tests

```bash
go test -cover
```
<img width="454" height="83" alt="image" src="https://github.com/user-attachments/assets/9ea683ae-e6c9-414c-88ef-5f4c657bd236" />
