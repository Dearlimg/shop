package dao

import (
	"database/sql"
	"shop/global/db"
	"shop/model"
)

// CreateOrder 创建订单
func CreateOrder(userID int, totalPrice float64) (int64, error) {
	result, err := db.DB.Exec(
		"INSERT INTO orders (user_id, total_price, status) VALUES (?, ?, 'pending')",
		userID, totalPrice,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// CreateOrderItem 创建订单项
func CreateOrderItem(orderID, productID, quantity int, price float64) error {
	_, err := db.DB.Exec(
		"INSERT INTO order_items (order_id, product_id, quantity, price) VALUES (?, ?, ?, ?)",
		orderID, productID, quantity, price,
	)
	return err
}

// GetOrdersByUserID 获取用户的订单列表
func GetOrdersByUserID(userID int) ([]model.Order, error) {
	rows, err := db.DB.Query(
		"SELECT id, user_id, total_price, status, created_at, updated_at FROM orders WHERE user_id = ? ORDER BY created_at DESC",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var order model.Order
		err := rows.Scan(&order.ID, &order.UserID, &order.TotalPrice, &order.Status, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

// GetOrderByID 根据ID获取订单
func GetOrderByID(orderID, userID int) (*model.Order, error) {
	var order model.Order
	err := db.DB.QueryRow(
		"SELECT id, user_id, total_price, status, created_at, updated_at FROM orders WHERE id = ? AND user_id = ?",
		orderID, userID,
	).Scan(&order.ID, &order.UserID, &order.TotalPrice, &order.Status, &order.CreatedAt, &order.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetOrderItems 获取订单项
func GetOrderItems(orderID int) ([]model.OrderItem, error) {
	rows, err := db.DB.Query(
		`SELECT oi.id, oi.order_id, oi.product_id, oi.quantity, oi.price, oi.created_at, oi.updated_at,
			p.id, p.name, p.description, p.price, p.image, p.stock, p.series, p.created_at, p.updated_at
		FROM order_items oi
		JOIN products p ON oi.product_id = p.id
		WHERE oi.order_id = ?`,
		orderID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.OrderItem
	for rows.Next() {
		var item model.OrderItem
		var product model.Product
		err := rows.Scan(
			&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.Price, &item.CreatedAt, &item.UpdatedAt,
			&product.ID, &product.Name, &product.Description, &product.Price, &product.Image, &product.Stock, &product.Series, &product.CreatedAt, &product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		item.Product = product
		items = append(items, item)
	}
	return items, nil
}

// GetCartItemWithProduct 获取购物车项及其商品信息
func GetCartItemWithProduct(cartItemID, userID int) (*model.CartItem, *model.Product, error) {
	var cartItem model.CartItem
	var product model.Product
	err := db.DB.QueryRow(
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
