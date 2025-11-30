package main

import (
	"context"
	"log"
	"os"

	"shop/config"
	"shop/database"
	"shop/handlers"
	"shop/middleware"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func main() {
	// 加载配置文件
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Printf("Loaded config from: %s", configPath)
	log.Printf("Database: %s@%s:%d/%s", cfg.Database.User, cfg.Database.Host, cfg.Database.Port, cfg.Database.Database)

	// 初始化数据库
	dsn := cfg.Database.GetDSN()
	if err := database.InitDB(dsn); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	// 创建表
	if err := database.CreateTables(); err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	// 初始化商品数据
	if err := database.SeedProducts(); err != nil {
		log.Printf("Warning: Failed to seed products: %v", err)
	}

	// 创建Hertz服务器
	serverAddr := cfg.Server.GetAddr()
	log.Printf("Server starting on %s", serverAddr)
	h := server.Default(server.WithHostPorts(serverAddr))

	// 静态文件服务
	h.StaticFS("/static", &app.FS{Root: "./static", PathRewrite: app.NewPathSlashesStripper(1)})

	// CORS中间件
	h.Use(func(ctx context.Context, c *app.RequestContext) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if string(c.Method()) == consts.MethodOptions {
			c.AbortWithStatus(consts.StatusNoContent)
			return
		}
		c.Next(ctx)
	})

	// API路由
	api := h.Group("/api")
	{
		// 公开路由
		api.POST("/register", handlers.Register)
		api.POST("/login", handlers.Login)
		api.GET("/products", handlers.GetProducts)
		api.GET("/products/:id", handlers.GetProduct)

		// 需要认证的路由
		auth := api.Group("/", middleware.AuthMiddleware())
		{
			// 购物车
			auth.GET("/cart", handlers.GetCart)
			auth.POST("/cart", handlers.AddToCart)
			auth.PUT("/cart/:id", handlers.UpdateCartItem)
			auth.DELETE("/cart/:id", handlers.DeleteCartItem)

			// 订单
			auth.POST("/orders", handlers.CreateOrder)
			auth.GET("/orders", handlers.GetOrders)
			auth.GET("/orders/:id", handlers.GetOrder)
		}
	}

	// 根路径重定向到前端
	h.GET("/", func(ctx context.Context, c *app.RequestContext) {
		c.Redirect(consts.StatusMovedPermanently, []byte("/static/index.html"))
	})

	h.Spin()
}
