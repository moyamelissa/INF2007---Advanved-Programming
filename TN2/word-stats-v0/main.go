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
