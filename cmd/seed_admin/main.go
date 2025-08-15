// cmd/seed_admin/main.go
//
// ส่วนนี้: สคริปต์ seed ผู้ใช้เริ่มต้นจากค่าใน .env
// - อ่าน ADMIN_EMAIL, ADMIN_PASSWORD, ADMIN_NAME (optional), ADMIN_ROLE (default=admin)
// - ถ้ามีอยู่แล้ว: ข้าม หรือจะให้เขียนทับได้ด้วย SEED_OVERWRITE=true
// - รัน: `go run cmd/seed_admin/main.go`
package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"go_bookstore_backend/config"
	"go_bookstore_backend/models"
	"go_bookstore_backend/utils"
)

func getenv(key, def string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	return v
}

func main() {
	// --- โหลด .env + ต่อฐานข้อมูล ---
	config.LoadEnv()
	config.ConnectDB() // คาดว่าเซ็ต config.DB ให้พร้อมใช้งาน

	// --- อ่านค่าจาก .env (มีค่าเริ่มต้นกันพลาด) ---
	email := getenv("ADMIN_EMAIL", "admin@gmail.com")
	password := getenv("ADMIN_PASSWORD", "admin1234")
	// name := getenv("ADMIN_NAME", "Administrator")
	role := getenv("ADMIN_ROLE", "admin")
	overwrite := strings.EqualFold(getenv("SEED_OVERWRITE", "false"), "true")

	// --- ตรวจว่ามีผู้ใช้อีเมลนี้อยู่หรือยัง ---
	var user models.User
	tx := config.DB.Where("email = ?", email).First(&user)
	found := tx.Error == nil

	switch {
	case found && !overwrite:
		log.Printf("Admin already exists: %s (overwrite disabled)\n", email)
		return

	case found && overwrite:
		// --- เขียนทับรหัสผ่าน/ชื่อ/บทบาท ---
		hashed, err := utils.HashPassword(password)
		if err != nil {
			log.Fatalf("hash password error: %v", err)
		}
		updates := map[string]any{
			"Password": hashed,
			// "Name":     name,
			"Role": role,
		}
		if err := config.DB.Model(&user).Updates(updates).Error; err != nil {
			log.Fatalf("update admin error: %v", err)
		}
		fmt.Printf("Admin updated: %s (role=%s)\n", email, role)
		return

	default:
		// --- สร้างผู้ใช้ใหม่ ---
		hashed, err := utils.HashPassword(password)
		if err != nil {
			log.Fatalf("hash password error: %v", err)
		}
		admin := models.User{
			// Name:     name,
			Email:    email,
			Password: hashed,
			Role:     role,
		}
		if err := config.DB.Create(&admin).Error; err != nil {
			log.Fatalf("create admin error: %v", err)
		}
		fmt.Printf("Admin created: %s / %s (role=%s)\n", email, password, role)
	}
}
