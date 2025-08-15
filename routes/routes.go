// routes/routes.go
//
// ส่วนนี้: ประกาศเส้นทาง API และกำหนดสิทธิ์ด้วย RequireRole แบบ "ชัดเจนตามบทบาท"
// โครงสร้างคอมเมนต์แบ่งเป็นกลุ่ม ๆ (Auth / Books / Cart / Orders / Wishlist)
// - Auth: /me เปิดให้เฉพาะ user และ admin
// - Books: อ่านได้ทั่วไป แต่สร้าง/แก้ไข/ลบ เฉพาะ admin
// - Cart/Orders/Wishlist: เฉพาะ role=user
package routes

import (
	"go_bookstore_backend/controllers"
	"go_bookstore_backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	// --- Auth Routes ---
	auth := api.Group("/auth")
	auth.Post("/register", controllers.Register)
	auth.Post("/login", controllers.Login)
	auth.Get("/me", middlewares.RequireRole("user", "admin"), controllers.CurrentUser)

	// --- Book Routes (admin protected for modification) ---
	books := api.Group("/books")
	books.Get("/", controllers.GetBooks)
	books.Get("/:id", controllers.GetBook)
	books.Post("/", middlewares.RequireRole("admin"), controllers.CreateBook)
	books.Put("/:id", middlewares.RequireRole("admin"), controllers.UpdateBook)
	books.Delete("/:id", middlewares.RequireRole("admin"), controllers.DeleteBook)

	// --- Cart Routes (user authenticated) ---
	cart := api.Group("/cart", middlewares.RequireRole("user"))
	cart.Get("/", controllers.ViewCart)
	cart.Post("/", controllers.AddToCart)
	cart.Delete("/:id", controllers.RemoveFromCart)

	// --- Order Routes (user authenticated) ---
	orders := api.Group("/orders", middlewares.RequireRole("user"))
	orders.Post("/", controllers.PlaceOrder)
	orders.Get("/", controllers.GetMyOrders)

	// --- Wishlist Routes (user authenticated) ---
	wishlist := api.Group("/wishlist", middlewares.RequireRole("user"))
	wishlist.Post("/toggle", controllers.ToggleWishlist)
	wishlist.Get("/", controllers.ViewWishlist)
}
