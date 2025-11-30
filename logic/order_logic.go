package logic

import (
	"fmt"

	"gorm.io/gorm"
	"shop/dao"
	"shop/global/db"
	"shop/model"
)

// CreateOrder 创建订单（使用购物车中所有商品）
func CreateOrder(userID int, req *model.CreateOrderRequest) (int64, float64, error) {
	// 开始事务
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 获取购物车中所有商品
	cartItems, err := dao.GetCartItemsFromRedis(userID)
	if err != nil {
		tx.Rollback()
		return 0, 0, fmt.Errorf("查询购物车失败: %w", err)
	}

	if len(cartItems) == 0 {
		tx.Rollback()
		return 0, 0, fmt.Errorf("购物车为空")
	}

	// 如果指定了商品ID列表，则只处理指定的商品；否则处理所有商品
	var itemsToProcess []model.CartItem
	if req != nil && req.CartItemIDs != nil && len(req.CartItemIDs) > 0 {
		// 只处理指定的商品
		productIDMap := make(map[int]bool)
		for _, id := range req.CartItemIDs {
			productIDMap[id] = true
		}
		for _, item := range cartItems {
			if productIDMap[item.ProductID] {
				itemsToProcess = append(itemsToProcess, item)
			}
		}
		if len(itemsToProcess) == 0 {
			tx.Rollback()
			return 0, 0, fmt.Errorf("指定的商品不在购物车中")
		}
	} else {
		// 处理购物车中所有商品
		itemsToProcess = cartItems
	}

	// 查询购物车项并计算总价
	var totalPrice float64
	var orderItems []struct {
		cartItemID int
		productID  int
		quantity   int
		price      float64
	}

	for _, cartItem := range itemsToProcess {
		// 获取商品信息
		product, err := dao.GetProductByID(fmt.Sprintf("%d", cartItem.ProductID))
		if err != nil || product == nil {
			tx.Rollback()
			return 0, 0, fmt.Errorf("商品不存在: %d", cartItem.ProductID)
		}

		// 检查库存
		if cartItem.Quantity > product.Stock {
			tx.Rollback()
			return 0, 0, fmt.Errorf("商品库存不足: %s (需要: %d, 库存: %d)", product.Name, cartItem.Quantity, product.Stock)
		}

		itemTotal := product.Price * float64(cartItem.Quantity)
		totalPrice += itemTotal
		orderItems = append(orderItems, struct {
			cartItemID int
			productID  int
			quantity   int
			price      float64
		}{cartItem.ProductID, cartItem.ProductID, cartItem.Quantity, product.Price})
	}

	// 创建订单
	order := model.Order{
		UserID:     userID,
		TotalPrice: totalPrice,
		Status:     "pending",
	}
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return 0, 0, fmt.Errorf("创建订单失败: %w", err)
	}

	// 创建订单项并更新库存
	for _, item := range orderItems {
		// 创建订单项
		orderItem := model.OrderItem{
			OrderID:   order.ID,
			ProductID: item.productID,
			Quantity:  item.quantity,
			Price:     item.price,
		}
		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			return 0, 0, fmt.Errorf("创建订单项失败: %w", err)
		}

		// 更新库存
		if err := tx.Model(&model.Product{}).
			Where("id = ?", item.productID).
			Update("stock", gorm.Expr("stock - ?", item.quantity)).Error; err != nil {
			tx.Rollback()
			return 0, 0, fmt.Errorf("更新库存失败: %w", err)
		}

		// 从Redis删除购物车项
		if err := dao.DeleteCartItemFromRedis(userID, item.productID); err != nil {
			// 这里不阻止订单创建，只记录错误
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return 0, 0, fmt.Errorf("提交订单失败: %w", err)
	}

	return int64(order.ID), totalPrice, nil
}

// GetOrders 获取订单历史
func GetOrders(userID int) ([]model.Order, error) {
	orders, err := dao.GetOrdersByUserID(userID)
	if err != nil {
		return nil, err
	}

	// 加载订单项
	for i := range orders {
		items, err := dao.GetOrderItems(orders[i].ID)
		if err == nil {
			orders[i].Items = items
		}
	}

	return orders, nil
}

// GetOrder 获取订单详情
func GetOrder(userID, orderID int) (*model.Order, error) {
	order, err := dao.GetOrderByID(orderID, userID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, fmt.Errorf("订单不存在")
	}

	// 加载订单项
	items, err := dao.GetOrderItems(order.ID)
	if err == nil {
		order.Items = items
	}

	return order, nil
}
