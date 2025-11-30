package model

import "time"

// CartItem 购物车项模型
type CartItem struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement;type:int"`
	UserID    int       `json:"user_id" gorm:"type:int;not null;index:idx_user_id"`
	ProductID int       `json:"product_id" gorm:"type:int;not null;index:idx_product_id"`
	Quantity  int       `json:"quantity" gorm:"type:int;default:1"`
	Product   Product   `json:"product" gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (CartItem) TableName() string {
	return "cart_items"
}

// AddToCartRequest 添加到购物车请求
type AddToCartRequest struct {
	ProductID int `json:"product_id" binding:"required"`
	Quantity  int `json:"quantity" binding:"required,min=1"`
}

// UpdateCartRequest 更新购物车请求
type UpdateCartRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}

// IncrementCartRequest 增量更新购物车请求
type IncrementCartRequest struct {
	Delta int `json:"delta" binding:"required"` // 增量值，+1 或 -1
}
