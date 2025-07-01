package router

import (
	"database/sql"

	"github.com/codepnw/go-cart-system/config"
	"github.com/codepnw/go-cart-system/internal/api/handler"
	"github.com/codepnw/go-cart-system/internal/repository"
	"github.com/codepnw/go-cart-system/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type routesConfig struct {
	app    *fiber.App
	db     *sql.DB
	config *config.EnvConfig
}

func NewAPIRoutes(app *fiber.App, db *sql.DB, config *config.EnvConfig) *routesConfig {
	return &routesConfig{
		app:    app,
		db:     db,
		config: config,
	}
}

func (r *routesConfig) CartRoutes() {
	routes := r.app.Group("/cart")

	repo := repository.NewCartRepository(r.db)
	uc := usecase.NewCartUsecase(repo)
	hdl := handler.NewCartHandler(uc)

	routes.Post("/", hdl.AddItems)
	routes.Get("/:cartID", hdl.GetCart)
	routes.Patch("/", hdl.UpdateQuantity)
	routes.Delete("/:id", hdl.DeleteItem)
}

func (r *routesConfig) ProductRoutes() {
	routes := r.app.Group("/products")

	repo := repository.NewProductRepository(r.db)
	uc := usecase.NewProductUsecase(repo)
	hdl := handler.NewProductHandler(uc)

	routes.Post("/", hdl.CreateProduct)
	routes.Get("/:id", hdl.GetProduct)
	routes.Get("/", hdl.ListProducts)
	routes.Patch("/:id", hdl.UpdateProduct)
	routes.Delete("/:id", hdl.DeleteProduct)
}

func (r *routesConfig) UserRoutes() {
	repo := repository.NewUserRepository(r.db)
	uc := usecase.NewUserUsecase(repo, *r.config)
	hdl := handler.NewUserHandler(uc)

	// Public
	pub := r.app.Group("/auth")
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
