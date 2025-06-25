package server

import (
	"log"

	"github.com/codepnw/go-cart-system/internal/database"
	"github.com/gofiber/fiber/v2"
)

type ServerConfig struct {
	DB_ADDR  string
	APP_PORT string
}

func NewServer(config ServerConfig) {
	app := fiber.New()

	db, err := database.NewPostgresDB(config.DB_ADDR)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err = app.Listen(":" + config.APP_PORT); err != nil {
		log.Fatal(err)
	}
}
