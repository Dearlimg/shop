package logic

import (
	"fmt"

	"shop/dao"
)

// IncrementCartItem 增量更新购物车商品数量（支持 +1 或 -1，也支持批量增量）
func IncrementCartItem(userID, productID, delta int) error {
	// 获取当前购物车项
	item, err := dao.GetCartItemByUserAndProductFromRedis(userID, productID)
	if err != nil {
		return fmt.Errorf("查询购物车项失败: %w", err)
	}
	if item == nil {
		// 如果不存在且是增加操作，则添加新项
		if delta > 0 {
			// 检查库存
			stock, err := dao.GetProductStock(productID)
			if err != nil {
				return fmt.Errorf("查询商品失败: %w", err)
			}
			if stock < delta {
				return fmt.Errorf("库存不足，当前库存: %d", stock)
			}
			return dao.AddCartItemToRedis(userID, productID, delta)
		}
		return fmt.Errorf("购物车项不存在")
	}

	// 计算新数量
	newQuantity := item.Quantity + delta

	// 如果数量为0或负数，删除该项
	if newQuantity <= 0 {
		return dao.DeleteCartItemFromRedis(userID, productID)
	}

	// 检查库存
	stock, err := dao.GetProductStock(productID)
	if err != nil {
		return fmt.Errorf("查询商品失败: %w", err)
	}

	if newQuantity > stock {
		return fmt.Errorf("库存不足，当前库存: %d", stock)
	}

	// 更新数量
	return dao.UpdateCartItemQuantityInRedis(userID, productID, newQuantity)
}
