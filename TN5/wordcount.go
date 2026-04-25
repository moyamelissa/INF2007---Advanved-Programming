package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// countWords compte le nombre de mots dans une chaîne de caractères.
// Un mot est une séquence de caractères séparée par des espaces.
func countWords(text string) int {
	return len(strings.Fields(text))
}

// splitIntoSegments divise le contenu en segments d'environ segmentSize caractères,
// en s'assurant de ne pas couper au milieu d'un mot. Chaque coupure est décalée
// vers l'espace le plus proche à droite de la position segmentSize.
func splitIntoSegments(content string, segmentSize int) []string {
	if len(content) == 0 {
		return nil
	}
	if segmentSize <= 0 {
		return []string{content}
	}

	var segments []string
	for len(content) > 0 {
		if len(content) <= segmentSize {
			segments = append(segments, content)
			break
		}

		// Chercher la fin du mot à partir de la position segmentSize
		end := segmentSize
		// Avancer jusqu'à la fin du mot courant (pas d'espace)
		for end < len(content) && content[end] != ' ' && content[end] != '\n' && content[end] != '\t' && content[end] != '\r' {
			end++
		}

		segments = append(segments, content[:end])
		// Sauter les espaces entre les segments
		for end < len(content) && (content[end] == ' ' || content[end] == '\n' || content[end] == '\t' || content[end] == '\r') {
			end++
		}
		content = content[end:]
	}

	return segments
}

// countWordsInSegment est la fonction exécutée par chaque goroutine.
// Elle compte les mots dans le segment donné et envoie le résultat sur le canal.
func countWordsInSegment(segment string, ch chan<- int) {
	ch <- countWords(segment)
}

// CountWordsConcurrent divise le contenu en segments et lance une goroutine par
// segment pour compter les mots en parallèle. Les résultats sont collectés via
// un canal et sommés dans la goroutine principale.
//
// Paramètres :
//   - content : le texte complet à analyser.
//   - segmentSize : taille approximative de chaque segment en caractères.
//
// Retour :
//   - int : le nombre total de mots dans le contenu.
func CountWordsConcurrent(content string, segmentSize int) int {
	segments := splitIntoSegments(content, segmentSize)
	if len(segments) == 0 {
		return 0
	}

	ch := make(chan int, len(segments))

	// Lancer une goroutine par segment
	for _, seg := range segments {
		go countWordsInSegment(seg, ch)
	}

	// Collecter les résultats depuis le canal
	total := 0
	for range segments {
		total += <-ch
	}

	return total
}

// run contient la logique principale du programme, extraite de main pour
// permettre les tests unitaires. Elle retourne le nombre total de mots et une erreur.
func run(args []string) (int, error) {
	if len(args) < 2 {
		return 0, fmt.Errorf("usage: go run wordcount.go <fichier> [taille_segment]")
	}

	filePath := args[1]

	segmentSize := 1000 // taille par défaut
	if len(args) >= 3 {
		s, err := strconv.Atoi(args[2])
		if err != nil || s <= 0 {
			return 0, fmt.Errorf("taille de segment invalide %q", args[2])
		}
		segmentSize = s
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return 0, fmt.Errorf("lecture du fichier : %v", err)
	}

	content := string(data)
	total := CountWordsConcurrent(content, segmentSize)

	fmt.Printf("Nombre total de mots : %d\n", total)
	fmt.Printf("Taille du fichier    : %d octets\n", len(data))
	fmt.Printf("Taille des segments  : %d caractères\n", segmentSize)

	return total, nil
}

// exitFunc permet de remplacer os.Exit dans les tests (cf. Ch. 11).
var exitFunc = os.Exit

func main() {
	total, err := run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erreur : %v\n", err)
		exitFunc(1)
	}
	_ = total
}
