package api

import (
	"context"

	"shop/logic"
	"shop/model"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

// Register 用户注册
func Register(ctx context.Context, c *app.RequestContext) {
	var req model.RegisterRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(400, utils.H{
			"error": "请求参数错误: " + err.Error(),
		})
		return
	}

	userID, err := logic.Register(&req)
	if err != nil {
		statusCode := 500
		if err.Error() == "用户名已存在" {
			statusCode = 400
		}
		c.JSON(statusCode, utils.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, utils.H{
		"message": "注册成功",
		"user_id": userID,
	})
}

// Login 用户登录
func Login(ctx context.Context, c *app.RequestContext) {
	var req model.LoginRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(400, utils.H{
			"error": "请求参数错误: " + err.Error(),
		})
		return
	}

	resp, err := logic.Login(&req)
	if err != nil {
		statusCode := 500
		if err.Error() == "用户名或密码错误" {
			statusCode = 401
		}
		c.JSON(statusCode, utils.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, resp)
}
