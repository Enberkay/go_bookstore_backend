package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID         uint        `gorm:"primaryKey" json:"id"`
	UserID     uint        `json:"user_id"`
	User       User        `gorm:"foreignKey:UserID" json:"user"`
	TotalPrice float64     `json:"total_price"`
	CreatedAt  time.Time   `json:"created_at"`
	Items      []OrderItem `gorm:"foreignKey:OrderID" json:"items"`
}

type OrderItem struct {
	ID       uint    `gorm:"primaryKey" json:"id"`
	OrderID  uint    `json:"order_id"`
	Order    Order   `gorm:"foreignKey:OrderID" json:"-"`
	BookID   uint    `json:"book_id"`
	Book     Book    `gorm:"foreignKey:BookID" json:"book"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"` // Price per unit at time of purchase
}

func MigrateOrders(db *gorm.DB) error {
	return db.AutoMigrate(&Order{}, &OrderItem{})
}
