package controllers

import (
	"go_bookstore_backend/config"
	"go_bookstore_backend/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// AddToCart godoc
// @Summary Add a book to user's cart
// @Description Add a book with quantity to the cart of the logged-in user
// @Tags carts
// @Accept json
// @Produce json
// @Param cart body AddToCartRequest true "Add to cart request"
// @Success 201 {object} models.Cart
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /cart [post]
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
// @Param cart body AddToCartRequest true "Add to cart request"
type AddToCartRequest struct {
	BookID   uint `json:"book_id"`
	Quantity int  `json:"quantity"`
}

// ViewCart godoc
// @Summary View all cart items for logged-in user
// @Description Retrieve all cart items including book and user info
// @Tags carts
// @Accept json
// @Produce json
// @Success 200 {array} models.Cart
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /cart [get]
func ViewCart(c *fiber.Ctx) error {
	// ดึง claims จาก JWT
	claims, ok := c.Locals("user").(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	userIDFloat, ok := claims["id"].(float64) // jwt.MapClaims แปลงตัวเลขเป็น float64
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid token data",
		})
	}

	userID := uint(userIDFloat)

	// Query carts พร้อม preload User และ Book
	var carts []models.Cart
	if err := config.DB.Preload("User").Preload("Book").Where("user_id = ?", userID).Find(&carts).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve cart",
		})
	}

	return c.JSON(carts)
}

// RemoveFromCart godoc
// @Summary Remove an item from cart by ID
// @Description Delete a cart item by cart ID if it belongs to the logged-in user
// @Tags carts
// @Accept json
// @Produce json
// @Param id path int true "Cart item ID"
// @Success 204 "No Content"
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /cart/{id} [delete]
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
