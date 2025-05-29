package main

import (
	"log"
	"net/http"
	"os"

	"github.com/xLuisPc/ProyectoGO/internal/db"
	"github.com/xLuisPc/ProyectoGO/internal/server"
)

func main() {
	// Conexión a la base de datos
	db.ConnectDB()
	log.Println("🚀 Base de datos conectada correctamente")

	// Configuración del router (debe devolver http.Handler)
	router := server.SetupRouter()

	// Obtener puerto desde variable de entorno (requerido por Render)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback local
	}

	log.Printf("Servidor corriendo en http://localhost:%s\n", port)

	// Iniciar servidor
	log.Fatal(http.ListenAndServe(":"+port, router))
}
