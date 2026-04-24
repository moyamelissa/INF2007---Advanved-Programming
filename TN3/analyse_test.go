package analyse

import (
	"testing"
)

// TestAnalyseDonneesValides vérifie le comportement avec des données entièrement valides
// pour un capteur donné.
func TestAnalyseDonneesValides(t *testing.T) {
	data := []uint32{
		0x00000105, // ID=5, bit 8=1
		0x00000205, // ID=5, bit 9=1
		0x00000105, // ID=5, bit 8=1
		0x00000003, // ID=3, bit 8=0 (pas le capteur 5)
		0x00000005, // ID=5, aucune mesure (bits 8–31 tous à 0)
		0x00800005, // ID=5, bit 23=1
	}

	counts, err := Analyse(data, 5)
	if err != nil {
		t.Fatalf("erreur inattendue : %v", err)
	}

	// bit 8 (index 0) : 2 fois, bit 9 (index 1) : 1 fois, bit 23 (index 15) : 1 fois
	expected := [24]int{}
	expected[0] = 2  // bit 8
	expected[1] = 1  // bit 9
	expected[15] = 1 // bit 23

	if counts != expected {
		t.Errorf("résultat incorrect\n  obtenu   : %v\n  attendu  : %v", counts, expected)
	}
}

// TestAnalyseBit7Invalide vérifie qu'une erreur est retournée lorsque le bit 7
// (bit de validation) est à 1 dans une entrée.
func TestAnalyseBit7Invalide(t *testing.T) {
	data := []uint32{
		0x00000105, // ID=5, valide
		0x00000080, // ID=0, bit 7=1 → invalide
	}

	_, err := Analyse(data, 5)
	if err == nil {
		t.Fatal("une erreur était attendue pour le bit 7 à 1, mais aucune erreur retournée")
	}
}

// TestAnalysePlusieursBitsValeur vérifie qu'une erreur est retournée lorsque plus
// d'un bit parmi les bits 8 à 31 est à 1.
func TestAnalysePlusieursBitsValeur(t *testing.T) {
	data := []uint32{
		0x00000305, // ID=5, bits 8 et 9 à 1 → invalide
	}

	_, err := Analyse(data, 5)
	if err == nil {
		t.Fatal("une erreur était attendue pour plusieurs bits à 1 (bits 8–31), mais aucune erreur retournée")
	}
}

// TestAnalyseCapteurInvalide vérifie qu'une erreur est retournée lorsque
// l'identifiant du capteur dépasse 127.
func TestAnalyseCapteurInvalide(t *testing.T) {
	data := []uint32{
		0x00000105,
	}

	_, err := Analyse(data, 200)
	if err == nil {
		t.Fatal("une erreur était attendue pour un capteur > 127, mais aucune erreur retournée")
	}
}

// TestAnalyseTableauVide vérifie que la fonction gère correctement un tableau vide.
func TestAnalyseTableauVide(t *testing.T) {
	data := []uint32{}
	counts, err := Analyse(data, 0)
	if err != nil {
		t.Fatalf("erreur inattendue : %v", err)
	}

	expected := [24]int{}
	if counts != expected {
		t.Errorf("résultat incorrect pour tableau vide\n  obtenu  : %v\n  attendu : %v", counts, expected)
	}
}

// TestAnalyseExempleEnonce reproduit l'exemple donné dans l'énoncé du travail.
func TestAnalyseExempleEnonce(t *testing.T) {
	data := []uint32{
		0x00000105, // ID=5, bit 7=0, bit 8=1 (valide)
		0x00000205, // ID=5, bit 7=0, bit 9=1 (valide)
		0x00000080, // ID=0, bit 7=1 (invalide)
		0x00000305, // ID=5, bit 7=0, bits 8 et 9=1 (invalide)
	}

	_, err := Analyse(data, 5)
	if err == nil {
		t.Fatal("une erreur était attendue pour l'exemple de l'énoncé, mais aucune erreur retournée")
	}
}
