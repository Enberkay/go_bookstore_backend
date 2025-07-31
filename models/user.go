package models

import "gorm.io/gorm"

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Email    string `gorm:"uniqueIndex" json:"email"`
	Password string `json:"-"`
	Role     string `json:"role"` // "admin" หรือ "user"
}

func MigrateUsers(db *gorm.DB) error {
	return db.AutoMigrate(&User{})
}
