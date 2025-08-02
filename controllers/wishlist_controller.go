package controllers

import (
	"go_bookstore_backend/config"
	"go_bookstore_backend/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// ToggleWishlistRequest ใช้เป็น request body สำหรับ toggle wishlist
type ToggleWishlistRequest struct {
	BookID uint `json:"book_id" example:"1"`
}

// MessageResponse สำหรับส่งข้อความเช่น "Removed from wishlist"
type MessageResponse struct {
	Message string `json:"message" example:"Removed from wishlist"`
}

// ErrorResponse ใช้สำหรับ error responses
type ErrorResponse struct {
	Error string `json:"error" example:"Unauthorized"`
}

// ToggleWishlist godoc
// @Summary Add or remove book from wishlist (toggle behavior)
// @Description Toggle a book in the user's wishlist. Adds if not exists, removes if exists.
// @Tags Wishlist
// @Accept json
// @Produce json
// @Param wishlist body ToggleWishlistRequest true "Book ID to toggle"
// @Success 201 {object} models.Wishlist
// @Success 200 {object} MessageResponse "Removed from wishlist"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /api/wishlist/toggle [post]
// @Security BearerAuth
func ToggleWishlist(c *fiber.Ctx) error {
	claims := c.Locals("user").(jwt.MapClaims)
	userID := uint(claims["id"].(float64))

	var body ToggleWishlistRequest
	if err := c.BodyParser(&body); err != nil || body.BookID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Invalid book ID"})
	}

	var existing models.Wishlist
	err := config.DB.
		Where("user_id = ? AND book_id = ?", userID, body.BookID).
		First(&existing).Error

	if err == nil {
		if err := config.DB.Delete(&existing).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: "Failed to remove from wishlist"})
		}
		return c.Status(fiber.StatusOK).JSON(MessageResponse{Message: "Removed from wishlist"})
	}

	wish := models.Wishlist{UserID: userID, BookID: body.BookID}
	if err := config.DB.Create(&wish).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: "Failed to add to wishlist"})
	}

	var result models.Wishlist
	if err := config.DB.Preload("Book").Preload("User").First(&result, wish.ID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: "Failed to load wishlist item"})
	}

	return c.Status(fiber.StatusCreated).JSON(result)
}

// ViewWishlist godoc
// @Summary View current user's wishlist
// @Description Retrieve all wishlist items of the authenticated user
// @Tags Wishlist
// @Produce json
// @Success 200 {array} models.Wishlist
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /api/wishlist [get]
// @Security BearerAuth
func ViewWishlist(c *fiber.Ctx) error {
	claims := c.Locals("user").(jwt.MapClaims)
	userID := uint(claims["id"].(float64))

	var wishlist []models.Wishlist
	if err := config.DB.Preload("Book").Preload("User").Where("user_id = ?", userID).Find(&wishlist).Error; err != nil {
		return c.Status(500).JSON(ErrorResponse{Error: "Failed to load wishlist"})
	}

	return c.JSON(wishlist)
}
