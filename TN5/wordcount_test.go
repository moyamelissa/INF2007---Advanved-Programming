package main

import (
	"strings"
	"testing"
)

// ========== Tests unitaires ==========

// TestCountWordsEmpty vérifie le comptage de mots pour un contenu vide.
func TestCountWordsEmpty(t *testing.T) {
	result := CountWordsConcurrent("", 100)
	if result != 0 {
		t.Errorf("attendu 0 mots pour un fichier vide, obtenu %d", result)
	}
}

// TestCountWordsSingleWord vérifie le comptage pour un seul mot.
func TestCountWordsSingleWord(t *testing.T) {
	result := CountWordsConcurrent("Bonjour", 100)
	if result != 1 {
		t.Errorf("attendu 1 mot, obtenu %d", result)
	}
}

// TestCountWordsMultipleLines vérifie le comptage pour plusieurs lignes.
func TestCountWordsMultipleLines(t *testing.T) {
	content := "Bonjour le monde\nComment allez vous\nTrès bien merci"
	// 3 + 3 + 3 = 9 mots
	result := CountWordsConcurrent(content, 10)
	if result != 9 {
		t.Errorf("attendu 9 mots, obtenu %d", result)
	}
}

// TestCountWordsMultipleSpaces vérifie la gestion des espaces multiples.
func TestCountWordsMultipleSpaces(t *testing.T) {
	content := "  mot1   mot2   mot3  "
	result := CountWordsConcurrent(content, 5)
	if result != 3 {
		t.Errorf("attendu 3 mots avec espaces multiples, obtenu %d", result)
	}
}

// TestCountWordsConsistency vérifie que le résultat concurrent est identique
// au résultat séquentiel pour différentes tailles de segment.
func TestCountWordsConsistency(t *testing.T) {
	content := "La programmation concurrente en Go repose sur les goroutines et les canaux pour une exécution efficace"
	expected := len(strings.Fields(content)) // 16 mots

	segmentSizes := []int{1, 5, 10, 20, 50, 100, 500}
	for _, size := range segmentSizes {
		result := CountWordsConcurrent(content, size)
		if result != expected {
			t.Errorf("segment=%d : attendu %d mots, obtenu %d", size, expected, result)
		}
	}
}

// TestSplitIntoSegments vérifie que les segments ne coupent pas les mots.
func TestSplitIntoSegments(t *testing.T) {
	content := "Hello world from Go"

	segments := splitIntoSegments(content, 7)
	// Chaque segment devrait contenir des mots complets
	totalWords := 0
	for _, seg := range segments {
		totalWords += countWords(seg)
	}
	if totalWords != 4 {
		t.Errorf("attendu 4 mots après split, obtenu %d (segments: %v)", totalWords, segments)
	}
}

// ========== Benchmarks ==========

// generateLargeContent génère un contenu de test avec environ n mots.
func generateLargeContent(nWords int) string {
	words := []string{"Lorem", "ipsum", "dolor", "sit", "amet", "consectetur",
		"adipiscing", "elit", "sed", "do", "eiusmod", "tempor", "incididunt",
		"ut", "labore", "et", "dolore", "magna", "aliqua", "programmation"}
	var b strings.Builder
	for i := 0; i < nWords; i++ {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(words[i%len(words)])
	}
	return b.String()
}

var benchContent = generateLargeContent(100_000) // ~100k mots

// BenchmarkSegmentSize mesure la performance en fonction de la taille des segments.
// Plus le segment est petit, plus il y a de goroutines lancées.
func BenchmarkSegmentSize(b *testing.B) {
	sizes := []struct {
		name string
		size int
	}{
		{"10chars", 10},
		{"100chars", 100},
		{"500chars", 500},
		{"1000chars", 1000},
		{"5000chars", 5000},
		{"10000chars", 10000},
		{"50000chars", 50000},
		{"100000chars", 100000},
		{"AllInOne", len(benchContent)},
	}

	for _, s := range sizes {
		b.Run(s.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				CountWordsConcurrent(benchContent, s.size)
			}
		})
	}
}

// BenchmarkSequentialVsConcurrent compare le comptage séquentiel vs concurrent.
func BenchmarkSequentialVsConcurrent(b *testing.B) {
	b.Run("Sequential", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			countWords(benchContent)
		}
	})

	b.Run("Concurrent_1000", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			CountWordsConcurrent(benchContent, 1000)
		}
	})

	b.Run("Concurrent_10000", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			CountWordsConcurrent(benchContent, 10000)
		}
	})

	b.Run("Concurrent_100000", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			CountWordsConcurrent(benchContent, 100000)
		}
	})
}
