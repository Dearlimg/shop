package model

import "time"

// Order 订单模型
type Order struct {
	ID         int         `json:"id" gorm:"primaryKey;autoIncrement;type:int"`
	UserID     int         `json:"user_id" gorm:"type:int;not null;index:idx_user_id"`
	TotalPrice float64     `json:"total_price" gorm:"type:decimal(10,2);not null"`
	Status     string      `json:"status" gorm:"type:varchar(20);default:'pending'"`
	Items      []OrderItem `json:"items" gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
	CreatedAt  time.Time   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time   `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (Order) TableName() string {
	return "orders"
}

// OrderItem 订单项模型
type OrderItem struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement;type:int"`
	OrderID   int       `json:"order_id" gorm:"type:int;not null;index:idx_order_id"`
	ProductID int       `json:"product_id" gorm:"type:int;not null;index:idx_product_id"`
	Quantity  int       `json:"quantity" gorm:"type:int;not null"`
	Price     float64   `json:"price" gorm:"type:decimal(10,2);not null"`
	Product   Product   `json:"product" gorm:"foreignKey:ProductID"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (OrderItem) TableName() string {
	return "order_items"
}

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	CartItemIDs []int `json:"cart_item_ids" binding:"required"`
}
