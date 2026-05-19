package main

import (
	"os"
	"path/filepath"
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
	// Chaque segment contient des mots complets — aucun mot n'est coupé
	totalWords := 0
	for _, seg := range segments {
		totalWords += countWords(seg)
	}
	if totalWords != 4 {
		t.Errorf("attendu 4 mots après split, obtenu %d (segments: %v)", totalWords, segments)
	}
}

// TestSplitIntoSegmentsNegativeSize vérifie que splitIntoSegments retourne
// tout le contenu en un seul segment quand la taille est <= 0.
func TestSplitIntoSegmentsNegativeSize(t *testing.T) {
	segments := splitIntoSegments("Hello world", -1)
	if len(segments) != 1 || segments[0] != "Hello world" {
		t.Errorf("attendu 1 segment complet, obtenu %v", segments)
	}
	segments = splitIntoSegments("Hello world", 0)
	if len(segments) != 1 || segments[0] != "Hello world" {
		t.Errorf("attendu 1 segment complet pour size=0, obtenu %v", segments)
	}
}

// TestRunValidFile vérifie que run fonctionne avec un fichier réel.
func TestRunValidFile(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "test.txt")
	os.WriteFile(f, []byte("Hello world from Go"), 0644)

	total, err := run([]string{"cmd", f})
	if err != nil {
		t.Fatalf("erreur inattendue : %v", err)
	}
	if total != 4 {
		t.Errorf("attendu 4 mots, obtenu %d", total)
	}
}

// TestRunWithSegmentSize vérifie que run accepte une taille de segment en argument.
func TestRunWithSegmentSize(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "test.txt")
	os.WriteFile(f, []byte("Hello world from Go"), 0644)

	total, err := run([]string{"cmd", f, "5"})
	if err != nil {
		t.Fatalf("erreur inattendue : %v", err)
	}
	if total != 4 {
		t.Errorf("attendu 4 mots, obtenu %d", total)
	}
}

// TestRunNoArgs vérifie que run retourne une erreur sans arguments.
func TestRunNoArgs(t *testing.T) {
	_, err := run([]string{"cmd"})
	if err == nil {
		t.Fatal("attendu une erreur sans arguments")
	}
}

// TestRunInvalidSegmentSize vérifie que run retourne une erreur pour une taille invalide.
func TestRunInvalidSegmentSize(t *testing.T) {
	_, err := run([]string{"cmd", "file.txt", "abc"})
	if err == nil {
		t.Fatal("attendu une erreur pour taille invalide")
	}
	_, err = run([]string{"cmd", "file.txt", "-5"})
	if err == nil {
		t.Fatal("attendu une erreur pour taille négative")
	}
}

// TestRunMissingFile vérifie que run retourne une erreur pour un fichier inexistant.
func TestRunMissingFile(t *testing.T) {
	_, err := run([]string{"cmd", "fichier_inexistant_xyz.txt"})
	if err == nil {
		t.Fatal("attendu une erreur pour fichier manquant")
	}
}

// TestMainFunction vérifie que main() s'exécute sans panique pour un fichier valide.
// On remplace exitFunc pour éviter que os.Exit ne tue le processus de test.
func TestMainFunction(t *testing.T) {
	oldArgs := os.Args
	oldExit := exitFunc
	defer func() { os.Args = oldArgs; exitFunc = oldExit }()

	dir := t.TempDir()
	f := filepath.Join(dir, "test.txt")
	os.WriteFile(f, []byte("Hello world"), 0644)

	exitFunc = func(code int) {}
	os.Args = []string{"cmd", f}
	main()
}

// TestMainFunctionError vérifie que main() appelle exitFunc(1) quand aucun fichier
// n'est fourni, ce qui couvre la branche d'erreur.
func TestMainFunctionError(t *testing.T) {
	oldArgs := os.Args
	oldExit := exitFunc
	defer func() { os.Args = oldArgs; exitFunc = oldExit }()

	exitCalled := false
	exitFunc = func(code int) { exitCalled = true }
	os.Args = []string{"cmd"}
	main()

	if !exitCalled {
		t.Fatal("attendu un appel à exitFunc")
	}
}

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

// ========== Benchmarks ==========

// BenchmarkSegmentSize mesure la performance en fonction de la taille des segments.
// Plus le segment est grand, moins il y a de goroutines lancées.
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
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				CountWordsConcurrent(benchContent, s.size)
			}
		})
	}
}

// BenchmarkSequentialVsConcurrent compare le comptage séquentiel vs concurrent.
func BenchmarkSequentialVsConcurrent(b *testing.B) {
	b.Run("Sequential", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			countWords(benchContent)
		}
	})

	b.Run("Concurrent_1000", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			CountWordsConcurrent(benchContent, 1000)
		}
	})

	b.Run("Concurrent_10000", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			CountWordsConcurrent(benchContent, 10000)
		}
	})

	b.Run("Concurrent_100000", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			CountWordsConcurrent(benchContent, 100000)
		}
	})
}

// TestCountWordsConcurrentN vérifie que CountWordsConcurrentN retourne le bon
// résultat pour différents nombres de goroutines, y compris les cas limites.
func TestCountWordsConcurrentN(t *testing.T) {
	content := "La programmation concurrente en Go repose sur les goroutines et les canaux pour une exécution efficace"
	expected := len(strings.Fields(content))

	for _, n := range []int{1, 2, 4, 8, 16} {
		result := CountWordsConcurrentN(content, n)
		if result != expected {
			t.Errorf("workers=%d : attendu %d mots, obtenu %d", n, expected, result)
		}
	}

	// Contenu vide → 0
	if CountWordsConcurrentN("", 4) != 0 {
		t.Error("attendu 0 pour contenu vide")
	}
	// Zéro workers → 0
	if CountWordsConcurrentN("hello world", 0) != 0 {
		t.Error("attendu 0 pour 0 workers")
	}
	// Plus de workers que de caractères → doit fonctionner
	if CountWordsConcurrentN("hi", 100) != 1 {
		t.Error("attendu 1 mot pour 'hi' avec 100 workers")
	}
}

// BenchmarkWorkerPool mesure la performance du worker pool réel en fonction
// du nombre de workers, pour comparer avec BenchmarkWorkerCount.
func BenchmarkWorkerPool(b *testing.B) {
	workers := []struct {
		name  string
		count int
	}{
		{"1worker", 1},
		{"2workers", 2},
		{"4workers", 4},
		{"8workers", 8},
		{"16workers", 16},
		{"32workers", 32},
	}
	for _, w := range workers {
		b.Run(w.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				CountWordsConcurrentPool(benchContent, w.count)
			}
		})
	}
}

// TestCountWordsConcurrentPool vérifie que CountWordsConcurrentPool retourne
// le bon résultat pour différents nombres de workers, y compris les cas limites.
func TestCountWordsConcurrentPool(t *testing.T) {
	content := "La programmation concurrente en Go repose sur les goroutines et les canaux pour une exécution efficace"
	expected := len(strings.Fields(content))

	for _, n := range []int{1, 2, 4, 8, 16} {
		result := CountWordsConcurrentPool(content, n)
		if result != expected {
			t.Errorf("pool workers=%d : attendu %d mots, obtenu %d", n, expected, result)
		}
	}

	// Contenu vide → 0
	if CountWordsConcurrentPool("", 4) != 0 {
		t.Error("attendu 0 pour contenu vide")
	}
	// Zéro workers → 0
	if CountWordsConcurrentPool("hello world", 0) != 0 {
		t.Error("attendu 0 pour 0 workers")
	}
	// Plus de workers que de segments → doit fonctionner
	// benchContent fait ~700k chars, 50k chars par segment = ~14 segments
	// avec 32 workers, certains workers ne recevront aucune tâche
	if CountWordsConcurrentPool("hi there", 100) != 2 {
		t.Error("attendu 2 mots avec plus de workers que de segments")
	}
}

// Répond directement à la question de l'énoncé : "Est-ce que la performance
// croît linéairement avec le nombre de goroutines ?"
func BenchmarkWorkerCount(b *testing.B) {
	workers := []struct {
		name  string
		count int
	}{
		{"1worker", 1},
		{"2workers", 2},
		{"4workers", 4},
		{"8workers", 8},
		{"16workers", 16},
		{"32workers", 32},
	}
	for _, w := range workers {
		b.Run(w.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				CountWordsConcurrentN(benchContent, w.count)
			}
		})
	}
}
