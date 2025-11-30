package logic

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
	"shop/dao"
	"shop/global/db"
	"shop/model"
)

// CreateOrder 创建订单
func CreateOrder(userID int, req *model.CreateOrderRequest) (int64, float64, error) {
	// 开始事务
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 查询购物车项并计算总价
	var totalPrice float64
	var orderItems []struct {
		cartItemID int
		productID  int
		quantity   int
		price      float64
	}

	for _, cartItemID := range req.CartItemIDs {
		cartItem, product, err := getCartItemWithProductTx(tx, cartItemID, userID)
		if err != nil {
			tx.Rollback()
			return 0, 0, fmt.Errorf("查询购物车失败: %w", err)
		}
		if cartItem == nil {
			tx.Rollback()
			return 0, 0, fmt.Errorf("购物车项不存在")
		}

		// 检查库存
		if cartItem.Quantity > product.Stock {
			tx.Rollback()
			return 0, 0, fmt.Errorf("商品库存不足: %s", product.Name)
		}

		itemTotal := product.Price * float64(cartItem.Quantity)
		totalPrice += itemTotal
		orderItems = append(orderItems, struct {
			cartItemID int
			productID  int
			quantity   int
			price      float64
		}{cartItem.ID, cartItem.ProductID, cartItem.Quantity, product.Price})
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

		// 删除购物车项
		if err := tx.Where("id = ? AND user_id = ?", item.cartItemID, userID).
			Delete(&model.CartItem{}).Error; err != nil {
			// 这里不阻止订单创建，只记录错误
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return 0, 0, fmt.Errorf("提交订单失败: %w", err)
	}

	return int64(order.ID), totalPrice, nil
}

// getCartItemWithProductTx 在事务中获取购物车项及其商品信息
func getCartItemWithProductTx(tx *gorm.DB, cartItemID, userID int) (*model.CartItem, *model.Product, error) {
	var cartItem model.CartItem
	err := tx.Preload("Product").Where("id = ? AND user_id = ?", cartItemID, userID).First(&cartItem).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil
		}
		return nil, nil, err
	}
	return &cartItem, &cartItem.Product, nil
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
