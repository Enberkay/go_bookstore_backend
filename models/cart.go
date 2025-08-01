package models

import "gorm.io/gorm"

type Cart struct {
	ID       uint `gorm:"primaryKey" json:"id"`
	UserID   uint `json:"user_id"`
	BookID   uint `json:"book_id"`
	Quantity int  `json:"Quantity"`

	User User `gorm:"foreignKey:UserID" json:"user"`
	Book Book `gorm:"foreignKey:BookID" json:"book"`
}

func MigrateCarts(db *gorm.DB) error {
	return db.AutoMigrate(&Cart{})
}
