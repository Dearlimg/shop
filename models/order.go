package models

import "time"

// Order 订单模型
type Order struct {
	ID         int         `json:"id" gorm:"primaryKey"`
	UserID     int         `json:"user_id" gorm:"not null;index"`
	TotalPrice float64     `json:"total_price" gorm:"not null"`
	Status     string      `json:"status" gorm:"default:'pending'"`
	Items      []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

// OrderItem 订单项模型
type OrderItem struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	OrderID   int       `json:"order_id" gorm:"not null;index"`
	ProductID int       `json:"product_id" gorm:"not null"`
	Quantity  int       `json:"quantity" gorm:"not null"`
	Price     float64   `json:"price" gorm:"not null"`
	Product   Product   `json:"product" gorm:"foreignKey:ProductID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	CartItemIDs []int `json:"cart_item_ids" binding:"required"`
}
