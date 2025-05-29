package handlers

import (
	"database/sql"
	"encoding/json"
	"github.com/xLuisPc/ProyectoGO/internal/db"
	"github.com/xLuisPc/ProyectoGO/internal/models"
	"github.com/xLuisPc/ProyectoGO/internal/services"
	"log"
	"net/http"
	"strconv"
)

func ObtenerClusters(w http.ResponseWriter, r *http.Request) {
	genero := r.URL.Query().Get("genero")
	if genero == "" {
		http.Error(w, "Parámetro 'genero' requerido", http.StatusBadRequest)
		return
	}

	k := 3 // valor por defecto
	if kStr := r.URL.Query().Get("k"); kStr != "" {
		if parsed, err := strconv.Atoi(kStr); err == nil && parsed >= 2 && parsed <= 10 {
			k = parsed
		} else {
			log.Println("⚠️ Valor inválido de k, usando 3")
		}
	}

	log.Println("✅ Generando clusters con género:", genero, "y K =", k)

	rows, err := db.DB.Query(`SELECT 
		id, carrera, genero_accion, genero_ciencia_ficcion, genero_comedia,
		genero_terror, genero_documental, genero_romance, genero_musicales,
		poo, calculo_multivariado, ctd, 
		ingenieria_software, bases_datos, control_analogo, circuitos_digitales,
		promedio 
		FROM dbpersonas`)
	if err != nil {
		log.Println("❌ Error al obtener personas:", err)
		http.Error(w, "Error al consultar la base de datos", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var personas []models.Persona
	for rows.Next() {
		var p models.Persona
		var ingSoft, bases, analogo, digitales sql.NullFloat64

		err := rows.Scan(
			&p.ID, &p.Carrera,
			&p.GeneroAccion, &p.GeneroCienciaFiccion, &p.GeneroComedia,
			&p.GeneroTerror, &p.GeneroDocumental, &p.GeneroRomance, &p.GeneroMusicales,
			&p.Poo, &p.CalculoMultivariado, &p.Ctd,
			&ingSoft, &bases, &analogo, &digitales,
			&p.Promedio,
		)
		if err != nil {
			log.Println("❌ Error al escanear persona:", err)
			http.Error(w, "Error al procesar resultados", http.StatusInternalServerError)
			return
		}

		// Convertir NullFloat64 a punteros
		if ingSoft.Valid {
			p.IngenieriaSoftware = &ingSoft.Float64
		}
		if bases.Valid {
			p.BasesDatos = &bases.Float64
		}
		if analogo.Valid {
			p.ControlAnalogo = &analogo.Float64
		}
		if digitales.Valid {
			p.CircuitosDigitales = &digitales.Float64
		}

		personas = append(personas, p)
	}

	clusters := services.KMeansPorGenero(personas, genero, k)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clusters)
}
