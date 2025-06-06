package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/xLuisPc/ProyectoGO/internal/models"
	_ "github.com/xLuisPc/ProyectoGO/internal/models"
)

// creando rama
// CreateTable crea la tabla dbpersonas permitiendo campos que pueden ser NULL.
func CreateTable(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS dbpersonas (
		id SERIAL PRIMARY KEY,
		carrera TEXT NOT NULL,
		genero_accion INTEGER NOT NULL,
		genero_ciencia_ficcion INTEGER NOT NULL,
		genero_comedia INTEGER NOT NULL,
		genero_terror INTEGER NOT NULL,
		genero_documental INTEGER NOT NULL,
		genero_romance INTEGER NOT NULL,
		genero_musicales INTEGER NOT NULL,
		poo REAL NOT NULL,
		ctd REAL NOT NULL,
		calculo_multivariado REAL NOT NULL,
		ingenieria_software REAL NULL,
		bases_datos REAL NULL,
		control_analogo REAL NULL,
		circuitos_digitales REAL NULL,
		promedio REAL NOT NULL
	);
	`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("❌ Error creando la tabla dbpersonas: %v", err)
	} else {
		log.Println("✅ Tabla dbpersonas creada exitosamente.")
	}
}

// DropTableSistemas elimina la tabla dbsistemas si existe.
func DropTableSistemas(conn *sql.DB) {
	_, err := conn.Exec("DROP TABLE IF EXISTS dbsistemas;")
	if err != nil {
		log.Fatalf("❌ Error eliminando la tabla dbsistemas: %v", err)
	} else {
		log.Println("✅ Tabla dbsistemas eliminada exitosamente.")
	}
}

// DropTableElectronica elimina la tabla dbelectronica si existe.
func DropTableElectronica(conn *sql.DB) {
	_, err := conn.Exec("DROP TABLE IF EXISTS dbelectronica;")
	if err != nil {
		log.Fatalf("❌ Error eliminando la tabla dbelectronica: %v", err)
	} else {
		log.Println("✅ Tabla dbelectronica eliminada exitosamente.")
	}
}

// Agregar200Personas lee un archivo JSON y agrega 200 personas a la base de datos.
func Agregar200Personas(db *sql.DB, jsonPath string) error {
	file, err := os.Open(jsonPath)
	if err != nil {
		log.Fatalf("❌ No se pudo abrir el archivo JSON: %v", err)
	}
	defer file.Close()
	log.Println("✅ Archivo JSON abierto correctamente")

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("error leyendo archivo JSON: %w", err)
	}

	var personas []models.Persona
	if err := json.Unmarshal(bytes, &personas); err != nil {
		return fmt.Errorf("error parseando JSON: %w", err)
	}

	for _, p := range personas {
		_, err := db.Exec(`INSERT INTO dbpersonas (
				id, carrera, genero_accion, genero_ciencia_ficcion, genero_comedia,
				genero_terror, genero_documental, genero_romance, genero_musicales,
				poo, ctd, calculo_multivariado, ingenieria_software, bases_datos,
				control_analogo, circuitos_digitales, promedio
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9,
				$10, $11, $12, $13, $14, $15, $16, $17
			)`,
			p.ID, p.Carrera, p.GeneroAccion, p.GeneroCienciaFiccion, p.GeneroComedia,
			p.GeneroTerror, p.GeneroDocumental, p.GeneroRomance, p.GeneroMusicales,
			p.Poo, p.Ctd, p.CalculoMultivariado,
			nullFloat64(p.IngenieriaSoftware),
			nullFloat64(p.BasesDatos),
			nullFloat64(p.ControlAnalogo),
			nullFloat64(p.CircuitosDigitales),
			p.Promedio)

		if err != nil {
			log.Printf("❌ Error insertando ID %d: %v", p.ID, err)
		}
	}

	fmt.Println("✅ Las 200 personas fueron insertadas exitosamente.")
	return nil
}

// nullFloat64 convierte *float64 en sql.NullFloat64
func nullFloat64(val *float64) sql.NullFloat64 {
	if val != nil {
		return sql.NullFloat64{Float64: *val, Valid: true}
	}
	return sql.NullFloat64{Valid: false}
}
