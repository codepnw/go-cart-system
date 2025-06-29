package router

import (
	"database/sql"

	"github.com/codepnw/go-cart-system/internal/api/handler"
	"github.com/codepnw/go-cart-system/internal/repository"
	"github.com/codepnw/go-cart-system/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, db *sql.DB) {
	cartRoutes(app, db)
}

func cartRoutes(app *fiber.App, db *sql.DB) {
	routes := app.Group("/cart")

	repo := repository.NewCartRepository(db)
	uc := usecase.NewCartUsecase(repo)
	hdl := handler.NewCartHandler(uc)

	routes.Post("/", hdl.AddItems)
	routes.Get("/:cartID", hdl.GetCart)
	routes.Patch("/", hdl.UpdateQuantity)
	routes.Delete("/:id", hdl.DeleteItem)
}
