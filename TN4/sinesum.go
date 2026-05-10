package main

import (
	"errors"
	"flag"
	"fmt"
	"math"
	"math/rand"
	"time"
)

const arraySize = 1_000_000

// generateIntArray génère un tableau de n entiers aléatoires dans [0, 1000]
// en utilisant une source reproductible (seed 42).
func generateIntArray(n int) []int {
	r := rand.New(rand.NewSource(42))
	arr := make([]int, n)
	for i := range arr {
		arr[i] = r.Intn(1001)
	}
	return arr
}

// generateFloatArray génère un tableau de n nombres à virgule flottante dans [0, 1)
// en utilisant une source reproductible (seed 42).
func generateFloatArray(n int) []float64 {
	r := rand.New(rand.NewSource(42))
	arr := make([]float64, n)
	for i := range arr {
		arr[i] = r.Float64()
	}
	return arr
}

// computeSineSumInt calcule la somme des sinus des éléments d'un tableau d'entiers.
func computeSineSumInt(data []int) float64 {
	var sum float64
	for _, v := range data {
		sum += math.Sin(float64(v))
	}
	return sum
}

// computeSineSumFloat calcule la somme des sinus des éléments d'un tableau de flottants.
func computeSineSumFloat(data []float64) float64 {
	var sum float64
	for _, v := range data {
		sum += math.Sin(v)
	}
	return sum
}

// computeSineSum dispatche le calcul vers la fonction typée appropriée.
// Le dispatch via interface{} et assertion de type introduit un surcoût
// négligeable (~1 ns) face à math.Sin (~30 ns), ce qui le rend acceptable
// pour l'interface utilisateur et les benchmarks (cf. Ch. 6).
func computeSineSum(dataType string, data interface{}) (float64, error) {
	switch dataType {
	case "int":
		arr, ok := data.([]int)
		if !ok {
			return 0, errors.New("données invalides : attendu []int pour type \"int\"")
		}
		return computeSineSumInt(arr), nil
	case "float":
		arr, ok := data.([]float64)
		if !ok {
			return 0, errors.New("données invalides : attendu []float64 pour type \"float\"")
		}
		return computeSineSumFloat(arr), nil
	default:
		return 0, fmt.Errorf("type de données invalide : %q (attendu \"int\" ou \"float\")", dataType)
	}
}

// RunResult contient le résultat numérique et les durées mesurées par run.
type RunResult struct {
	Result   float64
	GenTime  time.Duration
	CalcTime time.Duration
}

// run contient la logique pure du programme : génération du tableau et calcul
// de la somme des sinus. Elle ne produit aucun affichage — l'affichage est
// entièrement délégué à main, conformément à la séparation logique/affichage.
func run(dataType string) (RunResult, error) {
	var res RunResult

	switch dataType {
	case "int":
		start := time.Now()
		data := generateIntArray(arraySize)
		res.GenTime = time.Since(start)
		calcStart := time.Now()
		res.Result = computeSineSumInt(data)
		res.CalcTime = time.Since(calcStart)
	case "float":
		start := time.Now()
		data := generateFloatArray(arraySize)
		res.GenTime = time.Since(start)
		calcStart := time.Now()
		res.Result = computeSineSumFloat(data)
		res.CalcTime = time.Since(calcStart)
	default:
		return RunResult{}, fmt.Errorf("type invalide : %q (utilisez \"int\" ou \"float\")", dataType)
	}
	return res, nil
}

func main() {
	dataType := flag.String("type", "float", "Type de données : \"int\" ou \"float\"")
	flag.Parse()

	fmt.Printf("=== Somme des sinus — type=%s, taille=%d ===\n\n", *dataType, arraySize)
	res, err := run(*dataType)
	if err != nil {
		fmt.Printf("Erreur : %v\n", err)
		return
	}
	// Les mesures via time.Since servent uniquement à l'affichage console.
	// L'analyse de performance repose sur testing.B (cf. Ch. 6).
	fmt.Printf("Génération du tableau : %v\n", res.GenTime)
	fmt.Printf("Calcul de la somme    : %v\n", res.CalcTime)
	fmt.Printf("\nRésultat              : %.6f\n", res.Result)
	fmt.Printf("Temps total           : %v\n", res.GenTime+res.CalcTime)
}
