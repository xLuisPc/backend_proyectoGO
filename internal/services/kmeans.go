package services

import (
	"github.com/xLuisPc/ProyectoGO/internal/models"
	"math"
	"math/rand"
	"time"
)

type Cluster struct {
	ID       int
	Personas []models.Persona
}

// KMeansPorGenero agrupa estudiantes según afinidad a un género + notas, ignorando campos nulos
func KMeansPorGenero(personas []models.Persona, genero string, k int) []Cluster {
	rand.Seed(time.Now().UnixNano())

	var dataset [][]float64
	for _, p := range personas {
		var afinidad float64
		switch genero {
		case "genero_accion":
			afinidad = float64(p.GeneroAccion)
		case "genero_ciencia_ficcion":
			afinidad = float64(p.GeneroCienciaFiccion)
		case "genero_comedia":
			afinidad = float64(p.GeneroComedia)
		case "genero_terror":
			afinidad = float64(p.GeneroTerror)
		case "genero_documental":
			afinidad = float64(p.GeneroDocumental)
		case "genero_romance":
			afinidad = float64(p.GeneroRomance)
		case "genero_musicales":
			afinidad = float64(p.GeneroMusicales)
		default:
			afinidad = 0
		}

		entry := []float64{
			afinidad,
			p.Poo,
			p.CalculoMultivariado,
			p.Ctd,
			nullToNaN(p.IngenieriaSoftware),
			nullToNaN(p.BasesDatos),
			nullToNaN(p.ControlAnalogo),
			nullToNaN(p.CircuitosDigitales),
			p.Promedio,
		}
		dataset = append(dataset, entry)
	}

	centroids := make([][]float64, k)
	for i := 0; i < k; i++ {
		centroids[i] = dataset[rand.Intn(len(dataset))]
	}

	assignments := make([]int, len(dataset))
	for iter := 0; iter < 100; iter++ {
		for i, point := range dataset {
			minDist := math.MaxFloat64
			for j, centroid := range centroids {
				if dist := euclideanIgnoreNaN(point, centroid); dist < minDist {
					minDist = dist
					assignments[i] = j
				}
			}
		}

		newCentroids := make([][]float64, k)
		counts := make([]int, k)
		for i := 0; i < k; i++ {
			newCentroids[i] = make([]float64, len(dataset[0]))
		}

		for i, cluster := range assignments {
			for j := range dataset[i] {
				if !math.IsNaN(dataset[i][j]) {
					newCentroids[cluster][j] += dataset[i][j]
				}
			}
			counts[cluster]++
		}

		for i := 0; i < k; i++ {
			for j := range newCentroids[i] {
				if counts[i] > 0 {
					newCentroids[i][j] /= float64(counts[i])
				}
			}
		}
		centroids = newCentroids
	}

	clusters := make([]Cluster, k)
	for i := range clusters {
		clusters[i].ID = i
	}
	for i, idx := range assignments {
		clusters[idx].Personas = append(clusters[idx].Personas, personas[i])
	}

	return clusters
}

func PredecirClusterPorGustosConPromedios(dataset [][]float64, nuevo []float64, k int) (int, []int) {
	rand.Seed(time.Now().UnixNano())

	centroids := make([][]float64, k)
	for i := 0; i < k; i++ {
		centroids[i] = dataset[rand.Intn(len(dataset))]
	}

	assignments := make([]int, len(dataset))
	for iter := 0; iter < 100; iter++ {
		for i, punto := range dataset {
			minDist := math.MaxFloat64
			for j, centro := range centroids {
				if d := euclideanIgnoreNaN(punto, centro); d < minDist {
					minDist = d
					assignments[i] = j
				}
			}
		}

		newCentroids := make([][]float64, k)
		counts := make([]int, k)
		for i := 0; i < k; i++ {
			newCentroids[i] = make([]float64, len(dataset[0]))
		}
		for i, c := range assignments {
			for j := range dataset[i] {
				if !math.IsNaN(dataset[i][j]) {
					newCentroids[c][j] += dataset[i][j]
				}
			}
			counts[c]++
		}
		for i := 0; i < k; i++ {
			for j := range newCentroids[i] {
				if counts[i] > 0 {
					newCentroids[i][j] /= float64(counts[i])
				}
			}
		}
		centroids = newCentroids
	}

	// Predecir cluster para el nuevo vector
	minDist := math.MaxFloat64
	clusterAsignado := 0
	for j, centro := range centroids {
		if d := euclideanIgnoreNaN(nuevo, centro); d < minDist {
			minDist = d
			clusterAsignado = j
		}
	}
	return clusterAsignado, assignments
}

// Utilidades

func nullToNaN(f *float64) float64 {
	if f != nil {
		return *f
	}
	return math.NaN()
}

func euclideanIgnoreNaN(a, b []float64) float64 {
	var sum float64
	count := 0
	for i := range a {
		if math.IsNaN(a[i]) || math.IsNaN(b[i]) {
			continue
		}
		diff := a[i] - b[i]
		sum += diff * diff
		count++
	}
	if count == 0 {
		return math.MaxFloat64
	}
	return math.Sqrt(sum)
}
