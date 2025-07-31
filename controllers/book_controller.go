package controllers

import (
	"go_bookstore_backend/config"
	"go_bookstore_backend/models"

	"github.com/gofiber/fiber/v2"
)

func CreateBook(c *fiber.Ctx) error {
	var book models.Book
	if err := c.BodyParser(&book); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	config.DB.Create(&book)
	return c.JSON(book)
}

func GetBooks(c *fiber.Ctx) error {
	var books []models.Book
	config.DB.Find(&books)
	return c.JSON(books)
}

func GetBook(c *fiber.Ctx) error {
	id := c.Params("id")
	var book models.Book
	if err := config.DB.First(&book, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Book not found"})
	}
	return c.JSON(book)
}

func UpdateBook(c *fiber.Ctx) error {
	id := c.Params("id")
	var book models.Book
	if err := config.DB.First(&book, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Not found"})
	}
	if err := c.BodyParser(&book); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	config.DB.Save(&book)
	return c.JSON(book)
}

func DeleteBook(c *fiber.Ctx) error {
	id := c.Params("id")
	config.DB.Delete(&models.Book{}, id)
	return c.SendStatus(204)
}
