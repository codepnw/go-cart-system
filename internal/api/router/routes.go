package router

import (
	"database/sql"

	"github.com/codepnw/go-cart-system/config"
	"github.com/codepnw/go-cart-system/internal/api/handler"
	"github.com/codepnw/go-cart-system/internal/api/middleware"
	"github.com/codepnw/go-cart-system/internal/repository"
	"github.com/codepnw/go-cart-system/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type RoutesConfig struct {
	App       *fiber.App
	DB        *sql.DB
	Config    *config.EnvConfig
	Middlware *middleware.Middleware
}

func NewAPIRoutes(cfg *RoutesConfig) *RoutesConfig {
	return &RoutesConfig{
		App:       cfg.App,
		DB:        cfg.DB,
		Config:    cfg.Config,
		Middlware: cfg.Middlware,
	}
}

func (r *RoutesConfig) CartRoutes() {
	routes := r.App.Group("/cart", r.Middlware.Authorize())

	repo := repository.NewCartRepository(r.DB)
	uc := usecase.NewCartUsecase(repo)
	hdl := handler.NewCartHandler(uc)

	routes.Post("/", hdl.AddItems)
	routes.Get("/:cartID", hdl.GetCart)
	routes.Patch("/", hdl.UpdateQuantity)
	routes.Delete("/:id", hdl.DeleteItem)
}

func (r *RoutesConfig) ProductRoutes() {
	routes := r.App.Group("/products", r.Middlware.Authorize())

	repo := repository.NewProductRepository(r.DB)
	uc := usecase.NewProductUsecase(repo)
	hdl := handler.NewProductHandler(uc)

	routes.Post("/", hdl.CreateProduct)
	routes.Get("/:id", hdl.GetProduct)
	routes.Get("/", hdl.ListProducts)
	routes.Patch("/:id", hdl.UpdateProduct)
	routes.Delete("/:id", hdl.DeleteProduct)
}

func (r *RoutesConfig) UserRoutes() {
	repo := repository.NewUserRepository(r.DB)
	uc := usecase.NewUserUsecase(repo, *r.Config)
	hdl := handler.NewUserHandler(uc)

	// Public
	pub := r.App.Group("/auth")
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
