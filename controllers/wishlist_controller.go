package controllers

import (
	"go_bookstore_backend/config"
	"go_bookstore_backend/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AddToWishlist(c *fiber.Ctx) error {
	claims := c.Locals("user").(jwt.MapClaims)
	userID := uint(claims["id"].(float64))

	var body struct {
		BookID uint `json:"book_id"`
	}
	if err := c.BodyParser(&body); err != nil || body.BookID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid book ID"})
	}

	wish := models.Wishlist{UserID: userID, BookID: body.BookID}
	if err := config.DB.Create(&wish).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add to wishlist"})
	}

	return c.Status(fiber.StatusCreated).JSON(wish)
}

func ViewWishlist(c *fiber.Ctx) error {
	claims := c.Locals("user").(jwt.MapClaims)
	userID := uint(claims["id"].(float64))

	var wishlist []models.Wishlist
	if err := config.DB.Preload("Book").Where("user_id = ?", userID).Find(&wishlist).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to load wishlist"})
	}

	return c.JSON(wishlist)
}

func RemoveFromWishlist(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := config.DB.Delete(&models.Wishlist{}, id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to remove item"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
