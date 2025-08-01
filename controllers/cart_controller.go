package controllers

import (
	"go_bookstore_backend/config"
	"go_bookstore_backend/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AddToCart(c *fiber.Ctx) error {
	//Step 1: Extract user from JWT claims
	user, ok := c.Locals("user").(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	userID, ok := user["id"].(float64) // jwt uses float64 for numbers
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid token data",
		})
	}

	//Step 2: Parse request body
	var req AddToCartRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.BookID == 0 || req.Quantity <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "book_id and quantity must be valid",
		})
	}

	//Step 3: Create cart record
	cart := models.Cart{
		UserID:   uint(userID),
		BookID:   req.BookID,
		Quantity: req.Quantity,
	}

	if err := config.DB.Create(&cart).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add to cart",
		})
	}

	//Step 4: Preload book & user (optional)
	var created models.Cart
	if err := config.DB.Preload("Book").Preload("User").First(&created, cart.ID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to load cart details",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(created)
}

// AddToCartRequest is the expected request body
type AddToCartRequest struct {
	BookID   uint `json:"book_id"`
	Quantity int  `json:"quantity"`
}

func ViewCart(c *fiber.Ctx) error {
	claims, ok := c.Locals("user").(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	userID := uint(claims["id"].(float64))

	var cart []models.Cart
	if err := config.DB.Preload("Book").Where("user_id = ?", userID).Find(&cart).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve cart",
		})
	}

	return c.JSON(cart)
}

func RemoveFromCart(c *fiber.Ctx) error {
	id := c.Params("id")
	claims := c.Locals("user").(jwt.MapClaims)
	userID := uint(claims["id"].(float64))

	var cart models.Cart
	if err := config.DB.First(&cart, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Cart item not found",
		})
	}

	if cart.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You are not allowed to delete this cart item",
		})
	}

	if err := config.DB.Delete(&cart).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to remove cart item",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
