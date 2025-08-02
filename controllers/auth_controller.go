package controllers

import (
	"go_bookstore_backend/config"
	"go_bookstore_backend/models"
	"go_bookstore_backend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Register godoc
// @Summary Register a new user
// @Description Register by providing email and password
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "User info"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string
// @Router /users/register [post]
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

// Login godoc
// @Summary User login
// @Description Login with email and password to get JWT token
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "User login credentials"
// @Success 200 {object} map[string]string "JWT token"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 401 {object} map[string]string "Invalid email or password"
// @Router /users/login [post]
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

// CurrentUser godoc
// @Summary Get current user info
// @Description Get info of logged-in user by JWT token
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "User information"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /users/me [get]
// @Security ApiKeyAuth
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
