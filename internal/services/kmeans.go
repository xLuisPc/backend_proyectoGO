package services

// Función auxiliar para convertir *float64 en float64 o NaN
import "math"

func nullToNaN(f *float64) float64 {
	if f != nil {
		return *f
	}
	return math.NaN()
}
