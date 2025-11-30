package api

import (
	"context"

	"shop/logic"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

// GetProducts 获取商品列表
func GetProducts(ctx context.Context, c *app.RequestContext) {
	products, err := logic.GetProducts()
	if err != nil {
		c.JSON(500, utils.H{
			"error": "查询商品失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, utils.H{
		"products": products,
	})
}

// GetProduct 获取单个商品详情
func GetProduct(ctx context.Context, c *app.RequestContext) {
	productID := c.Param("id")

	product, err := logic.GetProduct(productID)
	if err != nil {
		c.JSON(500, utils.H{
			"error": "查询商品失败: " + err.Error(),
		})
		return
	}

	if product == nil {
		c.JSON(404, utils.H{
			"error": "商品不存在",
		})
		return
	}

	c.JSON(200, product)
}
