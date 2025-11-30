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

// GetCart 获取购物车
func GetCart(ctx context.Context, c *app.RequestContext) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, utils.H{
			"error": "未授权",
		})
		return
	}

	rows, err := database.DB.Query(
		`SELECT ci.id, ci.user_id, ci.product_id, ci.quantity, ci.created_at, ci.updated_at,
			p.id, p.name, p.description, p.price, p.image, p.stock, p.series, p.created_at, p.updated_at
		FROM cart_items ci
		JOIN products p ON ci.product_id = p.id
		WHERE ci.user_id = ?`,
		userID,
	)
	if err != nil {
		c.JSON(500, utils.H{
			"error": "查询购物车失败: " + err.Error(),
		})
		return
	}
	defer rows.Close()

	var items []models.CartItem
	for rows.Next() {
		var item models.CartItem
		var product models.Product
		err := rows.Scan(
			&item.ID, &item.UserID, &item.ProductID, &item.Quantity, &item.CreatedAt, &item.UpdatedAt,
			&product.ID, &product.Name, &product.Description, &product.Price, &product.Image, &product.Stock, &product.Series, &product.CreatedAt, &product.UpdatedAt,
		)
		if err != nil {
			c.JSON(500, utils.H{
				"error": "读取购物车数据失败",
			})
			return
		}
		item.Product = product
		items = append(items, item)
	}

	c.JSON(200, utils.H{
		"items": items,
	})
}

// AddToCart 添加到购物车
func AddToCart(ctx context.Context, c *app.RequestContext) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, utils.H{
			"error": "未授权",
		})
		return
	}

	var req models.AddToCartRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(400, utils.H{
			"error": "请求参数错误: " + err.Error(),
		})
		return
	}

	// 检查商品是否存在
	var stock int
	err := database.DB.QueryRow("SELECT stock FROM products WHERE id = ?", req.ProductID).Scan(&stock)
	if err == sql.ErrNoRows {
		c.JSON(404, utils.H{
			"error": "商品不存在",
		})
		return
	}
	if err != nil {
		c.JSON(500, utils.H{
			"error": "查询商品失败",
		})
		return
	}

	if stock < req.Quantity {
		c.JSON(400, utils.H{
			"error": "库存不足",
		})
		return
	}

	// 检查购物车中是否已有该商品
	var existingID int
	var existingQuantity int
	err = database.DB.QueryRow(
		"SELECT id, quantity FROM cart_items WHERE user_id = ? AND product_id = ?",
		userID, req.ProductID,
	).Scan(&existingID, &existingQuantity)

	if err == nil {
		// 更新数量
		newQuantity := existingQuantity + req.Quantity
		if newQuantity > stock {
			c.JSON(400, utils.H{
				"error": "库存不足",
			})
			return
		}
		_, err = database.DB.Exec(
			"UPDATE cart_items SET quantity = ? WHERE id = ?",
			newQuantity, existingID,
		)
		if err != nil {
			c.JSON(500, utils.H{
				"error": "更新购物车失败",
			})
			return
		}
		c.JSON(200, utils.H{
			"message": "购物车更新成功",
		})
		return
	}

	// 添加新商品到购物车
	_, err = database.DB.Exec(
		"INSERT INTO cart_items (user_id, product_id, quantity) VALUES (?, ?, ?)",
		userID, req.ProductID, req.Quantity,
	)
	if err != nil {
		c.JSON(500, utils.H{
			"error": "添加到购物车失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, utils.H{
		"message": "添加到购物车成功",
	})
}

// UpdateCartItem 更新购物车商品数量
func UpdateCartItem(ctx context.Context, c *app.RequestContext) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, utils.H{
			"error": "未授权",
		})
		return
	}

	cartItemID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, utils.H{
			"error": "无效的购物车项ID",
		})
		return
	}

	var req models.UpdateCartRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(400, utils.H{
			"error": "请求参数错误: " + err.Error(),
		})
		return
	}

	// 检查购物车项是否属于当前用户
	var productID int
	err = database.DB.QueryRow(
		"SELECT product_id FROM cart_items WHERE id = ? AND user_id = ?",
		cartItemID, userID,
	).Scan(&productID)
	if err == sql.ErrNoRows {
		c.JSON(404, utils.H{
			"error": "购物车项不存在",
		})
		return
	}

	// 检查库存
	var stock int
	err = database.DB.QueryRow("SELECT stock FROM products WHERE id = ?", productID).Scan(&stock)
	if err != nil {
		c.JSON(500, utils.H{
			"error": "查询商品失败",
		})
		return
	}

	if req.Quantity > stock {
		c.JSON(400, utils.H{
			"error": "库存不足",
		})
		return
	}

	// 更新数量
	_, err = database.DB.Exec(
		"UPDATE cart_items SET quantity = ? WHERE id = ?",
		req.Quantity, cartItemID,
	)
	if err != nil {
		c.JSON(500, utils.H{
			"error": "更新购物车失败",
		})
		return
	}

	c.JSON(200, utils.H{
		"message": "更新成功",
	})
}

// DeleteCartItem 删除购物车商品
func DeleteCartItem(ctx context.Context, c *app.RequestContext) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, utils.H{
			"error": "未授权",
		})
		return
	}

	cartItemID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, utils.H{
			"error": "无效的购物车项ID",
		})
		return
	}

	// 检查购物车项是否属于当前用户
	var count int
	err = database.DB.QueryRow(
		"SELECT COUNT(*) FROM cart_items WHERE id = ? AND user_id = ?",
		cartItemID, userID,
	).Scan(&count)
	if err != nil || count == 0 {
		c.JSON(404, utils.H{
			"error": "购物车项不存在",
		})
		return
	}

	// 删除
	_, err = database.DB.Exec("DELETE FROM cart_items WHERE id = ?", cartItemID)
	if err != nil {
		c.JSON(500, utils.H{
			"error": "删除失败",
		})
		return
	}

	c.JSON(200, utils.H{
		"message": "删除成功",
	})
}
