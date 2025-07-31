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

func Login(c *fiber.Ctx) error {
	var input models.User
	var user models.User

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid email"})
	}
	if err := utils.CheckPassword(user.Password, input.Password); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid password"})
	}

	token, _ := utils.GenerateJWT(user.ID, user.Email, user.Role)
	return c.JSON(fiber.Map{"token": token})
}

func CurrentUser(c *fiber.Ctx) error {
	user := c.Locals("user").(map[string]interface{})
	return c.JSON(user)
}
