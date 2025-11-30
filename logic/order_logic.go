package logic

import (
	"database/sql"
	"fmt"

	"shop/dao"
	"shop/global/db"
	"shop/model"
)

// CreateOrder 创建订单
func CreateOrder(userID int, req *model.CreateOrderRequest) (int64, float64, error) {
	// 开始事务
	tx, err := db.DB.Begin()
	if err != nil {
		return 0, 0, fmt.Errorf("创建订单失败: %w", err)
	}
	defer tx.Rollback()

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
			return 0, 0, fmt.Errorf("查询购物车失败: %w", err)
		}
		if cartItem == nil {
			return 0, 0, fmt.Errorf("购物车项不存在")
		}

		// 检查库存
		if cartItem.Quantity > product.Stock {
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
	result, err := tx.Exec(
		"INSERT INTO orders (user_id, total_price, status) VALUES (?, ?, 'pending')",
		userID, totalPrice,
	)
	if err != nil {
		return 0, 0, fmt.Errorf("创建订单失败: %w", err)
	}

	orderID, _ := result.LastInsertId()

	// 创建订单项并更新库存
	for _, item := range orderItems {
		_, err = tx.Exec(
			"INSERT INTO order_items (order_id, product_id, quantity, price) VALUES (?, ?, ?, ?)",
			orderID, item.productID, item.quantity, item.price,
		)
		if err != nil {
			return 0, 0, fmt.Errorf("创建订单项失败: %w", err)
		}

		// 更新库存
		_, err = tx.Exec(
			"UPDATE products SET stock = stock - ? WHERE id = ?",
			item.quantity, item.productID,
		)
		if err != nil {
			return 0, 0, fmt.Errorf("更新库存失败: %w", err)
		}

		// 删除购物车项
		_, err = tx.Exec("DELETE FROM cart_items WHERE id = ? AND user_id = ?", item.cartItemID, userID)
		if err != nil {
			// 这里不阻止订单创建，只记录错误
		}
	}

	// 提交事务
	if err = tx.Commit(); err != nil {
		return 0, 0, fmt.Errorf("提交订单失败: %w", err)
	}

	return orderID, totalPrice, nil
}

// getCartItemWithProductTx 在事务中获取购物车项及其商品信息
func getCartItemWithProductTx(tx *sql.Tx, cartItemID, userID int) (*model.CartItem, *model.Product, error) {
	var cartItem model.CartItem
	var product model.Product
	err := tx.QueryRow(
		`SELECT ci.id, ci.user_id, ci.product_id, ci.quantity,
			p.id, p.name, p.price, p.stock
		FROM cart_items ci
		JOIN products p ON ci.product_id = p.id
		WHERE ci.id = ? AND ci.user_id = ?`,
		cartItemID, userID,
	).Scan(
		&cartItem.ID, &cartItem.UserID, &cartItem.ProductID, &cartItem.Quantity,
		&product.ID, &product.Name, &product.Price, &product.Stock,
	)

	if err == sql.ErrNoRows {
		return nil, nil, nil
	}
	if err != nil {
		return nil, nil, err
	}
	return &cartItem, &product, nil
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
