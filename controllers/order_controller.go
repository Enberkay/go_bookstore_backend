package controllers

import (
	"go_bookstore_backend/config"
	"go_bookstore_backend/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func PlaceOrder(c *fiber.Ctx) error {
	claims, ok := c.Locals("user").(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	userID := uint(claims["id"].(float64))

	// Step 1: ดึงตะกร้าสินค้า
	var cartItems []models.Cart
	if err := config.DB.Preload("Book").Where("user_id = ?", userID).Find(&cartItems).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch cart"})
	}
	if len(cartItems) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Cart is empty"})
	}

	// Step 2: เตรียม Order และ OrderItems
	var orderItems []models.OrderItem
	var total float64 = 0

	for _, item := range cartItems {
		if item.Book.Stock < item.Quantity {
			return c.Status(400).JSON(fiber.Map{
				"error": "Not enough stock for book: " + item.Book.Title,
			})
		}

		total += float64(item.Book.Price) * float64(item.Quantity)

		orderItems = append(orderItems, models.OrderItem{
			BookID:    item.BookID,
			Quantity:  item.Quantity,
			UnitPrice: item.Book.Price,
		})

		// หัก stock
		item.Book.Stock -= item.Quantity
		config.DB.Save(&item.Book)
	}

	order := models.Order{
		UserID:     userID,
		TotalPrice: total,
		Items:      orderItems,
	}

	// Step 3: บันทึก Order และล้าง cart
	if err := config.DB.Create(&order).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to place order"})
	}
	config.DB.Where("user_id = ?", userID).Delete(&models.Cart{})

	// Step 4: โหลดข้อมูลออเดอร์แบบเต็ม
	var fullOrder models.Order
	if err := config.DB.Preload("Items.Book").Preload("User").First(&fullOrder, order.ID).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to load order"})
	}

	return c.Status(201).JSON(fullOrder)
}

func GetMyOrders(c *fiber.Ctx) error {
	claims, ok := c.Locals("user").(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	userID := uint(claims["id"].(float64))

	var orders []models.Order
	if err := config.DB.
		Preload("User").
		Preload("Items").
		Preload("Items.Book").
		Where("user_id = ?", userID).
		Find(&orders).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch orders"})
	}

	return c.JSON(orders)
}
