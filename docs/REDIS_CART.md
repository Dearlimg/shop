# Redis 购物车实现说明

## 为什么使用 Redis 实现购物车？

### 优势

1. **高性能**
   - Redis 是内存数据库，读写速度极快
   - 购物车操作频繁，使用 Redis 可以显著提升性能

2. **适合临时数据**
   - 购物车通常是临时数据，不需要永久存储
   - Redis 支持设置过期时间，自动清理过期数据

3. **减少数据库压力**
   - 购物车操作（添加、删除、修改）非常频繁
   - 使用 Redis 可以减轻 MySQL 数据库的负担

4. **更好的用户体验**
   - 更快的响应速度
   - 支持高并发访问

### 对比

| 特性 | MySQL | Redis |
|------|-------|-------|
| 存储位置 | 磁盘 | 内存 |
| 读写速度 | 较慢 | 极快 |
| 数据持久化 | 永久 | 可设置过期 |
| 适合场景 | 永久数据 | 临时/缓存数据 |

## 实现方案

### 数据结构

购物车在 Redis 中的存储结构：

```
Key: cart:user:{user_id}:product:{product_id}
Value: Hash {
  product_id: {product_id}
  quantity: {quantity}
  user_id: {user_id}
}
```

### 过期时间

- 默认过期时间：**30天**
- 每次操作购物车时，会自动续期

### 主要改动

1. **DAO 层** (`dao/cart_redis_dao.go`)
   - 使用 Redis Hash 存储购物车项
   - 使用 `cart:user:{user_id}:product:{product_id}` 作为键

2. **Logic 层** (`logic/cart_logic.go`)
   - 调用 Redis DAO 方法
   - 参数改为使用 `product_id` 而不是 `cart_item_id`

3. **Controller 层** (`controller/api/cart_controller.go`)
   - API 参数改为 `product_id`
   - 路径参数 `/cart/:id` 中的 `id` 现在是 `product_id`

4. **前端** (`static/index.html`)
   - 更新购物车操作，使用 `product_id` 而不是 `cart_item_id`

## 配置说明

在 `config.yaml` 中添加 Redis 配置：

```yaml
redis:
  addr: localhost:6379      # Redis地址
  password: ""              # Redis密码，无密码留空
  db: 0                     # Redis数据库编号
```

## 使用说明

### 1. 安装 Redis

**macOS:**
```bash
brew install redis
brew services start redis
```

**Linux:**
```bash
sudo apt-get install redis-server
sudo systemctl start redis
```

**Docker:**
```bash
docker run -d -p 6379:6379 redis:latest
```

### 2. 验证 Redis 连接

```bash
redis-cli ping
# 应该返回: PONG
```

### 3. 运行程序

程序启动时会自动连接 Redis，如果连接失败会报错。

## API 变更说明

### 重要变更

**注意**: Redis 版本的购物车 API 参数有所变化：

1. **更新购物车项** (`PUT /api/cart/:id`)
   - 之前：`id` 是 `cart_item_id`
   - 现在：`id` 是 `product_id`

2. **删除购物车项** (`DELETE /api/cart/:id`)
   - 之前：`id` 是 `cart_item_id`
   - 现在：`id` 是 `product_id`

3. **创建订单** (`POST /api/orders`)
   - 之前：`cart_item_ids` 是购物车项ID数组
   - 现在：`cart_item_ids` 是商品ID数组（实际是 `product_ids`）

### 前端适配

前端代码已更新，使用 `product_id` 进行操作。

## 数据迁移

如果之前使用 MySQL 存储购物车，需要：

1. **保留旧代码**：`dao/cart_dao.go` 和 `logic/cart_logic.go` 已保留
2. **切换实现**：修改 `logic/cart_logic.go` 中的函数调用即可

## 性能对比

### MySQL 版本
- 每次操作需要：查询数据库 → 更新数据库 → 返回结果
- 响应时间：~10-50ms

### Redis 版本
- 每次操作需要：操作内存 → 返回结果
- 响应时间：~1-5ms

**性能提升：约 5-10 倍**

## 注意事项

1. **数据持久化**
   - Redis 数据默认在内存中，重启可能丢失
   - 可以配置 Redis 持久化（RDB/AOF）
   - 购物车是临时数据，丢失影响不大

2. **内存管理**
   - 设置合理的过期时间（当前 30 天）
   - 监控 Redis 内存使用情况

3. **高可用**
   - 生产环境建议使用 Redis 集群或哨兵模式
   - 配置主从复制

4. **数据一致性**
   - 商品信息仍从 MySQL 读取
   - 购物车数量存储在 Redis
   - 创建订单时从 Redis 读取并写入 MySQL

## 监控建议

1. **Redis 连接数**
2. **内存使用率**
3. **命令执行时间**
4. **键的数量**（购物车数量）

## 故障处理

如果 Redis 连接失败：
- 程序启动时会报错并退出
- 检查 Redis 服务是否运行
- 检查配置文件中的地址和端口

如果 Redis 数据丢失：
- 用户需要重新添加商品到购物车
- 不影响已创建的订单

