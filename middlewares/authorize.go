// middlewares/authorize.go
//
// ส่วนนี้: Middleware รวมขั้นตอนตรวจ JWT + ตรวจ role ให้อยู่ในฟังก์ชันเดียว
// - ต้องมี Authorization: Bearer <token> ที่ถูกต้องก่อน
// - ถ้าระบุ roles เข้ามา จะต้องมี claims["role"] ตรงกับหนึ่งใน roles นั้น
// - แนบ claims ไว้ที่ c.Locals("user") เผื่อ controller ใช้งานต่อ
package middlewares

import (
	"strings"

	"go_bookstore_backend/utils"

	"github.com/gofiber/fiber/v2"
)

// RequireRole ตรวจทั้ง token และ role (ถ้ามี)
// ตัวอย่าง:
//
//	RequireRole("user")              // เฉพาะผู้ใช้ role=user
//	RequireRole("admin")             // เฉพาะผู้ใช้ role=admin
//	RequireRole("user","admin")      // สองบทบาทนี้เท่านั้น
func RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// --- ตรวจรูปแบบ Authorization header ---
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing or invalid token",
			})
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// --- แกะ JWT -> claims (utils.ParseJWT ต้องคืน map[string]any/claims ที่อ่านได้) ---
		claims, err := utils.ParseJWT(tokenStr)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}
		// แนบ claims ให้ controller เรียกใช้ต่อ (เช่น id, email, role)
		c.Locals("user", claims)

		// --- ตรวจ role (เฉพาะเมื่อมีการกำหนด roles) ---
		if len(roles) > 0 {
			rawRole, ok := claims["role"]
			if !ok {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error": "Access denied",
				})
			}
			roleStr, ok := rawRole.(string)
			if !ok || !in(roleStr, roles) {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error": "Access denied",
				})
			}
		}

		return c.Next()
	}
}

// in: helper เช็คสมาชิกแบบตรงไปตรงมา (ไม่ซับซ้อน)
func in(v string, list []string) bool {
	for _, s := range list {
		if v == s {
			return true
		}
	}
	return false
}
