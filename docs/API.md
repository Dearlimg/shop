# 购物车系统 API 接口文档

## 基础信息

- **Base URL**: `http://localhost:8080/api`
- **Content-Type**: `application/json`
- **认证方式**: Bearer Token（在请求头中添加 `Authorization: Bearer {token}`）

## 认证说明

大部分接口需要用户登录后才能访问。登录成功后，服务器会返回一个 `token`，后续请求需要在请求头中携带：

```
Authorization: Bearer {token}
```

---

## 1. 用户相关接口

### 1.1 用户注册

**接口地址**: `POST /api/register`

**接口描述**: 注册新用户

**请求参数**:

```json
{
  "username": "string",  // 必填，用户名
  "password": "string",   // 必填，密码（至少6位）
  "email": "string"      // 必填，邮箱地址
}
```

**响应示例**:

```json
{
  "message": "注册成功",
  "user_id": 1
}
```

**错误响应**:

```json
{
  "error": "用户名已存在"
}
```

**状态码**:
- `200`: 注册成功
- `400`: 请求参数错误或用户名已存在

---

### 1.2 用户登录

**接口地址**: `POST /api/login`

**接口描述**: 用户登录，获取访问令牌

**请求参数**:

```json
{
  "username": "string",  // 必填，用户名
  "password": "string"   // 必填，密码
}
```

**响应示例**:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

**错误响应**:

```json
{
  "error": "用户名或密码错误"
}
```

**状态码**:
- `200`: 登录成功
- `400`: 请求参数错误
- `401`: 用户名或密码错误

---

## 2. 商品相关接口

### 2.1 获取商品列表

**接口地址**: `GET /api/products`

**接口描述**: 获取所有商品列表（无需认证）

**请求参数**: 无

**响应示例**:

```json
{
  "products": [
    {
      "id": 1,
      "name": "拉布布盲盒-经典款",
      "description": "经典款拉布布盲盒",
      "price": 59.00,
      "image": "https://example.com/image.jpg",
      "stock": 100,
      "series": "拉布布",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

**状态码**:
- `200`: 查询成功

---

### 2.2 获取商品详情

**接口地址**: `GET /api/products/:id`

**接口描述**: 获取单个商品的详细信息（无需认证）

**路径参数**:
- `id`: 商品ID（整数）

**响应示例**:

```json
{
  "id": 1,
  "name": "拉布布盲盒-经典款",
  "description": "经典款拉布布盲盒",
  "price": 59.00,
  "image": "https://example.com/image.jpg",
  "stock": 100,
  "series": "拉布布",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

**错误响应**:

```json
{
  "error": "商品不存在"
}
```

**状态码**:
- `200`: 查询成功
- `404`: 商品不存在

---

## 3. 购物车相关接口

> ⚠️ **注意**: 以下所有接口都需要认证（在请求头中携带 token）

### 3.1 获取购物车

**接口地址**: `GET /api/cart`

**接口描述**: 获取当前用户的购物车中所有商品

**请求参数**: 无

**请求头**:
```
Authorization: Bearer {token}
```

**响应示例**:

```json
{
  "items": [
    {
      "user_id": 1,
      "product_id": 1,
      "quantity": 2,
      "product": {
        "id": 1,
        "name": "拉布布盲盒-经典款",
        "description": "经典款拉布布盲盒",
        "price": 59.00,
        "image": "https://example.com/image.jpg",
        "stock": 100,
        "series": "拉布布",
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    }
  ]
}
```

**状态码**:
- `200`: 查询成功
- `401`: 未授权（未登录或token无效）

---

### 3.2 添加商品到购物车

**接口地址**: `POST /api/cart`

**接口描述**: 将商品添加到购物车，如果商品已存在则增加数量

**请求参数**:

```json
{
  "product_id": 1,  // 必填，商品ID
  "quantity": 1     // 必填，数量（至少为1）
}
```

**请求头**:
```
Authorization: Bearer {token}
```

**响应示例**:

```json
{
  "message": "添加到购物车成功"
}
```

**错误响应**:

```json
{
  "error": "库存不足"
}
```

**状态码**:
- `200`: 添加成功
- `400`: 请求参数错误或库存不足
- `401`: 未授权

---

### 3.3 增量更新购物车商品数量

**接口地址**: `PATCH /api/cart/:id/increment`

**接口描述**: 对购物车中的商品进行增量更新（+1 或 -1），推荐使用此接口实现购物车的 +1/-1 功能

**路径参数**:
- `id`: 商品ID（整数）

**请求参数**:

```json
{
  "delta": 1  // 必填，增量值（可以是正数或负数，范围：-100 到 100）
}
```

**请求头**:
```
Authorization: Bearer {token}
```

**响应示例**:

```json
{
  "message": "更新成功"
}
```

**错误响应**:

```json
{
  "error": "库存不足，当前库存: 10"
}
```

**状态码**:
- `200`: 更新成功
- `400`: 请求参数错误、库存不足或购物车项不存在
- `401`: 未授权

**使用场景**: 
- 点击 `+` 按钮：`{ "delta": 1 }`
- 点击 `-` 按钮：`{ "delta": -1 }`
- 如果数量减到 0 或以下，会自动删除该购物车项

---

### 3.4 更新购物车商品数量

**接口地址**: `PUT /api/cart/:id`

**接口描述**: 直接设置购物车中商品的数量

**路径参数**:
- `id`: 商品ID（整数）

**请求参数**:

```json
{
  "quantity": 3  // 必填，新的数量（至少为1）
}
```

**请求头**:
```
Authorization: Bearer {token}
```

**响应示例**:

```json
{
  "message": "更新成功"
}
```

**错误响应**:

```json
{
  "error": "库存不足"
}
```

**状态码**:
- `200`: 更新成功
- `400`: 请求参数错误、库存不足或购物车项不存在
- `401`: 未授权

---

### 3.5 删除购物车商品

**接口地址**: `DELETE /api/cart/:id`

**接口描述**: 从购物车中删除指定商品

**路径参数**:
- `id`: 商品ID（整数）

**请求参数**: 无

**请求头**:
```
Authorization: Bearer {token}
```

**响应示例**:

```json
{
  "message": "删除成功"
}
```

**错误响应**:

```json
{
  "error": "购物车项不存在"
}
```

**状态码**:
- `200`: 删除成功
- `404`: 购物车项不存在
- `401`: 未授权

---

## 4. 订单相关接口

> ⚠️ **注意**: 以下所有接口都需要认证（在请求头中携带 token）

### 4.1 创建订单

**接口地址**: `POST /api/orders`

**接口描述**: 创建订单，默认使用购物车中所有商品进行结算

**请求参数**（可选）:

```json
{
  "cart_item_ids": [1, 2, 3]  // 可选，商品ID数组。如果不提供或为空，则使用购物车中所有商品
}
```

**请求头**:
```
Authorization: Bearer {token}
```

**响应示例**:

```json
{
  "message": "订单创建成功",
  "order_id": 1,
  "total_price": 118.00
}
```

**错误响应**:

```json
{
  "error": "购物车为空"
}
```

或

```json
{
  "error": "商品库存不足: 拉布布盲盒-经典款 (需要: 5, 库存: 3)"
}
```

**状态码**:
- `200`: 订单创建成功
- `400`: 购物车为空、库存不足或请求参数错误
- `401`: 未授权
- `404`: 购物车项不存在

**重要说明**:
- 如果请求体为空 `{}` 或不提供 `cart_item_ids`，系统会自动使用购物车中所有商品进行结算
- 订单创建成功后，购物车中对应的商品会被自动删除
- 商品库存会在订单创建时自动扣减

---

### 4.2 获取订单列表

**接口地址**: `GET /api/orders`

**接口描述**: 获取当前用户的所有订单历史

**请求参数**: 无

**请求头**:
```
Authorization: Bearer {token}
```

**响应示例**:

```json
{
  "orders": [
    {
      "id": 1,
      "user_id": 1,
      "total_price": 118.00,
      "status": "pending",
      "items": [
        {
          "id": 1,
          "order_id": 1,
          "product_id": 1,
          "quantity": 2,
          "price": 59.00,
          "product": {
            "id": 1,
            "name": "拉布布盲盒-经典款",
            "description": "经典款拉布布盲盒",
            "price": 59.00,
            "image": "https://example.com/image.jpg",
            "stock": 98,
            "series": "拉布布",
            "created_at": "2024-01-01T00:00:00Z",
            "updated_at": "2024-01-01T00:00:00Z"
          },
          "created_at": "2024-01-01T00:00:00Z",
          "updated_at": "2024-01-01T00:00:00Z"
        }
      ],
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

**状态码**:
- `200`: 查询成功
- `401`: 未授权

---

### 4.3 获取订单详情

**接口地址**: `GET /api/orders/:id`

**接口描述**: 获取指定订单的详细信息

**路径参数**:
- `id`: 订单ID（整数）

**请求参数**: 无

**请求头**:
```
Authorization: Bearer {token}
```

**响应示例**:

```json
{
  "id": 1,
  "user_id": 1,
  "total_price": 118.00,
  "status": "pending",
  "items": [
    {
      "id": 1,
      "order_id": 1,
      "product_id": 1,
      "quantity": 2,
      "price": 59.00,
      "product": {
        "id": 1,
        "name": "拉布布盲盒-经典款",
        "description": "经典款拉布布盲盒",
        "price": 59.00,
        "image": "https://example.com/image.jpg",
        "stock": 98,
        "series": "拉布布",
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      },
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

**错误响应**:

```json
{
  "error": "订单不存在"
}
```

**状态码**:
- `200`: 查询成功
- `404`: 订单不存在
- `401`: 未授权

---

## 5. 数据模型

### 5.1 User（用户）

```json
{
  "id": 1,
  "username": "testuser",
  "email": "test@example.com",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### 5.2 Product（商品）

```json
{
  "id": 1,
  "name": "拉布布盲盒-经典款",
  "description": "经典款拉布布盲盒",
  "price": 59.00,
  "image": "https://example.com/image.jpg",
  "stock": 100,
  "series": "拉布布",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### 5.3 CartItem（购物车项）

```json
{
  "user_id": 1,
  "product_id": 1,
  "quantity": 2,
  "product": {
    // Product 对象（见上方）
  }
}
```

### 5.4 Order（订单）

```json
{
  "id": 1,
  "user_id": 1,
  "total_price": 118.00,
  "status": "pending",
  "items": [
    // OrderItem 数组
  ],
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### 5.5 OrderItem（订单项）

```json
{
  "id": 1,
  "order_id": 1,
  "product_id": 1,
  "quantity": 2,
  "price": 59.00,
  "product": {
    // Product 对象（见上方）
  },
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

---

## 6. 错误码说明

| 状态码 | 说明 |
|--------|------|
| 200 | 请求成功 |
| 400 | 请求参数错误、业务逻辑错误（如库存不足、用户名已存在等） |
| 401 | 未授权（未登录或token无效） |
| 404 | 资源不存在（商品、订单、购物车项等） |
| 500 | 服务器内部错误 |

---

## 7. 请求示例

### 7.1 使用 curl 测试

#### 注册用户
```bash
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "123456",
    "email": "test@example.com"
  }'
```

#### 用户登录
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "123456"
  }'
```

#### 获取商品列表
```bash
curl -X GET http://localhost:8080/api/products
```

#### 添加到购物车
```bash
curl -X POST http://localhost:8080/api/cart \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "product_id": 1,
    "quantity": 1
  }'
```

#### 增量更新购物车（+1）
```bash
curl -X PATCH http://localhost:8080/api/cart/1/increment \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "delta": 1
  }'
```

#### 创建订单（使用购物车中所有商品）
```bash
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{}'
```

---

## 8. 注意事项

1. **购物车机制**：
   - 每个用户只有一个购物车
   - 购物车使用 Redis 存储，30天自动过期
   - 购物车中的商品以 `product_id` 为唯一标识

2. **下单流程**：
   - 默认使用购物车中所有商品进行结算
   - 订单创建成功后，购物车中对应的商品会被自动删除
   - 商品库存会在订单创建时自动扣减

3. **库存检查**：
   - 添加到购物车时会检查库存
   - 更新购物车数量时会检查库存
   - 创建订单时会再次检查库存（防止并发问题）

4. **认证 Token**：
   - Token 通过登录接口获取
   - Token 需要保存在前端（如 localStorage）
   - 所有需要认证的接口都需要在请求头中携带 token

5. **CORS 支持**：
   - 后端已配置 CORS，支持跨域请求
   - 允许的请求方法：GET, POST, PUT, PATCH, DELETE, OPTIONS

---

## 9. 常见问题

**Q: 如何实现购物车的 +1/-1 功能？**  
A: 使用 `PATCH /api/cart/:id/increment` 接口，传递 `delta: 1` 或 `delta: -1`。

**Q: 下单时如何选择部分商品？**  
A: 在创建订单时，传递 `cart_item_ids` 数组，指定要结算的商品ID。如果不传递或为空，则使用购物车中所有商品。

**Q: 购物车数据会丢失吗？**  
A: 购物车数据存储在 Redis 中，30天自动过期。如果 Redis 服务重启，数据可能会丢失（建议生产环境配置 Redis 持久化）。

**Q: 如何判断用户是否已登录？**  
A: 调用需要认证的接口（如获取购物车），如果返回 401 状态码，说明未登录或 token 无效。

---

## 10. 更新日志

- **v1.0.0** (2024-01-01)
  - 初始版本
  - 支持用户注册、登录
  - 支持商品浏览
  - 支持购物车管理（增删改查）
  - 支持订单创建和查询
  - 支持购物车增量更新（+1/-1）

---

**文档版本**: 1.0.0  
**最后更新**: 2024-01-01

