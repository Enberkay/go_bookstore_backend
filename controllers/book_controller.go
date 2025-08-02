package controllers

import (
	"go_bookstore_backend/config"
	"go_bookstore_backend/models"

	"github.com/gofiber/fiber/v2"
)

// CreateBook godoc
// @Summary Create a new book
// @Description Create a new book with details
// @Tags books
// @Accept json
// @Produce json
// @Param book body models.Book true "Book data"
// @Success 200 {object} models.Book
// @Failure 400 {object} map[string]string
// @Router /books [post]
func CreateBook(c *fiber.Ctx) error {
	var book models.Book
	if err := c.BodyParser(&book); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	config.DB.Create(&book)
	return c.JSON(book)
}

// GetBooks godoc
// @Summary Get list of books
// @Description Get all books available in the store
// @Tags books
// @Accept json
// @Produce json
// @Success 200 {array} models.Book
// @Router /books [get]
func GetBooks(c *fiber.Ctx) error {
	var books []models.Book
	config.DB.Find(&books)
	return c.JSON(books)
}

// GetBook godoc
// @Summary Get book by ID
// @Description Get detailed information about a book by its ID
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} models.Book
// @Failure 404 {object} map[string]string
// @Router /books/{id} [get]
func GetBook(c *fiber.Ctx) error {
	id := c.Params("id")
	var book models.Book
	if err := config.DB.First(&book, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Book not found"})
	}
	return c.JSON(book)
}

// UpdateBook godoc
// @Summary Update a book by ID
// @Description Update book details by its ID
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Param book body models.Book true "Updated book data"
// @Success 200 {object} models.Book
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /books/{id} [put]
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

// DeleteBook godoc
// @Summary Delete a book by ID
// @Description Delete a book from the store by its ID
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Success 204 "No Content"
// @Router /books/{id} [delete]
func DeleteBook(c *fiber.Ctx) error {
	id := c.Params("id")
	config.DB.Delete(&models.Book{}, id)
	return c.SendStatus(204)
}
