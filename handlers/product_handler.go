package handlers

import (
	"context"
	"database/sql"

	"shop/database"
	"shop/models"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

// GetProducts 获取商品列表
func GetProducts(ctx context.Context, c *app.RequestContext) {
	rows, err := database.DB.Query("SELECT id, name, description, price, image, stock, series, created_at, updated_at FROM products ORDER BY id DESC")
	if err != nil {
		c.JSON(500, utils.H{
			"error": "查询商品失败: " + err.Error(),
		})
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Image, &p.Stock, &p.Series, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			c.JSON(500, utils.H{
				"error": "读取商品数据失败",
			})
			return
		}
		products = append(products, p)
	}

	c.JSON(200, utils.H{
		"products": products,
	})
}

// GetProduct 获取单个商品详情
func GetProduct(ctx context.Context, c *app.RequestContext) {
	productID := c.Param("id")

	var product models.Product
	err := database.DB.QueryRow(
		"SELECT id, name, description, price, image, stock, series, created_at, updated_at FROM products WHERE id = ?",
		productID,
	).Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Image, &product.Stock, &product.Series, &product.CreatedAt, &product.UpdatedAt)

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

	c.JSON(200, product)
}
