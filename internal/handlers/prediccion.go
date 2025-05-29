package handlers

import (
	"database/sql"
	"encoding/json"
	"github.com/xLuisPc/ProyectoGO/internal/db"
	"github.com/xLuisPc/ProyectoGO/internal/models"
	"github.com/xLuisPc/ProyectoGO/internal/services"
	"log"
	"net/http"
)

func PredecirCluster(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var nuevo models.Persona
	err := json.NewDecoder(r.Body).Decode(&nuevo)
	if err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	rows, err := db.DB.Query(`
		SELECT genero_accion, genero_ciencia_ficcion, genero_comedia,
		       genero_terror, genero_documental, genero_romance,
		       genero_musicales, promedio 
		FROM dbpersonas`)
	if err != nil {
		log.Println("❌ Error al obtener personas:", err)
		http.Error(w, "Error al consultar base de datos", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var dataset [][]float64
	var promedios []float64
	for rows.Next() {
		var g1, g2, g3, g4, g5, g6, g7 int
		var promedio float64
		if err := rows.Scan(&g1, &g2, &g3, &g4, &g5, &g6, &g7, &promedio); err != nil {
			continue
		}
		dataset = append(dataset, []float64{float64(g1), float64(g2), float64(g3), float64(g4), float64(g5), float64(g6), float64(g7)})
		promedios = append(promedios, promedio)
	}

	perfilNuevo := []float64{
		float64(nuevo.GeneroAccion),
		float64(nuevo.GeneroCienciaFiccion),
		float64(nuevo.GeneroComedia),
		float64(nuevo.GeneroTerror),
		float64(nuevo.GeneroDocumental),
		float64(nuevo.GeneroRomance),
		float64(nuevo.GeneroMusicales),
	}

	clusterID, asignaciones := services.PredecirClusterPorGustosConPromedios(dataset, perfilNuevo, 3)

	var suma float64
	var cantidad int
	for i, asignado := range asignaciones {
		if asignado == clusterID {
			suma += promedios[i]
			cantidad++
		}
	}
	promedioEstimado := 0.0
	if cantidad > 0 {
		promedioEstimado = suma / float64(cantidad)
	}

	// Calcular notas estimadas por materia
	rows2, _ := db.DB.Query(`
		SELECT poo, ctd, calculo_multivariado,
		       ingenieria_software, bases_datos,
		       control_analogo, circuitos_digitales
		FROM dbpersonas`)
	defer rows2.Close()

	campos := []string{"poo", "ctd", "calculo_multivariado", "ingenieria_software", "bases_datos", "control_analogo", "circuitos_digitales"}
	materias := map[string][]float64{}
	i := 0
	for rows2.Next() {
		var m [7]sql.NullFloat64
		if err := rows2.Scan(&m[0], &m[1], &m[2], &m[3], &m[4], &m[5], &m[6]); err != nil {
			continue
		}
		if i < len(asignaciones) && asignaciones[i] == clusterID {
			for j, key := range campos {
				if m[j].Valid {
					materias[key] = append(materias[key], m[j].Float64)
				}
			}
		}
		i++
	}

	estimadas := map[string]float64{}
	for k, valores := range materias {
		var suma float64
		for _, v := range valores {
			suma += v
		}
		if len(valores) > 0 {
			estimadas[k] = suma / float64(len(valores))
		}
	}

	// Respuesta final
	response := map[string]interface{}{
		"cluster":  clusterID,
		"promedio": promedioEstimado,
		"materias": estimadas,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
