package dao

import (
	"errors"
	"gorm.io/gorm"
	"shop/global/db"
	"shop/model"
)

// GetCartItems 获取用户的购物车项
func GetCartItems(userID int) ([]model.CartItem, error) {
	var items []model.CartItem
	err := db.DB.Preload("Product").Where("user_id = ?", userID).Find(&items).Error
	return items, err
}

// GetCartItemByID 根据ID获取购物车项
func GetCartItemByID(cartItemID, userID int) (*model.CartItem, error) {
	var item model.CartItem
	err := db.DB.Where("id = ? AND user_id = ?", cartItemID, userID).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

// GetCartItemByUserAndProduct 获取用户和商品的购物车项
func GetCartItemByUserAndProduct(userID, productID int) (*model.CartItem, error) {
	var item model.CartItem
	err := db.DB.Where("user_id = ? AND product_id = ?", userID, productID).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

// AddCartItem 添加购物车项
func AddCartItem(userID, productID, quantity int) error {
	cartItem := model.CartItem{
		UserID:    userID,
		ProductID: productID,
		Quantity:  quantity,
	}
	return db.DB.Create(&cartItem).Error
}

// UpdateCartItemQuantity 更新购物车项数量
func UpdateCartItemQuantity(cartItemID, quantity int) error {
	return db.DB.Model(&model.CartItem{}).
		Where("id = ?", cartItemID).
		Update("quantity", quantity).Error
}

// DeleteCartItem 删除购物车项
func DeleteCartItem(cartItemID int) error {
	return db.DB.Delete(&model.CartItem{}, cartItemID).Error
}
