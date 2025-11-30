package logic

import (
	"shop/dao"
	"shop/model"
)

// GetProducts 获取商品列表
func GetProducts() ([]model.Product, error) {
	return dao.GetProducts()
}

// GetProduct 获取商品详情
func GetProduct(id string) (*model.Product, error) {
	return dao.GetProductByID(id)
}
