package logic

import (
	"fmt"

	"shop/dao"
	"shop/model"
)

// GetCart 获取购物车
func GetCart(userID int) ([]model.CartItem, error) {
	return dao.GetCartItems(userID)
}

// AddToCart 添加到购物车
func AddToCart(userID int, req *model.AddToCartRequest) error {
	// 检查商品是否存在
	stock, err := dao.GetProductStock(req.ProductID)
	if err != nil {
		return fmt.Errorf("查询商品失败: %w", err)
	}
	if stock < req.Quantity {
		return fmt.Errorf("库存不足")
	}

	// 检查购物车中是否已有该商品
	existingItem, err := dao.GetCartItemByUserAndProduct(userID, req.ProductID)
	if err != nil {
		return fmt.Errorf("查询购物车失败: %w", err)
	}

	if existingItem != nil {
		// 更新数量
		newQuantity := existingItem.Quantity + req.Quantity
		if newQuantity > stock {
			return fmt.Errorf("库存不足")
		}
		return dao.UpdateCartItemQuantity(existingItem.ID, newQuantity)
	}

	// 添加新商品到购物车
	return dao.AddCartItem(userID, req.ProductID, req.Quantity)
}

// UpdateCartItem 更新购物车商品数量
func UpdateCartItem(userID, cartItemID int, req *model.UpdateCartRequest) error {
	// 检查购物车项是否属于当前用户
	item, err := dao.GetCartItemByID(cartItemID, userID)
	if err != nil {
		return fmt.Errorf("查询购物车项失败: %w", err)
	}
	if item == nil {
		return fmt.Errorf("购物车项不存在")
	}

	// 检查库存
	stock, err := dao.GetProductStock(item.ProductID)
	if err != nil {
		return fmt.Errorf("查询商品失败: %w", err)
	}

	if req.Quantity > stock {
		return fmt.Errorf("库存不足")
	}

	// 更新数量
	return dao.UpdateCartItemQuantity(cartItemID, req.Quantity)
}

// DeleteCartItem 删除购物车商品
func DeleteCartItem(userID, cartItemID int) error {
	// 检查购物车项是否属于当前用户
	item, err := dao.GetCartItemByID(cartItemID, userID)
	if err != nil {
		return fmt.Errorf("查询购物车项失败: %w", err)
	}
	if item == nil {
		return fmt.Errorf("购物车项不存在")
	}

	return dao.DeleteCartItem(cartItemID)
}
