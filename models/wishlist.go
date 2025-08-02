package models

import "gorm.io/gorm"

type Wishlist struct {
	ID     uint `gorm:"primaryKey"`
	UserID uint `gorm:"uniqueIndex:idx_user_book"`
	BookID uint `gorm:"uniqueIndex:idx_user_book"`

	User User `gorm:"foreignKey:UserID" json:"user"`
	Book Book `gorm:"foreignKey:BookID" json:"book"`
}

func MigrateWishlists(db *gorm.DB) error {
	return db.AutoMigrate(&Wishlist{})
}
