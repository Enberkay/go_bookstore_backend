package models

import "gorm.io/gorm"

type Wishlist struct {
	ID     uint `gorm:"primaryKey" json:"id"`
	UserID uint `json:"user_id"`
	BookID uint `json:"book_id"`

	User User `gorm:"foreignKey:UserID" json:"user"`
	Book Book `gorm:"foreignKey:BookID" json:"book"`
}

func MigrateWishlists(db *gorm.DB) error {
	return db.AutoMigrate(&Wishlist{})
}
