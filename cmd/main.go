package main

import (
	"log"

	"github.com/codepnw/go-cart-system/config"
	"github.com/codepnw/go-cart-system/internal/api/middleware"
	"github.com/codepnw/go-cart-system/internal/api/router"
	"github.com/codepnw/go-cart-system/internal/database"
	"github.com/gofiber/fiber/v2"
)

const envPath = "dev.env"

func main() {
	// config
	cfg, err := config.InitEnvConfig(envPath)
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	// database
	db, err := database.NewPostgresDB(cfg.DBAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mid := middleware.NewMiddleware(cfg)

	// router
	routesConfig := &router.RoutesConfig{
		App:       app,
		DB:        db,
		Config:    cfg,
		Middlware: mid,
	}
	routes := router.NewAPIRoutes(routesConfig)
	routes.CartRoutes()
	routes.ProductRoutes()
	routes.UserRoutes()
	routes.OrderRoutes()

	// run server
	if err = app.Listen(cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
