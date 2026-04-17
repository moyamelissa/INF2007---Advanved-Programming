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
