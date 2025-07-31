package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"

	"go_bookstore_backend/config"
	"go_bookstore_backend/models"
	"go_bookstore_backend/routes"

	_ "go_bookstore_backend/docs" // สำหรับ Swagger
)

// @title Go Bookstore API
// @version 1.0
// @description A RESTful API for managing a bookstore
// @host localhost:3000
// @BasePath /

func main() {
	config.LoadEnv()
	config.ConnectDB()

	// AutoMigrate
	models.MigrateUsers(config.DB)
	models.MigrateBooks(config.DB)
	models.MigrateCarts(config.DB)

	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())

	app.Get("/swagger/*", swagger.HandlerDefault)
	routes.SetupRoutes(app)

	app.Listen(":3000")
}
