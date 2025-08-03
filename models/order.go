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

// ---------- Response DTO ----------

type OrderItemResponse struct {
	ID        uint    `json:"id"`
	OrderID   uint    `json:"order_id"`
	BookID    uint    `json:"book_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
	Book      Book    `json:"book"`
}

type OrderResponse struct {
	ID         uint                `json:"id"`
	UserID     uint                `json:"user_id"`
	TotalPrice float64             `json:"total_price"`
	User       UserResponse        `json:"user"`
	Items      []OrderItemResponse `json:"items"`
}

func (o *Order) ToResponse() OrderResponse {
	var items []OrderItemResponse
	for _, item := range o.Items {
		items = append(items, OrderItemResponse{
			ID:        item.ID,
			OrderID:   item.OrderID,
			BookID:    item.BookID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
			Book:      item.Book,
		})
	}

	return OrderResponse{
		ID:         o.ID,
		UserID:     o.UserID,
		TotalPrice: o.TotalPrice,
		User:       o.User.ToResponse(), // ไม่มี password
		Items:      items,
	}
}
