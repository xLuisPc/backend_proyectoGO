package handlers

import (
	"encoding/json"
	"github.com/xLuisPc/ProyectoGO/internal/db"
	"github.com/xLuisPc/ProyectoGO/internal/models"
	"github.com/xLuisPc/ProyectoGO/internal/services"
	"log"
	"math"
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
		SELECT 
			genero_accion, genero_ciencia_ficcion, genero_comedia,
			genero_terror, genero_documental, genero_romance,
			genero_musicales, poo, calculo_multivariado, promedio
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
		var poo, calc, promedio float64

		if err := rows.Scan(&g1, &g2, &g3, &g4, &g5, &g6, &g7, &poo, &calc, &promedio); err != nil {
			continue
		}

		vec := []float64{
			float64(g1), float64(g2), float64(g3),
			float64(g4), float64(g5), float64(g6), float64(g7),
			poo, calc,
		}
		dataset = append(dataset, vec)
		promedios = append(promedios, promedio)
	}

	// Construir vector del nuevo perfil
	vecNuevo := []float64{
		float64(nuevo.GeneroAccion),
		float64(nuevo.GeneroCienciaFiccion),
		float64(nuevo.GeneroComedia),
		float64(nuevo.GeneroTerror),
		float64(nuevo.GeneroDocumental),
		float64(nuevo.GeneroRomance),
		float64(nuevo.GeneroMusicales),
		-1, // por defecto sin notas
		-1,
	}

	if nuevo.Poo >= 0 && nuevo.Poo <= 5 {
		vecNuevo[7] = nuevo.Poo
	}
	if nuevo.CalculoMultivariado >= 0 && nuevo.CalculoMultivariado <= 5 {
		vecNuevo[8] = nuevo.CalculoMultivariado
	}

	// Llamar a KNN
	promedioEstimado := services.KNNPredecirPromedio(dataset, promedios, vecNuevo, 5)

	// Respuesta final
	response := map[string]interface{}{
		"promedio": math.Round(promedioEstimado*100) / 100, // redondear a 2 decimales
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
