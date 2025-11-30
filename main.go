package main

import (
	"log"
	"os"

	"shop/config"
	"shop/global/db"
	"shop/global/redis"
	"shop/routers"

	"github.com/cloudwego/hertz/pkg/app/server"
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
	if err := db.InitDB(dsn); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.CloseDB()

	// 初始化Redis
	if err := redis.InitRedis(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB); err != nil {
		log.Fatalf("Failed to initialize redis: %v", err)
	}
	defer redis.CloseRedis()

	// 创建表
	if err := db.CreateTables(); err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	// 初始化商品数据
	if err := db.SeedProducts(); err != nil {
		log.Printf("Warning: Failed to seed products: %v", err)
	}

	// 创建Hertz服务器
	serverAddr := cfg.Server.GetAddr()
	log.Printf("Server starting on %s", serverAddr)
	h := server.Default(server.WithHostPorts(serverAddr))

	// 初始化路由
	routers.InitRouter(h)

	h.Spin()
}
