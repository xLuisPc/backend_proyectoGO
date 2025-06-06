package db

import (
	"github.com/xLuisPc/ProyectoGO/internal/models"
	"log"
)

func ObtenerPersonas() ([]models.Persona, error) {
	rows, err := DB.Query(`SELECT 
		id, carrera, genero_accion, genero_ciencia_ficcion, genero_comedia, genero_terror,
		genero_documental, genero_romance, genero_musicales,
		poo, calculo_multivariado, ctd,
		ingenieria_software, bases_datos, control_analogo, circuitos_digitales, promedio 
		FROM dbpersonas`)
	if err != nil {
		log.Println("Error al consultar personas:", err)
		return nil, err
	}
	defer rows.Close()

	var personas []models.Persona
	for rows.Next() {
		var p models.Persona
		err := rows.Scan(
			&p.ID, &p.Carrera,
			&p.GeneroAccion, &p.GeneroCienciaFiccion, &p.GeneroComedia, &p.GeneroTerror,
			&p.GeneroDocumental, &p.GeneroRomance, &p.GeneroMusicales,
			&p.Poo, &p.CalculoMultivariado, &p.Ctd,
			&p.IngenieriaSoftware, &p.BasesDatos, &p.ControlAnalogo, &p.CircuitosDigitales,
			&p.Promedio,
		)
		if err != nil {
			log.Println("Error escaneando persona:", err)
			continue
		}
		personas = append(personas, p)
	}
	return personas, nil
}
