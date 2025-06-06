package handlers

import (
	"encoding/json"
	"github.com/xLuisPc/ProyectoGO/internal/db"
	"log"
	"math"
	"net/http"
)

// CorrelacionResponse define la estructura JSON de salida
type CorrelacionResponse struct {
	Labels []string    `json:"labels"`
	Matrix [][]float64 `json:"matrix"`
}

// Handler para /api/correlacion
func ObtenerCorrelacion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	rows, err := db.DB.Query(`
		SELECT 
			genero_accion, genero_ciencia_ficcion, genero_comedia, genero_terror,
			genero_documental, genero_romance, genero_musicales,
			poo, calculo_multivariado, promedio
		FROM dbpersonas`)
	if err != nil {
		log.Println("ERROR QUERY:", err)
		http.Error(w, "Error al consultar la base de datos", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var datos [][]float64

	for rows.Next() {
		var g1, g2, g3, g4, g5, g6, g7 int
		var poo, calc, promedio float64

		err := rows.Scan(&g1, &g2, &g3, &g4, &g5, &g6, &g7, &poo, &calc, &promedio)
		if err != nil {
			continue
		}

		if poo >= 0 && calc >= 0 {
			vector := []float64{
				float64(g1), float64(g2), float64(g3), float64(g4),
				float64(g5), float64(g6), float64(g7),
				poo, calc, promedio,
			}
			datos = append(datos, vector)
		}
	}

	if len(datos) == 0 {
		http.Error(w, "No hay datos suficientes", http.StatusInternalServerError)
		return
	}

	matrix := calcularMatrizCorrelacion(datos)

	labels := []string{
		"Acción", "Ciencia Ficción", "Comedia", "Terror",
		"Documental", "Romance", "Musicales",
		"POO", "Cálculo", "Promedio",
	}

	respuesta := CorrelacionResponse{
		Labels: labels,
		Matrix: matrix,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respuesta)
}

// Calcula la matriz de correlación de Pearson
func calcularMatrizCorrelacion(data [][]float64) [][]float64 {
	n := len(data)
	m := len(data[0])
	matrix := make([][]float64, m)

	// Calcular promedios
	promedios := make([]float64, m)
	for i := 0; i < m; i++ {
		var sum float64
		for j := 0; j < n; j++ {
			sum += data[j][i]
		}
		promedios[i] = sum / float64(n)
	}

	// Calcular correlaciones
	for i := 0; i < m; i++ {
		matrix[i] = make([]float64, m)
		for j := 0; j < m; j++ {
			num, denA, denB := 0.0, 0.0, 0.0
			for k := 0; k < n; k++ {
				a := data[k][i] - promedios[i]
				b := data[k][j] - promedios[j]
				num += a * b
				denA += a * a
				denB += b * b
			}
			if denA > 0 && denB > 0 {
				matrix[i][j] = num / math.Sqrt(denA*denB)
			} else {
				matrix[i][j] = 0
			}
		}
	}
	return matrix
}
