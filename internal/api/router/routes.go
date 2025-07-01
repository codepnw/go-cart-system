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
	productRoutes(app, db)
	userRoutes(app, db)
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

func productRoutes(app *fiber.App, db *sql.DB) {
	routes := app.Group("/products")

	repo := repository.NewProductRepository(db)
	uc := usecase.NewProductUsecase(repo)
	hdl := handler.NewProductHandler(uc)

	routes.Post("/", hdl.CreateProduct)
	routes.Get("/:id", hdl.GetProduct)
	routes.Get("/", hdl.ListProducts)
	routes.Patch("/:id", hdl.UpdateProduct)
	routes.Delete("/:id", hdl.DeleteProduct)
}

func userRoutes(app *fiber.App, db *sql.DB) {
	repo := repository.NewUserRepository(db)
	uc := usecase.NewUserUsecase(repo)
	hdl := handler.NewUserHandler(uc)

	// Public
	pub := app.Group("/auth")
	pub.Post("/register", hdl.Register)
	pub.Post("/login", hdl.Login)

	// Private
	// 	get profile
	// 	update profile
	//	logout
	//	refresh token

	// Admin
	//	list users
	//	get user
	//	update role
	//	delete user
}
