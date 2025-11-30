package middleware

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

// AuthMiddleware 认证中间件
func AuthMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		authHeader := string(c.Request.Header.Get("Authorization"))
		if authHeader == "" {
			c.JSON(401, utils.H{
				"error": "未授权，请先登录",
			})
			c.Abort()
			return
		}

		// 简单的token验证（实际项目中应使用JWT）
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(401, utils.H{
				"error": "无效的认证格式",
			})
			c.Abort()
			return
		}

		token := parts[1]
		// 从token中提取用户ID（简化版，实际应验证JWT）
		userID := extractUserIDFromToken(token)
		if userID == 0 {
			c.JSON(401, utils.H{
				"error": "无效的token",
			})
			c.Abort()
			return
		}

		// 将用户ID存储到上下文中
		c.Set("user_id", userID)
		c.Next(ctx)
	}
}

// extractUserIDFromToken 从token中提取用户ID（简化实现）
func extractUserIDFromToken(token string) int {
	// 这里简化处理，实际应该验证JWT token
	// 假设token格式为 "user_{id}"
	if strings.HasPrefix(token, "user_") {
		var userID int
		_, err := fmt.Sscanf(token, "user_%d", &userID)
		if err == nil {
			return userID
		}
	}
	return 0
}
