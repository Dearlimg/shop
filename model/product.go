package model

import "time"

// Product 商品模型
type Product struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	Price       float64   `json:"price" gorm:"not null"`
	Image       string    `json:"image"`
	Stock       int       `json:"stock" gorm:"default:0"`
	Series      string    `json:"series" gorm:"default:'拉布布'"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
