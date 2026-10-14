package main

import (
	"log"

	_ "shop-api/docs" // сгенерированный swag

	"shop-api/internal/di"
)

// @title           Shop API
// @version         1.0
// @description     REST API магазина бытовой техники
// @termsOfService  http://swagger.io/terms/
// @contact.name    Shop API Support
// @host            localhost:8080
// @BasePath        /api/v1
// @schemes         http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT-токен в формате "Bearer <token>"
func main() {
	app, err := di.InitializeApp()
	if err != nil {
		log.Fatalf("init app: %v", err)
	}
	defer app.Close()

	if err := app.Serve(); err != nil {
		log.Fatalf("serve: %v", err)
	}
}
