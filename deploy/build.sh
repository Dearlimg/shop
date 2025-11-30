#!/bin/bash

# 构建脚本
# 用于快速构建 Docker 镜像

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}开始构建 Docker 镜像...${NC}"

# 获取版本号（可选）
VERSION=${1:-latest}

# 构建镜像
echo -e "${YELLOW}构建镜像: shop:${VERSION}${NC}"
docker build -t shop:${VERSION} -t shop:latest ..

if [ $? -eq 0 ]; then
    echo -e "${GREEN}构建成功！${NC}"
    echo -e "${GREEN}镜像标签: shop:${VERSION}, shop:latest${NC}"
    
    # 显示镜像信息
    echo -e "\n${YELLOW}镜像信息:${NC}"
    docker images | grep shop
    
    echo -e "\n${GREEN}可以使用以下命令运行容器:${NC}"
    echo -e "docker run -d --name shop-app -p 8080:8080 -v \$(pwd)/config.yaml:/app/config.yaml:ro shop:latest"
else
    echo -e "${RED}构建失败！${NC}"
    exit 1
fi

