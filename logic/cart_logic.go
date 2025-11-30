package logic

import (
	"fmt"

	"shop/dao"
	"shop/model"
)

// GetCart 获取购物车（使用Redis）
func GetCart(userID int) ([]model.CartItem, error) {
	return dao.GetCartItemsFromRedis(userID)
}

// AddToCart 添加到购物车（使用Redis）
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
	existingItem, err := dao.GetCartItemByUserAndProductFromRedis(userID, req.ProductID)
	if err != nil {
		return fmt.Errorf("查询购物车失败: %w", err)
	}

	if existingItem != nil {
		// 更新数量
		newQuantity := existingItem.Quantity + req.Quantity
		if newQuantity > stock {
			return fmt.Errorf("库存不足")
		}
		return dao.UpdateCartItemQuantityInRedis(userID, req.ProductID, newQuantity)
	}

	// 添加新商品到购物车
	return dao.AddCartItemToRedis(userID, req.ProductID, req.Quantity)
}

// UpdateCartItem 更新购物车商品数量（使用Redis）
// 注意：参数改为 productID 而不是 cartItemID
func UpdateCartItem(userID, productID int, req *model.UpdateCartRequest) error {
	// 检查购物车项是否存在
	item, err := dao.GetCartItemByUserAndProductFromRedis(userID, productID)
	if err != nil {
		return fmt.Errorf("查询购物车项失败: %w", err)
	}
	if item == nil {
		return fmt.Errorf("购物车项不存在")
	}

	// 检查库存
	stock, err := dao.GetProductStock(productID)
	if err != nil {
		return fmt.Errorf("查询商品失败: %w", err)
	}

	if req.Quantity > stock {
		return fmt.Errorf("库存不足")
	}

	// 更新数量
	return dao.UpdateCartItemQuantityInRedis(userID, productID, req.Quantity)
}

// DeleteCartItem 删除购物车商品（使用Redis）
// 注意：参数改为 productID 而不是 cartItemID
func DeleteCartItem(userID, productID int) error {
	// 检查购物车项是否存在
	item, err := dao.GetCartItemByUserAndProductFromRedis(userID, productID)
	if err != nil {
		return fmt.Errorf("查询购物车项失败: %w", err)
	}
	if item == nil {
		return fmt.Errorf("购物车项不存在")
	}

	return dao.DeleteCartItemFromRedis(userID, productID)
}
