# 部署说明

本文档说明如何使用 Docker 部署购物车系统。

## 前置要求

- Docker (版本 20.10+)
- Docker Compose (版本 2.0+)
- MySQL 数据库（已配置，无需 Docker 管理）
- Redis（已配置，无需 Docker 管理）

## 快速开始

> **注意**: 配置信息已硬编码到程序中，如需修改配置，请直接修改 `main.go` 文件中的配置值后重新构建镜像。

### 构建和运行

#### 方式一：使用 Docker Compose（推荐）

```bash
# 进入 deploy 目录
cd deploy

# 构建并启动
docker-compose up -d --build

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

#### 方式二：使用 Docker 命令

```bash
# 构建镜像
docker build -t shop:latest .

# 运行容器
docker run -d \
  --name shop-app \
  -p 8080:8080 \
  shop:latest

# 查看日志
docker logs -f shop-app

# 停止容器
docker stop shop-app
docker rm shop-app
```

#### 方式三：使用部署脚本

```bash
cd deploy
./build.sh    # 构建镜像
./run.sh      # 运行容器
```

## 配置说明

### 修改配置

配置信息已硬编码在 `main.go` 文件中，如需修改配置：

1. 编辑 `main.go` 文件中的配置值：
   ```go
   databaseHost := "your-mysql-host"
   databasePort := 3306
   databaseUser := "root"
   databasePassword := "your-password"
   databaseName := "your-database"
   
   redisAddr := "your-redis-host:6379"
   redisPassword := "your-redis-password"
   ```

2. 重新构建镜像：
   ```bash
   docker build -t shop:latest .
   ```

### 环境变量

- `TZ`: 时区设置（默认：Asia/Shanghai）

### 挂载卷

- `static`: 静态文件目录（可选，如果需要动态更新）

## 健康检查

容器包含健康检查，可以通过以下命令查看：

```bash
docker inspect shop-app | grep -A 10 Health
```

## 查看日志

```bash
# Docker Compose
docker-compose logs -f shop

# Docker
docker logs -f shop-app
```

## 更新部署

```bash
# 停止旧容器
docker-compose down

# 重新构建
docker-compose build --no-cache

# 启动新容器
docker-compose up -d
```

## 生产环境建议

1. **使用环境变量管理敏感信息**：
   ```bash
   docker run -d \
     --name shop-app \
     -p 8080:8080 \
     -e DB_HOST=your-host \
     -e DB_PASSWORD=your-password \
     shop:latest
   ```

2. **使用 Docker Secrets**（Docker Swarm）或 Kubernetes Secrets

3. **配置反向代理**（Nginx/Traefik）：
   ```nginx
   server {
       listen 80;
       server_name your-domain.com;
       
       location / {
           proxy_pass http://localhost:8080;
           proxy_set_header Host $host;
           proxy_set_header X-Real-IP $remote_addr;
       }
   }
   ```

4. **启用 HTTPS**：使用 Let's Encrypt 或配置 SSL 证书

5. **监控和日志**：集成 Prometheus、Grafana 等监控工具

6. **资源限制**：
   ```yaml
   deploy:
     resources:
       limits:
         cpus: '1'
         memory: 512M
       reservations:
         cpus: '0.5'
         memory: 256M
   ```

## 故障排查

### 容器无法启动

1. 检查配置文件是否正确：
   ```bash
   docker exec shop-app cat /app/config.yaml
   ```

2. 查看容器日志：
   ```bash
   docker logs shop-app
   ```

### 无法连接数据库/Redis

1. 确认数据库和 Redis 地址可从容器内访问
2. 检查防火墙规则
3. 验证配置文件中的连接信息

### 端口冲突

如果 8080 端口被占用，修改 `docker-compose.yml` 中的端口映射：
```yaml
ports:
  - "8081:8080"  # 将主机端口改为 8081
```

## 清理

```bash
# 停止并删除容器
docker-compose down

# 删除镜像
docker rmi shop:latest

# 清理未使用的资源
docker system prune -a
```

