package api

import (
	"context"
	"strconv"

	"shop/logic"
	"shop/model"

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

	var req model.CreateOrderRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(400, utils.H{
			"error": "请求参数错误: " + err.Error(),
		})
		return
	}

	orderID, totalPrice, err := logic.CreateOrder(userID.(int), &req)
	if err != nil {
		statusCode := 500
		if err.Error() == "购物车项不存在" {
			statusCode = 404
		} else if err.Error() == "商品库存不足" || err.Error() == "库存不足" {
			statusCode = 400
		}
		c.JSON(statusCode, utils.H{
			"error": err.Error(),
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

	orders, err := logic.GetOrders(userID.(int))
	if err != nil {
		c.JSON(500, utils.H{
			"error": "查询订单失败: " + err.Error(),
		})
		return
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

	order, err := logic.GetOrder(userID.(int), orderID)
	if err != nil {
		statusCode := 500
		if err.Error() == "订单不存在" {
			statusCode = 404
		}
		c.JSON(statusCode, utils.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, order)
}
