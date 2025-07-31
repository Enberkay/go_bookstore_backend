package models

import "gorm.io/gorm"

type Book struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	Title       string  `json:"title"`
	Author      string  `json:"author"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

func MigrateBooks(db *gorm.DB) error {
	return db.AutoMigrate(&Book{})
}
