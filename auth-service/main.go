package main

import (
	"fmt"
	"net/http"

	"auth-service/config"
	"auth-service/db"
	"auth-service/routes"
)

func main() {
	config.LoadConfig()
	db.InitDB()
	routes.SetupRoutes()

	fmt.Println("Server running on port 3001")
	http.ListenAndServe(":3001", nil)
}
