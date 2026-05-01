package main

import (
	"flag"
	"math"
	"os"
	"testing"
)

const tolerance = 1e-9

// ========== Tests unitaires ==========

// TestComputeSineSumInt vérifie le calcul pour [1, 2, 3] en comparant au résultat
// attendu avec une tolérance de 1e-9. Important car la conversion int → float64
// (instruction CVTSI2SD) pourrait introduire une erreur si mal implémentée.
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

// TestComputeSineSumFloat vérifie le calcul pour [0.1, 0.2, 0.3] sans conversion
// de type. Valide que math.Sin reçoit directement des float64 et que l'accumulation
// par addition ne dépasse pas la tolérance de 1e-9 sur 3 éléments.
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

// TestComputeSineSumInvalidType passe "complex" comme type et vérifie que la
// fonction retourne une erreur. Sans ce test, un type non supporté pourrait
// silencieusement retourner 0 au lieu de signaler le problème.
func TestComputeSineSumInvalidType(t *testing.T) {
	data := []int{1, 2, 3}
	_, err := computeSineSum("complex", data)
	if err == nil {
		t.Fatal("une erreur était attendue pour un type invalide, mais aucune erreur retournée")
	}
}

// TestComputeSineSumEmpty vérifie que la somme est 0 pour un tableau vide,
// tant pour les entiers que pour les flottants. Cas limite qui garantit
// l'absence de panique ou d'index hors bornes sur une slice de taille 0.
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

// TestComputeSineSumNegative vérifie que les entiers négatifs sont gérés
// correctement. math.Sin est une fonction impaire, donc sin(-1) = -sin(1),
// et la somme sin(-1)+sin(0)+sin(1) doit valoir exactement sin(0) = 0.
func TestComputeSineSumNegative(t *testing.T) {
	data := []int{-1, 0, 1}
	expected := math.Sin(-1) + math.Sin(0) + math.Sin(1)

	result, err := computeSineSum("int", data)
	if err != nil {
		t.Fatalf("erreur inattendue : %v", err)
	}
	if math.Abs(result-expected) > tolerance {
		t.Errorf("résultat incorrect : obtenu %.10f, attendu %.10f", result, expected)
	}
}

// TestComputeSineSumLargeFloat vérifie la stabilité numérique avec des valeurs
// proches de la limite de précision float64. Des valeurs comme 1e15 forcent
// la réduction d'argument de math.Sin à travailler sur de grands multiples de pi,
// ce qui peut amplifier les erreurs d'arrondi.
func TestComputeSineSumLargeFloat(t *testing.T) {
	data := []float64{1e15, 1e-15, 0.0}
	expected := math.Sin(1e15) + math.Sin(1e-15) + math.Sin(0.0)

	result, err := computeSineSum("float", data)
	if err != nil {
		t.Fatalf("erreur inattendue : %v", err)
	}
	if math.Abs(result-expected) > tolerance {
		t.Errorf("résultat incorrect : obtenu %.10f, attendu %.10f", result, expected)
	}
}

// TestComputeSineSumIntWrongData vérifie que passer des données du mauvais type
// ([]float64 au lieu de []int) retourne une erreur au lieu de paniquer.
func TestComputeSineSumIntWrongData(t *testing.T) {
	_, err := computeSineSum("int", []float64{0.1})
	if err == nil {
		t.Fatal("attendu une erreur pour type \"int\" avec []float64")
	}
}

// TestComputeSineSumFloatWrongData vérifie que passer des données du mauvais type
// ([]int au lieu de []float64) retourne une erreur.
func TestComputeSineSumFloatWrongData(t *testing.T) {
	_, err := computeSineSum("float", []int{1})
	if err == nil {
		t.Fatal("attendu une erreur pour type \"float\" avec []int")
	}
}

// TestRunInt vérifie que la fonction run s'exécute correctement avec le type int.
func TestRunInt(t *testing.T) {
	_, err := run("int")
	if err != nil {
		t.Fatalf("erreur inattendue : %v", err)
	}
}

// TestRunFloat vérifie que la fonction run s'exécute correctement avec le type float.
func TestRunFloat(t *testing.T) {
	_, err := run("float")
	if err != nil {
		t.Fatalf("erreur inattendue : %v", err)
	}
}

// TestRunInvalidType vérifie que run retourne une erreur pour un type invalide.
func TestRunInvalidType(t *testing.T) {
	_, err := run("complex")
	if err == nil {
		t.Fatal("attendu une erreur pour un type invalide")
	}
}

// TestMainFunction vérifie que main() s'exécute sans panique pour le cas nominal.
// On réinitialise flag.CommandLine car main() appelle flag.String et flag.Parse.
func TestMainFunction(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	os.Args = []string{"cmd", "-type", "float"}
	main()
}

// TestMainFunctionError vérifie que main() gère un type invalide sans panique.
func TestMainFunctionError(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	os.Args = []string{"cmd", "-type", "complex"}
	main()
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

// BenchmarkSineSumInt mesure le temps de computeSineSum pour des entiers
// à différentes tailles de tableau. On passe par le dispatch pour benchmarker
// la fonction telle qu'elle est appelée en production (cf. Ch. 6).
func BenchmarkSineSumInt(b *testing.B) {
	b.ReportAllocs()
	for _, p := range percentages {
		size := int(float64(arraySize) * p.percent)
		slice := benchIntArray[:size]
		b.Run(p.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				computeSineSum("int", slice)
			}
		})
	}
}

// BenchmarkSineSumFloat mesure le temps de computeSineSum pour des flottants
// à différentes tailles de tableau. Le dispatch via interface{} ajoute un coût
// négligeable par rapport à math.Sin, mais on le mesure quand même (cf. Ch. 6).
func BenchmarkSineSumFloat(b *testing.B) {
	b.ReportAllocs()
	for _, p := range percentages {
		size := int(float64(arraySize) * p.percent)
		slice := benchFloatArray[:size]
		b.Run(p.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				computeSineSum("float", slice)
			}
		})
	}
}

// BenchmarkSineSumProfile est un benchmark dédié au profilage CPU (pprof).
// Il appelle computeSineSumFloat directement (sans dispatch) sur le tableau
// complet pour maximiser le temps passé dans math.Sin et obtenir un profil
// statistiquement significatif. Lancer avec :
//
//	go test -bench=BenchmarkSineSumProfile -cpuprofile=cpu.prof -run=^$
//	go tool pprof -png cpu.prof > cpu-profile.png
//	go tool pprof -http=:8080 cpu.prof   (flamegraph interactif)
func BenchmarkSineSumProfile(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		computeSineSumFloat(benchFloatArray)
	}
}
