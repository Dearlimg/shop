package handlers

import (
	"context"
	"database/sql"
	"strconv"

	"shop/database"
	"shop/models"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

// CreateOrder 创建订单
func CreateOrder(ctx context.Context, c *app.RequestContext) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, utils.H{
			"error": "未授权",
		})
		return
	}

	var req models.CreateOrderRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(400, utils.H{
			"error": "请求参数错误: " + err.Error(),
		})
		return
	}

	// 开始事务
	tx, err := database.DB.Begin()
	if err != nil {
		c.JSON(500, utils.H{
			"error": "创建订单失败",
		})
		return
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
		var cartItem models.CartItem
		var product models.Product
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
			tx.Rollback()
			c.JSON(404, utils.H{
				"error": "购物车项不存在",
			})
			return
		}
		if err != nil {
			tx.Rollback()
			c.JSON(500, utils.H{
				"error": "查询购物车失败",
			})
			return
		}

		// 检查库存
		if cartItem.Quantity > product.Stock {
			tx.Rollback()
			c.JSON(400, utils.H{
				"error": "商品库存不足: " + product.Name,
			})
			return
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
		tx.Rollback()
		c.JSON(500, utils.H{
			"error": "创建订单失败",
		})
		return
	}

	orderID, _ := result.LastInsertId()

	// 创建订单项并更新库存
	for _, item := range orderItems {
		_, err = tx.Exec(
			"INSERT INTO order_items (order_id, product_id, quantity, price) VALUES (?, ?, ?, ?)",
			orderID, item.productID, item.quantity, item.price,
		)
		if err != nil {
			tx.Rollback()
			c.JSON(500, utils.H{
				"error": "创建订单项失败",
			})
			return
		}

		// 更新库存
		_, err = tx.Exec(
			"UPDATE products SET stock = stock - ? WHERE id = ?",
			item.quantity, item.productID,
		)
		if err != nil {
			tx.Rollback()
			c.JSON(500, utils.H{
				"error": "更新库存失败",
			})
			return
		}

		// 删除购物车项
		_, err = tx.Exec("DELETE FROM cart_items WHERE id = ? AND user_id = ?", item.cartItemID, userID)
		if err != nil {
			// 这里不阻止订单创建，只记录错误
		}
	}

	// 提交事务
	if err = tx.Commit(); err != nil {
		c.JSON(500, utils.H{
			"error": "提交订单失败",
		})
		return
	}

	c.JSON(200, utils.H{
		"message":     "订单创建成功",
		"order_id":    orderID,
		"total_price": totalPrice,
	})
}

// GetOrders 获取订单历史
func GetOrders(ctx context.Context, c *app.RequestContext) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, utils.H{
			"error": "未授权",
		})
		return
	}

	rows, err := database.DB.Query(
		"SELECT id, user_id, total_price, status, created_at, updated_at FROM orders WHERE user_id = ? ORDER BY created_at DESC",
		userID,
	)
	if err != nil {
		c.JSON(500, utils.H{
			"error": "查询订单失败: " + err.Error(),
		})
		return
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(&order.ID, &order.UserID, &order.TotalPrice, &order.Status, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			c.JSON(500, utils.H{
				"error": "读取订单数据失败",
			})
			return
		}

		// 查询订单项
		itemRows, err := database.DB.Query(
			`SELECT oi.id, oi.order_id, oi.product_id, oi.quantity, oi.price, oi.created_at, oi.updated_at,
				p.id, p.name, p.description, p.price, p.image, p.stock, p.series, p.created_at, p.updated_at
			FROM order_items oi
			JOIN products p ON oi.product_id = p.id
			WHERE oi.order_id = ?`,
			order.ID,
		)
		if err == nil {
			defer itemRows.Close()
			for itemRows.Next() {
				var item models.OrderItem
				var product models.Product
				err := itemRows.Scan(
					&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.Price, &item.CreatedAt, &item.UpdatedAt,
					&product.ID, &product.Name, &product.Description, &product.Price, &product.Image, &product.Stock, &product.Series, &product.CreatedAt, &product.UpdatedAt,
				)
				if err == nil {
					item.Product = product
					order.Items = append(order.Items, item)
				}
			}
		}

		orders = append(orders, order)
	}

	c.JSON(200, utils.H{
		"orders": orders,
	})
}

// GetOrder 获取单个订单详情
func GetOrder(ctx context.Context, c *app.RequestContext) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, utils.H{
			"error": "未授权",
		})
		return
	}

	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, utils.H{
			"error": "无效的订单ID",
		})
		return
	}

	var order models.Order
	err = database.DB.QueryRow(
		"SELECT id, user_id, total_price, status, created_at, updated_at FROM orders WHERE id = ? AND user_id = ?",
		orderID, userID,
	).Scan(&order.ID, &order.UserID, &order.TotalPrice, &order.Status, &order.CreatedAt, &order.UpdatedAt)

	if err == sql.ErrNoRows {
		c.JSON(404, utils.H{
			"error": "订单不存在",
		})
		return
	}
	if err != nil {
		c.JSON(500, utils.H{
			"error": "查询订单失败",
		})
		return
	}

	// 查询订单项
	itemRows, err := database.DB.Query(
		`SELECT oi.id, oi.order_id, oi.product_id, oi.quantity, oi.price, oi.created_at, oi.updated_at,
			p.id, p.name, p.description, p.price, p.image, p.stock, p.series, p.created_at, p.updated_at
		FROM order_items oi
		JOIN products p ON oi.product_id = p.id
		WHERE oi.order_id = ?`,
		order.ID,
	)
	if err == nil {
		defer itemRows.Close()
		for itemRows.Next() {
			var item models.OrderItem
			var product models.Product
			err := itemRows.Scan(
				&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.Price, &item.CreatedAt, &item.UpdatedAt,
				&product.ID, &product.Name, &product.Description, &product.Price, &product.Image, &product.Stock, &product.Series, &product.CreatedAt, &product.UpdatedAt,
			)
			if err == nil {
				item.Product = product
				order.Items = append(order.Items, item)
			}
		}
	}

	c.JSON(200, order)
}
