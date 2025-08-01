package routes

import (
	"go_bookstore_backend/controllers"
	"go_bookstore_backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	// Auth Routes
	auth := api.Group("/auth")
	auth.Post("/register", controllers.Register)
	auth.Post("/login", controllers.Login)
	auth.Get("/me", middlewares.RequireAuth, controllers.CurrentUser)

	// Book Routes
	books := api.Group("/books")
	books.Get("/", controllers.GetBooks)
	books.Get("/:id", controllers.GetBook)
	books.Post("/", middlewares.RequireAuth, middlewares.RequireRole("admin"), controllers.CreateBook)
	books.Put("/:id", middlewares.RequireAuth, middlewares.RequireRole("admin"), controllers.UpdateBook)
	books.Delete("/:id", middlewares.RequireAuth, middlewares.RequireRole("admin"), controllers.DeleteBook)

	// Cart Routes
	cart := api.Group("/cart", middlewares.RequireAuth)
	cart.Get("/", controllers.ViewCart)
	cart.Post("/", controllers.AddToCart)
	cart.Delete("/:id", controllers.RemoveFromCart)

}
