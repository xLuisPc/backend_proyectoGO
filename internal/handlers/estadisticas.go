package handlers

import (
	"encoding/json"
	"github.com/xLuisPc/ProyectoGO/internal/db"
	"github.com/xLuisPc/ProyectoGO/internal/services"
	"github.com/xLuisPc/ProyectoGO/internal/utils"
	"log"
	"net/http"
	"strconv"
)

func ObtenerClusters(w http.ResponseWriter, r *http.Request) {
	if utils.EnableCORS(w, r) {
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	personas, err := db.ObtenerPersonas()
	if err != nil {
		log.Println("Error obteniendo personas:", err)
		http.Error(w, "Error en la base de datos", http.StatusInternalServerError)
		return
	}

	genero := r.URL.Query().Get("genero")
	kStr := r.URL.Query().Get("k")
	k, err := strconv.Atoi(kStr)
	if err != nil || k < 2 || k > 10 {
		k = 3
	}

	clusters := services.KMeansPorGenero(personas, genero, k)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clusters)
}
