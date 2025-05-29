package main

import (
	"github.com/xLuisPc/ProyectoGO/internal/db"
	"github.com/xLuisPc/ProyectoGO/internal/server"
	"log"
)

func main() {
	db.ConnectDB()
	log.Println("🚀 Base de datos conectada correctamente")

	router := server.SetupRouter()
	log.Println("Servidor corriendo en http://localhost:8080")
	log.Fatal(router.ListenAndServe())
}
