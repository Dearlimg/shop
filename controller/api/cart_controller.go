package api

import (
	"context"
	"strconv"

	"shop/logic"
	"shop/model"

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

	items, err := logic.GetCart(userID.(int))
	if err != nil {
		c.JSON(500, utils.H{
			"error": "查询购物车失败: " + err.Error(),
		})
		return
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

	var req model.AddToCartRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(400, utils.H{
			"error": "请求参数错误: " + err.Error(),
		})
		return
	}

	err := logic.AddToCart(userID.(int), &req)
	if err != nil {
		statusCode := 500
		if err.Error() == "库存不足" {
			statusCode = 400
		}
		c.JSON(statusCode, utils.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, utils.H{
		"message": "添加到购物车成功",
	})
}

// UpdateCartItem 更新购物车商品数量（使用Redis，参数改为product_id）
func UpdateCartItem(ctx context.Context, c *app.RequestContext) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, utils.H{
			"error": "未授权",
		})
		return
	}

	// Redis版本使用product_id而不是cart_item_id
	productID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, utils.H{
			"error": "无效的商品ID",
		})
		return
	}

	var req model.UpdateCartRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(400, utils.H{
			"error": "请求参数错误: " + err.Error(),
		})
		return
	}

	err = logic.UpdateCartItem(userID.(int), productID, &req)
	if err != nil {
		statusCode := 500
		if err.Error() == "库存不足" || err.Error() == "购物车项不存在" {
			statusCode = 400
		}
		c.JSON(statusCode, utils.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, utils.H{
		"message": "更新成功",
	})
}

// DeleteCartItem 删除购物车商品（使用Redis，参数改为product_id）
func DeleteCartItem(ctx context.Context, c *app.RequestContext) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, utils.H{
			"error": "未授权",
		})
		return
	}

	// Redis版本使用product_id而不是cart_item_id
	productID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, utils.H{
			"error": "无效的商品ID",
		})
		return
	}

	err = logic.DeleteCartItem(userID.(int), productID)
	if err != nil {
		statusCode := 500
		if err.Error() == "购物车项不存在" {
			statusCode = 404
		}
		c.JSON(statusCode, utils.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, utils.H{
		"message": "删除成功",
	})
}
