package handlers

import (
	"context"
	"database/sql"
	"fmt"

	"shop/database"
	"shop/models"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"golang.org/x/crypto/bcrypt"
)

// Register 用户注册
func Register(ctx context.Context, c *app.RequestContext) {
	var req models.RegisterRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(400, utils.H{
			"error": "请求参数错误: " + err.Error(),
		})
		return
	}

	// 检查用户名是否已存在
	var count int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", req.Username).Scan(&count)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(500, utils.H{
			"error": "数据库查询错误",
		})
		return
	}
	if count > 0 {
		c.JSON(400, utils.H{
			"error": "用户名已存在",
		})
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, utils.H{
			"error": "密码加密失败",
		})
		return
	}

	// 插入用户
	result, err := database.DB.Exec(
		"INSERT INTO users (username, password, email) VALUES (?, ?, ?)",
		req.Username, string(hashedPassword), req.Email,
	)
	if err != nil {
		c.JSON(500, utils.H{
			"error": "注册失败: " + err.Error(),
		})
		return
	}

	userID, _ := result.LastInsertId()
	c.JSON(200, utils.H{
		"message": "注册成功",
		"user_id": userID,
	})
}

// Login 用户登录
func Login(ctx context.Context, c *app.RequestContext) {
	var req models.LoginRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(400, utils.H{
			"error": "请求参数错误: " + err.Error(),
		})
		return
	}

	// 查询用户
	var user models.User
	var passwordHash string
	err := database.DB.QueryRow(
		"SELECT id, username, password, email, created_at, updated_at FROM users WHERE username = ?",
		req.Username,
	).Scan(&user.ID, &user.Username, &passwordHash, &user.Email, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		c.JSON(401, utils.H{
			"error": "用户名或密码错误",
		})
		return
	}
	if err != nil {
		c.JSON(500, utils.H{
			"error": "数据库查询错误",
		})
		return
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password))
	if err != nil {
		c.JSON(401, utils.H{
			"error": "用户名或密码错误",
		})
		return
	}

	// 生成token（简化版，实际应使用JWT）
	token := fmt.Sprintf("user_%d", user.ID)

	c.JSON(200, models.LoginResponse{
		Token: token,
		User:  user,
	})
}
