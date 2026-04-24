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

// computeSineSum calcule la somme des sinus pour un tableau de type "int" ou "float".
// Retourne la somme et une erreur si le type est invalide ou si data n'est pas du bon type.
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

func main() {
	dataType := flag.String("type", "float", "Type de données : \"int\" ou \"float\"")
	flag.Parse()

	fmt.Printf("=== Somme des sinus — type=%s, taille=%d ===\n\n", *dataType, arraySize)

	start := time.Now()

	var result float64
	var err error

	switch *dataType {
	case "int":
		data := generateIntArray(arraySize)
		genTime := time.Since(start)
		fmt.Printf("Génération du tableau : %v\n", genTime)

		calcStart := time.Now()
		result, err = computeSineSum("int", data)
		calcTime := time.Since(calcStart)
		fmt.Printf("Calcul de la somme    : %v\n", calcTime)

	case "float":
		data := generateFloatArray(arraySize)
		genTime := time.Since(start)
		fmt.Printf("Génération du tableau : %v\n", genTime)

		calcStart := time.Now()
		result, err = computeSineSum("float", data)
		calcTime := time.Since(calcStart)
		fmt.Printf("Calcul de la somme    : %v\n", calcTime)

	default:
		fmt.Printf("Type invalide : %q. Utilisez \"int\" ou \"float\".\n", *dataType)
		return
	}

	if err != nil {
		fmt.Printf("Erreur : %v\n", err)
		return
	}

	totalTime := time.Since(start)
	fmt.Printf("\nRésultat              : %.6f\n", result)
	fmt.Printf("Temps total           : %v\n", totalTime)
}
