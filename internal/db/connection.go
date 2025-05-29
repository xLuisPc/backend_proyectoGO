package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {
	// Obtener variables desde el entorno
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Verificar que ninguna variable esté vacía
	if host == "" || port == "" || user == "" || password == "" || dbname == "" {
		log.Fatal("❌ Faltan variables de entorno para la conexión a la base de datos")
	}

	// Formatear la cadena de conexión
	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	// Conectar a PostgreSQL
	var err error
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("❌ Error al abrir la base de datos: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("❌ No se pudo conectar a la base de datos: %v", err)
	}

	log.Println("✅ Conexión exitosa a PostgreSQL")
}
