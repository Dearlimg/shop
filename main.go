package main

import (
	"fmt"
	"log"

	"shop/global/db"
	"shop/global/redis"
	"shop/routers"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	// 硬编码配置信息
	databaseHost := "47.118.19.28"
	databasePort := 3307
	databaseUser := "root"
	databasePassword := "sta_go"
	databaseName := "durlim"
	databaseCharset := "utf8mb4"

	redisAddr := "47.118.19.28:6379"
	redisPassword := "sta_go"
	redisDB := 0

	serverHost := "0.0.0.0"
	serverPort := 8080

	log.Printf("Database: %s@%s:%d/%s", databaseUser, databaseHost, databasePort, databaseName)
	log.Printf("Redis: %s", redisAddr)

	// 初始化数据库
	dsn := getDSN(databaseUser, databasePassword, databaseHost, databasePort, databaseName, databaseCharset)
	if err := db.InitDB(dsn); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.CloseDB()

	// 初始化Redis
	if err := redis.InitRedis(redisAddr, redisPassword, redisDB); err != nil {
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
	serverAddr := getServerAddr(serverHost, serverPort)
	log.Printf("Server starting on %s", serverAddr)
	h := server.Default(server.WithHostPorts(serverAddr))

	// 初始化路由
	routers.InitRouter(h)

	h.Spin()
}

// getDSN 生成数据库连接字符串
func getDSN(user, password, host string, port int, database, charset string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		user, password, host, port, database, charset)
}

// getServerAddr 生成服务器地址
func getServerAddr(host string, port int) string {
	return fmt.Sprintf("%s:%d", host, port)
}
