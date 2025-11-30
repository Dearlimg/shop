package dao

import (
	"errors"
	"gorm.io/gorm"
	"shop/global/db"
	"shop/model"
)

// CreateOrder 创建订单
func CreateOrder(userID int, totalPrice float64) (*model.Order, error) {
	order := model.Order{
		UserID:     userID,
		TotalPrice: totalPrice,
		Status:     "pending",
	}
	err := db.DB.Create(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// CreateOrderItem 创建订单项
func CreateOrderItem(orderID, productID, quantity int, price float64) error {
	orderItem := model.OrderItem{
		OrderID:   orderID,
		ProductID: productID,
		Quantity:  quantity,
		Price:     price,
	}
	return db.DB.Create(&orderItem).Error
}

// GetOrdersByUserID 获取用户的订单列表
func GetOrdersByUserID(userID int) ([]model.Order, error) {
	var orders []model.Order
	err := db.DB.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&orders).Error
	return orders, err
}

// GetOrderByID 根据ID获取订单
func GetOrderByID(orderID, userID int) (*model.Order, error) {
	var order model.Order
	err := db.DB.Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

// GetOrderItems 获取订单项
func GetOrderItems(orderID int) ([]model.OrderItem, error) {
	var items []model.OrderItem
	err := db.DB.Preload("Product").Where("order_id = ?", orderID).Find(&items).Error
	return items, err
}

// GetCartItemWithProduct 获取购物车项及其商品信息
func GetCartItemWithProduct(cartItemID, userID int) (*model.CartItem, *model.Product, error) {
	var cartItem model.CartItem
	err := db.DB.Preload("Product").Where("id = ? AND user_id = ?", cartItemID, userID).First(&cartItem).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil
		}
		return nil, nil, err
	}
	return &cartItem, &cartItem.Product, nil
}
