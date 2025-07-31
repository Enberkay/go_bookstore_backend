package controllers

import (
	"go_bookstore_backend/config"
	"go_bookstore_backend/models"
	"go_bookstore_backend/utils"

	"github.com/gofiber/fiber/v2"
)

func Register(c *fiber.Ctx) error {
	var input models.User
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	hash, _ := utils.HashPassword(input.Password)
	user := models.User{
		Email:    input.Email,
		Password: hash,
		Role:     "user",
	}
	config.DB.Create(&user)
	return c.JSON(user)
}
