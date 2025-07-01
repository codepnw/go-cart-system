package main

import (
	"log"

	"github.com/codepnw/go-cart-system/config"
	"github.com/codepnw/go-cart-system/internal/api/router"
	"github.com/codepnw/go-cart-system/internal/database"
	"github.com/gofiber/fiber/v2"
)

const envPath = "dev.env"

func main() {
	// config
	config, err := config.InitEnvConfig(envPath)
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	// database
	db, err := database.NewPostgresDB(config.DBAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// router
	routes := router.NewAPIRoutes(app, db, config)
	routes.CartRoutes()
	routes.ProductRoutes()
	routes.UserRoutes()

	// run server
	if err = app.Listen(config.AppPort); err != nil {
		log.Fatal(err)
	}
}
