package services

import (
	"math"
	"sort"
)

// KNNPredecirPromedio calcula el promedio estimado con KNN
func KNNPredecirPromedio(dataset [][]float64, targets []float64, nuevo []float64, k int) float64 {
	type vecino struct {
		distancia float64
		promedio  float64
	}

	var vecinos []vecino

	for i, datos := range dataset {
		d := distanciaIgnorandoValoresInvalidos(datos, nuevo)
		vecinos = append(vecinos, vecino{
			distancia: d,
			promedio:  targets[i],
		})
	}

	// Ordenar por distancia ascendente
	sort.Slice(vecinos, func(i, j int) bool {
		return vecinos[i].distancia < vecinos[j].distancia
	})

	// Tomar los K más cercanos
	k = int(math.Min(float64(k), float64(len(vecinos))))
	var suma float64
	for i := 0; i < k; i++ {
		suma += vecinos[i].promedio
	}

	if k > 0 {
		return suma / float64(k)
	}
	return 0
}

// distanciaIgnorandoValoresInvalidos calcula la distancia euclidiana
// ignorando cualquier valor igual a -1 o NaN
func distanciaIgnorandoValoresInvalidos(a, b []float64) float64 {
	var suma float64
	var count int

	for i := range a {
		if a[i] == -1 || b[i] == -1 || math.IsNaN(a[i]) || math.IsNaN(b[i]) {
			continue
		}
		diff := a[i] - b[i]
		suma += diff * diff
		count++
	}

	if count == 0 {
		return math.MaxFloat64
	}
	return math.Sqrt(suma)
}
