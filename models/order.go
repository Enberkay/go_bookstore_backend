package models

import "gorm.io/gorm"

type Order struct {
	ID         uint        `gorm:"primaryKey" json:"id"`
	UserID     uint        `json:"user_id"`
	TotalPrice float64     `json:"total_price"`
	User       User        `gorm:"foreignKey:UserID" json:"user"`
	Items      []OrderItem `gorm:"foreignKey:OrderID" json:"items"`
}

type OrderItem struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	OrderID   uint    `json:"order_id"`
	BookID    uint    `json:"book_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
	Book      Book    `gorm:"foreignKey:BookID" json:"book"`
}

func MigrateOrders(db *gorm.DB) error {
	return db.AutoMigrate(&Order{}, &OrderItem{})
}
