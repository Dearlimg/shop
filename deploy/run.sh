#!/bin/bash

# 运行脚本
# 用于快速启动容器

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 检查配置文件是否存在
if [ ! -f "../config.yaml" ]; then
    echo -e "${RED}错误: 配置文件 config.yaml 不存在！${NC}"
    echo -e "${YELLOW}请先创建配置文件或复制 config.yaml.example${NC}"
    exit 1
fi

# 检查容器是否已存在
if [ "$(docker ps -a -q -f name=shop-app)" ]; then
    echo -e "${YELLOW}容器 shop-app 已存在，正在停止并删除...${NC}"
    docker stop shop-app > /dev/null 2>&1 || true
    docker rm shop-app > /dev/null 2>&1 || true
fi

# 检查镜像是否存在
if [ -z "$(docker images -q shop:latest)" ]; then
    echo -e "${YELLOW}镜像 shop:latest 不存在，正在构建...${NC}"
    ./build.sh
fi

echo -e "${GREEN}启动容器...${NC}"

# 运行容器
docker run -d \
    --name shop-app \
    --restart unless-stopped \
    -p 8080:8080 \
    -v $(pwd)/../config.yaml:/app/config.yaml:ro \
    -e CONFIG_PATH=/app/config.yaml \
    -e TZ=Asia/Shanghai \
    shop:latest

if [ $? -eq 0 ]; then
    echo -e "${GREEN}容器启动成功！${NC}"
    echo -e "\n${YELLOW}容器信息:${NC}"
    docker ps | grep shop-app
    
    echo -e "\n${GREEN}查看日志:${NC}"
    echo -e "docker logs -f shop-app"
    
    echo -e "\n${GREEN}停止容器:${NC}"
    echo -e "docker stop shop-app"
    
    echo -e "\n${GREEN}访问地址:${NC}"
    echo -e "http://localhost:8080"
else
    echo -e "${RED}容器启动失败！${NC}"
    exit 1
fi

