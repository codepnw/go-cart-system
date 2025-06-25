package main

import (
	"log"
	"os"

	"github.com/codepnw/go-cart-system/internal/api/server"
	"github.com/joho/godotenv"
)

const envFile = "dev.env"

func main() {
	if err := godotenv.Load(envFile); err != nil {
		log.Fatal(err)
	}

	server.NewServer(server.ServerConfig{
		DB_ADDR: os.Getenv("DB_ADDR"),
		APP_PORT: os.Getenv("APP_PORT"),
	})
}