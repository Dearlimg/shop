#!/bin/bash

# 检查 MySQL 是否运行
echo "检查 MySQL 连接..."

# 设置数据库连接（可根据实际情况修改）
export DB_DSN="${DB_DSN:-root:password@tcp(localhost:3306)/shop?charset=utf8mb4&parseTime=True&loc=Local}"

# 运行程序
echo "启动服务器..."
go run main.go

