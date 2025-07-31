package controllers

import (
	"go_bookstore_backend/config"
	"go_bookstore_backend/models"

	"github.com/gofiber/fiber/v2"
)

func AddToCart(c *fiber.Ctx) error {
	user := c.Locals("user").(map[string]interface{})
	var input struct {
		BookID uint `json:"book_id"`
		Qty    int  `json:"qty"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	item := models.Cart{
		UserID: uint(user["id"].(float64)),
		BookID: input.BookID,
		Qty:    input.Qty,
	}
	config.DB.Create(&item)
	return c.JSON(item)
}

func ViewCart(c *fiber.Ctx) error {
	user := c.Locals("user").(map[string]interface{})
	var cart []models.Cart
	config.DB.Preload("Book").Where("user_id = ?", uint(user["id"].(float64))).Find(&cart)
	return c.JSON(cart)
}

func RemoveFromCart(c *fiber.Ctx) error {
	id := c.Params("id")
	config.DB.Delete(&models.Cart{}, id)
	return c.SendStatus(204)
}
