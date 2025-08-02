package models

import "gorm.io/gorm"

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Email    string `gorm:"uniqueIndex" json:"email"`
	Password string `json:"password"` // ใช้ตอนรับ input
	Role     string `json:"role"`     // "admin" หรือ "user"
}

func MigrateUsers(db *gorm.DB) error {
	return db.AutoMigrate(&User{})
}

// DTO สำหรับ Response (ไม่ expose password)
type UserResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// Helper function สำหรับแปลงจาก User เป็น UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:    u.ID,
		Email: u.Email,
		Role:  u.Role,
	}
}
