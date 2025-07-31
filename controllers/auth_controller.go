package controllers

import (
	"go_bookstore_backend/config"
	"go_bookstore_backend/models"
	"go_bookstore_backend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
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
	claims := c.Locals("user").(jwt.MapClaims)

	// Extract ID
	id := uint(claims["id"].(float64))

	// Query from DB
	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(fiber.Map{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
	})
}
