package models

import "gorm.io/gorm"

type Cart struct {
	ID       uint `gorm:"primaryKey" json:"id"`
	UserID   uint `json:"user_id"`
	BookID   uint `json:"book_id"`
	Quantity int  `json:"quantity"`

	User User `gorm:"foreignKey:UserID;references:ID" json:"user"`
	Book Book `gorm:"foreignKey:BookID;references:ID" json:"book"`
}

func MigrateCarts(db *gorm.DB) error {
	return db.AutoMigrate(&Cart{})
}

// DTO สำหรับ response (ไม่มี password)
type CartResponse struct {
	ID       uint         `json:"id"`
	UserID   uint         `json:"user_id"`
	BookID   uint         `json:"book_id"`
	Quantity int          `json:"quantity"`
	User     UserResponse `json:"user"`
	Book     Book         `json:"book"`
}

func (c *Cart) ToResponse() CartResponse {
	return CartResponse{
		ID:       c.ID,
		UserID:   c.UserID,
		BookID:   c.BookID,
		Quantity: c.Quantity,
		User:     c.User.ToResponse(),
		Book:     c.Book,
	}
}
