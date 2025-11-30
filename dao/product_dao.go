package dao

import (
	"errors"
	"gorm.io/gorm"
	"shop/global/db"
	"shop/model"
)

// GetProducts 获取所有商品
func GetProducts() ([]model.Product, error) {
	var products []model.Product
	err := db.DB.Order("id DESC").Find(&products).Error
	return products, err
}

// GetProductByID 根据ID获取商品
func GetProductByID(id string) (*model.Product, error) {
	var product model.Product
	err := db.DB.First(&product, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

// GetProductStock 获取商品库存
func GetProductStock(productID int) (int, error) {
	var product model.Product
	err := db.DB.Select("stock").First(&product, productID).Error
	if err != nil {
		return 0, err
	}
	return product.Stock, nil
}

// UpdateProductStock 更新商品库存
func UpdateProductStock(productID, quantity int) error {
	return db.DB.Model(&model.Product{}).
		Where("id = ?", productID).
		Update("stock", gorm.Expr("stock - ?", quantity)).Error
}
