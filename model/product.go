package model

import "time"

// Product 商品模型
type Product struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement;type:int"`
	Name        string    `json:"name" gorm:"type:varchar(200);not null"`
	Description string    `json:"description" gorm:"type:text"`
	Price       float64   `json:"price" gorm:"type:decimal(10,2);not null"`
	Image       string    `json:"image" gorm:"type:varchar(500)"`
	Stock       int       `json:"stock" gorm:"type:int;default:0"`
	Series      string    `json:"series" gorm:"type:varchar(50);default:'拉布布'"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (Product) TableName() string {
	return "products"
}
