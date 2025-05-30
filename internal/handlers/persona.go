package handlers

import (
	"database/sql"
	"encoding/json"
	"github.com/xLuisPc/ProyectoGO/internal/db"
	"github.com/xLuisPc/ProyectoGO/internal/models"
	"log"
	"net/http"
)

func CrearPersona(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var persona models.Persona
	err := json.NewDecoder(r.Body).Decode(&persona)
	if err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Calcular promedio de campos no nulos
	var suma float64
	var cuenta int

	suma += persona.Poo
	cuenta++
	suma += persona.Ctd
	cuenta++
	suma += persona.CalculoMultivariado
	cuenta++

	if persona.IngenieriaSoftware != nil {
		suma += *persona.IngenieriaSoftware
		cuenta++
	}
	if persona.BasesDatos != nil {
		suma += *persona.BasesDatos
		cuenta++
	}
	if persona.ControlAnalogo != nil {
		suma += *persona.ControlAnalogo
		cuenta++
	}
	if persona.CircuitosDigitales != nil {
		suma += *persona.CircuitosDigitales
		cuenta++
	}
	persona.Promedio = suma / 5

	// Obtener nuevo ID
	var ultimoID int
	err = db.DB.QueryRow("SELECT COALESCE(MAX(id), 0) FROM dbpersonas").Scan(&ultimoID)
	if err != nil {
		log.Println("ERROR OBTENER ID:", err)
		http.Error(w, "Error al obtener el último ID", http.StatusInternalServerError)
		return
	}
	nuevoID := ultimoID + 1

	// Insertar con NullFloat64
	query := `
        INSERT INTO dbpersonas (
            id, carrera, genero_accion, genero_ciencia_ficcion, genero_comedia,
            genero_terror, genero_documental, genero_romance, genero_musicales,
            poo, calculo_multivariado, ctd, ingenieria_software, bases_datos,
            control_analogo, circuitos_digitales, promedio
        ) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)
    `
	_, err = db.DB.Exec(query,
		nuevoID,
		persona.Carrera,
		persona.GeneroAccion,
		persona.GeneroCienciaFiccion,
		persona.GeneroComedia,
		persona.GeneroTerror,
		persona.GeneroDocumental,
		persona.GeneroRomance,
		persona.GeneroMusicales,
		persona.Poo,
		persona.CalculoMultivariado,
		persona.Ctd,
		nullFloat64(persona.IngenieriaSoftware),
		nullFloat64(persona.BasesDatos),
		nullFloat64(persona.ControlAnalogo),
		nullFloat64(persona.CircuitosDigitales),
		persona.Promedio,
	)
	if err != nil {
		log.Println("ERROR INSERT:", err)
		http.Error(w, "Error al insertar en la base de datos", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Persona agregada correctamente"))
}

func ListarPersonas(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	rows, err := db.DB.Query(`SELECT 
        id, carrera, genero_accion, genero_ciencia_ficcion, genero_comedia, genero_terror,
        genero_documental, genero_romance, genero_musicales,
        poo, calculo_multivariado, ctd, ingenieria_software, bases_datos,
        control_analogo, circuitos_digitales, promedio FROM dbpersonas`)
	if err != nil {
		log.Println("ERROR CONSULTA:", err)
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
			&p.GeneroAccion, &p.GeneroCienciaFiccion, &p.GeneroComedia, &p.GeneroTerror,
			&p.GeneroDocumental, &p.GeneroRomance, &p.GeneroMusicales,
			&p.Poo, &p.CalculoMultivariado, &p.Ctd,
			&ingSoft, &bases, &analogo, &digitales,
			&p.Promedio,
		)
		if err != nil {
			log.Println("ERROR SCAN:", err)
			http.Error(w, "Error al leer resultados", http.StatusInternalServerError)
			return
		}

		// Asignar si es válido
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(personas)
}

// Función auxiliar para convertir *float64 en sql.NullFloat64
func nullFloat64(val *float64) sql.NullFloat64 {
	if val != nil {
		return sql.NullFloat64{Float64: *val, Valid: true}
	}
	return sql.NullFloat64{Valid: false}
}
