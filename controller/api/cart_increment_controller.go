package api

import (
	"context"
	"strconv"

	"shop/logic"
	"shop/model"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

// IncrementCartItem 增量更新购物车商品数量（+1 或 -1）
func IncrementCartItem(ctx context.Context, c *app.RequestContext) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, utils.H{
			"error": "未授权",
		})
		return
	}

	productID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, utils.H{
			"error": "无效的商品ID",
		})
		return
	}

	var req model.IncrementCartRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(400, utils.H{
			"error": "请求参数错误: " + err.Error(),
		})
		return
	}

	// 限制增量值范围（允许批量操作，但限制在合理范围内）
	if req.Delta == 0 {
		c.JSON(400, utils.H{
			"error": "增量值不能为 0",
		})
		return
	}
	if req.Delta > 100 || req.Delta < -100 {
		c.JSON(400, utils.H{
			"error": "增量值超出范围（-100 到 100）",
		})
		return
	}

	err = logic.IncrementCartItem(userID.(int), productID, req.Delta)
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
