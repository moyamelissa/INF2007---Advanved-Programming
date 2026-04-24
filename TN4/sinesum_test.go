package main

import (
	"math"
	"testing"
)

const tolerance = 1e-9

// ========== Tests unitaires ==========

// TestComputeSineSumInt vérifie la correction de computeSineSum pour un petit tableau d'entiers.
func TestComputeSineSumInt(t *testing.T) {
	data := []int{1, 2, 3}
	expected := math.Sin(1) + math.Sin(2) + math.Sin(3)

	result, err := computeSineSum("int", data)
	if err != nil {
		t.Fatalf("erreur inattendue : %v", err)
	}
	if math.Abs(result-expected) > tolerance {
		t.Errorf("résultat incorrect : obtenu %.10f, attendu %.10f", result, expected)
	}
}

// TestComputeSineSumFloat vérifie la correction de computeSineSum pour un petit tableau de flottants.
func TestComputeSineSumFloat(t *testing.T) {
	data := []float64{0.1, 0.2, 0.3}
	expected := math.Sin(0.1) + math.Sin(0.2) + math.Sin(0.3)

	result, err := computeSineSum("float", data)
	if err != nil {
		t.Fatalf("erreur inattendue : %v", err)
	}
	if math.Abs(result-expected) > tolerance {
		t.Errorf("résultat incorrect : obtenu %.10f, attendu %.10f", result, expected)
	}
}

// TestComputeSineSumInvalidType vérifie qu'une erreur est retournée pour un type invalide.
func TestComputeSineSumInvalidType(t *testing.T) {
	data := []int{1, 2, 3}
	_, err := computeSineSum("complex", data)
	if err == nil {
		t.Fatal("une erreur était attendue pour un type invalide, mais aucune erreur retournée")
	}
}

// TestComputeSineSumEmpty vérifie que la somme est 0 pour un tableau vide.
func TestComputeSineSumEmpty(t *testing.T) {
	resultInt, err := computeSineSum("int", []int{})
	if err != nil {
		t.Fatalf("erreur inattendue (int) : %v", err)
	}
	if resultInt != 0 {
		t.Errorf("attendu 0 pour tableau vide d'entiers, obtenu %f", resultInt)
	}

	resultFloat, err := computeSineSum("float", []float64{})
	if err != nil {
		t.Fatalf("erreur inattendue (float) : %v", err)
	}
	if resultFloat != 0 {
		t.Errorf("attendu 0 pour tableau vide de flottants, obtenu %f", resultFloat)
	}
}

// ========== Benchmarks ==========

// Pourcentages du tableau à tester
var percentages = []struct {
	name    string
	percent float64
}{
	{"1pct", 0.01},
	{"10pct", 0.10},
	{"20pct", 0.20},
	{"30pct", 0.30},
	{"40pct", 0.40},
	{"50pct", 0.50},
	{"60pct", 0.60},
	{"70pct", 0.70},
	{"80pct", 0.80},
	{"90pct", 0.90},
	{"100pct", 1.00},
}

// Tableaux pré-générés pour les benchmarks (évite de régénérer à chaque itération)
var benchIntArray = generateIntArray(arraySize)
var benchFloatArray = generateFloatArray(arraySize)

// BenchmarkSineSumInt mesure le temps de computeSineSumInt pour différentes tailles de tableau.
func BenchmarkSineSumInt(b *testing.B) {
	for _, p := range percentages {
		size := int(float64(arraySize) * p.percent)
		slice := benchIntArray[:size]
		b.Run(p.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				computeSineSumInt(slice)
			}
		})
	}
}

// BenchmarkSineSumFloat mesure le temps de computeSineSumFloat pour différentes tailles de tableau.
func BenchmarkSineSumFloat(b *testing.B) {
	for _, p := range percentages {
		size := int(float64(arraySize) * p.percent)
		slice := benchFloatArray[:size]
		b.Run(p.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				computeSineSumFloat(slice)
			}
		})
	}
}
