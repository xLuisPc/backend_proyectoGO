package handlers

import (
	"encoding/json"
	"github.com/xLuisPc/ProyectoGO/internal/db"
	"github.com/xLuisPc/ProyectoGO/internal/utils"
	"log"
	"math"
	"net/http"
	"strings"
)

type CorrelacionResponse struct {
	Labels []string    `json:"labels"`
	Matrix [][]float64 `json:"matrix"`
}

var columnasDisponibles = []string{
	"genero_accion", "genero_ciencia_ficcion", "genero_comedia", "genero_terror",
	"genero_documental", "genero_romance", "genero_musicales",
	"poo", "calculo_multivariado", "promedio",
}

func ObtenerCorrelacion(w http.ResponseWriter, r *http.Request) {
	if utils.EnableCORS(w, r) {
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Leer las variables desde el query param
	seleccionadas := columnasDisponibles
	query := r.URL.Query().Get("vars")
	if query != "" {
		seleccionadas = []string{}
		for _, v := range strings.Split(query, ",") {
			v = strings.TrimSpace(v)
			for _, col := range columnasDisponibles {
				if v == col {
					seleccionadas = append(seleccionadas, v)
				}
			}
		}
	}

	// Construir query SQL solo con las columnas seleccionadas
	sql := "SELECT " + strings.Join(seleccionadas, ", ") + " FROM dbpersonas"

	rows, err := db.DB.Query(sql)
	if err != nil {
		log.Println("ERROR QUERY:", err)
		http.Error(w, "Error consultando columnas seleccionadas", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var datos [][]float64
	for rows.Next() {
		values := make([]interface{}, len(seleccionadas))
		floatRefs := make([]float64, len(seleccionadas))
		for i := range floatRefs {
			values[i] = &floatRefs[i]
		}
		if err := rows.Scan(values...); err != nil {
			continue
		}
		valid := true
		for _, v := range floatRefs {
			if v == -1 {
				valid = false
				break
			}
		}
		if valid {
			datos = append(datos, floatRefs)
		}
	}

	if len(datos) == 0 {
		http.Error(w, "Sin datos válidos", http.StatusInternalServerError)
		return
	}

	matrix := calcularMatrizCorrelacion(datos)
	respuesta := CorrelacionResponse{
		Labels: seleccionadas,
		Matrix: matrix,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respuesta)
}

func calcularMatrizCorrelacion(data [][]float64) [][]float64 {
	n := len(data)
	m := len(data[0])
	matrix := make([][]float64, m)

	// Promedios
	prom := make([]float64, m)
	for i := range prom {
		for j := 0; j < n; j++ {
			prom[i] += data[j][i]
		}
		prom[i] /= float64(n)
	}

	// Correlación
	for i := 0; i < m; i++ {
		matrix[i] = make([]float64, m)
		for j := 0; j < m; j++ {
			var num, denA, denB float64
			for k := 0; k < n; k++ {
				a := data[k][i] - prom[i]
				b := data[k][j] - prom[j]
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
