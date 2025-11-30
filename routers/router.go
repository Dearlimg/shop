package routers

import (
	"context"

	"shop/controller/api"
	"shop/middleware"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// InitRouter 初始化路由
func InitRouter(h *server.Hertz) {
	// CORS中间件（需要在所有路由之前）
	h.Use(func(ctx context.Context, c *app.RequestContext) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if string(c.Method()) == consts.MethodOptions {
			c.AbortWithStatus(consts.StatusNoContent)
			return
		}
		c.Next(ctx)
	})

	// 静态文件服务
	h.StaticFS("/static", &app.FS{Root: "./static", PathRewrite: app.NewPathSlashesStripper(1)})

	// API路由
	apiGroup := h.Group("/api")
	{
		// 公开路由
		apiGroup.POST("/register", api.Register)
		apiGroup.POST("/login", api.Login)
		apiGroup.GET("/products", api.GetProducts)
		apiGroup.GET("/products/:id", api.GetProduct)

		// 需要认证的路由
		authGroup := apiGroup.Group("/", middleware.AuthMiddleware())
		{
			// 购物车
			authGroup.GET("/cart", api.GetCart)
			authGroup.POST("/cart", api.AddToCart)
			authGroup.PATCH("/cart/:id/increment", api.IncrementCartItem) // 增量更新（+1/-1）
			authGroup.PUT("/cart/:id", api.UpdateCartItem)
			authGroup.DELETE("/cart/:id", api.DeleteCartItem)

			// 订单
			authGroup.POST("/orders", api.CreateOrder)
			authGroup.GET("/orders", api.GetOrders)
			authGroup.GET("/orders/:id", api.GetOrder)
		}
	}

	// 根路径重定向到前端
	h.GET("/", func(ctx context.Context, c *app.RequestContext) {
		c.Redirect(consts.StatusMovedPermanently, []byte("/static/index.html"))
	})

	// 测试页面
	h.GET("/test", func(ctx context.Context, c *app.RequestContext) {
		c.Redirect(consts.StatusMovedPermanently, []byte("/static/test.html"))
	})
}
