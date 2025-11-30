# 泡泡玛特 - 拉布布商城 API 接口文档

## 基础信息

- **Base URL**: `http://localhost:8080/api`
- **协议**: HTTP/HTTPS
- **数据格式**: JSON
- **字符编码**: UTF-8

## 认证说明

部分接口需要认证，需要在请求头中添加：

```
Authorization: Bearer {token}
```

登录成功后，服务器会返回 `token`，后续请求需要携带此 token。

---

## 1. 用户相关接口

### 1.1 用户注册

**接口地址**: `POST /api/register`

**接口描述**: 用户注册新账号

**请求头**:
```
Content-Type: application/json
```

**请求参数**:
```json
{
  "username": "string",  // 用户名，必填
  "password": "string",  // 密码，必填，最少6位
  "email": "string"      // 邮箱，必填，需符合邮箱格式
}
```

**请求示例**:
```bash
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "123456",
    "email": "test@example.com"
  }'
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
- `500`: 服务器内部错误

---

### 1.2 用户登录

**接口地址**: `POST /api/login`

**接口描述**: 用户登录获取 token

**请求头**:
```
Content-Type: application/json
```

**请求参数**:
```json
{
  "username": "string",  // 用户名，必填
  "password": "string"   // 密码，必填
}
```

**请求示例**:
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "123456"
  }'
```

**响应示例**:
```json
{
  "token": "user_1",
  "user": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "created_at": "2025-11-30T13:00:00Z",
    "updated_at": "2025-11-30T13:00:00Z"
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
- `500`: 服务器内部错误

---

## 2. 商品相关接口

### 2.1 获取商品列表

**接口地址**: `GET /api/products`

**接口描述**: 获取所有商品列表

**请求头**: 无需认证

**请求参数**: 无

**请求示例**:
```bash
curl -X GET http://localhost:8080/api/products
```

**响应示例**:
```json
{
  "products": [
    {
      "id": 1,
      "name": "拉布布 经典款",
      "description": "经典拉布布盲盒，随机款式",
      "price": 59.00,
      "image": "https://via.placeholder.com/300x300?text=拉布布经典款",
      "stock": 100,
      "series": "拉布布",
      "created_at": "2025-11-30T13:00:00Z",
      "updated_at": "2025-11-30T13:00:00Z"
    }
  ]
}
```

**状态码**:
- `200`: 成功
- `500`: 服务器内部错误

---

### 2.2 获取商品详情

**接口地址**: `GET /api/products/:id`

**接口描述**: 根据商品ID获取商品详情

**请求头**: 无需认证

**路径参数**:
- `id`: 商品ID（整数）

**请求示例**:
```bash
curl -X GET http://localhost:8080/api/products/1
```

**响应示例**:
```json
{
  "id": 1,
  "name": "拉布布 经典款",
  "description": "经典拉布布盲盒，随机款式",
  "price": 59.00,
  "image": "https://via.placeholder.com/300x300?text=拉布布经典款",
  "stock": 100,
  "series": "拉布布",
  "created_at": "2025-11-30T13:00:00Z",
  "updated_at": "2025-11-30T13:00:00Z"
}
```

**错误响应**:
```json
{
  "error": "商品不存在"
}
```

**状态码**:
- `200`: 成功
- `404`: 商品不存在
- `500`: 服务器内部错误

---

## 3. 购物车相关接口

> ⚠️ **注意**: 以下接口需要认证，请在请求头中添加 `Authorization: Bearer {token}`

### 3.1 获取购物车

**接口地址**: `GET /api/cart`

**接口描述**: 获取当前用户的购物车列表

**请求头**:
```
Authorization: Bearer {token}
```

**请求参数**: 无

**请求示例**:
```bash
curl -X GET http://localhost:8080/api/cart \
  -H "Authorization: Bearer user_1"
```

**响应示例**:
```json
{
  "items": [
    {
      "id": 1,
      "user_id": 1,
      "product_id": 1,
      "quantity": 2,
      "product": {
        "id": 1,
        "name": "拉布布 经典款",
        "description": "经典拉布布盲盒，随机款式",
        "price": 59.00,
        "image": "https://via.placeholder.com/300x300?text=拉布布经典款",
        "stock": 100,
        "series": "拉布布",
        "created_at": "2025-11-30T13:00:00Z",
        "updated_at": "2025-11-30T13:00:00Z"
      },
      "created_at": "2025-11-30T13:00:00Z",
      "updated_at": "2025-11-30T13:00:00Z"
    }
  ]
}
```

**错误响应**:
```json
{
  "error": "未授权"
}
```

**状态码**:
- `200`: 成功
- `401`: 未授权，需要登录
- `500`: 服务器内部错误

---

### 3.2 添加到购物车

**接口地址**: `POST /api/cart`

**接口描述**: 添加商品到购物车

**请求头**:
```
Content-Type: application/json
Authorization: Bearer {token}
```

**请求参数**:
```json
{
  "product_id": 1,  // 商品ID，必填
  "quantity": 1     // 数量，必填，最少1
}
```

**请求示例**:
```bash
curl -X POST http://localhost:8080/api/cart \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer user_1" \
  -d '{
    "product_id": 1,
    "quantity": 2
  }'
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
- `200`: 成功
- `400`: 请求参数错误或库存不足
- `401`: 未授权，需要登录
- `500`: 服务器内部错误

---

### 3.3 更新购物车商品数量

**接口地址**: `PUT /api/cart/:id`

**接口描述**: 更新购物车中商品的数量

**请求头**:
```
Content-Type: application/json
Authorization: Bearer {token}
```

**路径参数**:
- `id`: 购物车项ID（整数）

**请求参数**:
```json
{
  "quantity": 3  // 新的数量，必填，最少1
}
```

**请求示例**:
```bash
curl -X PUT http://localhost:8080/api/cart/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer user_1" \
  -d '{
    "quantity": 3
  }'
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
- `200`: 成功
- `400`: 请求参数错误、库存不足或购物车项不存在
- `401`: 未授权，需要登录
- `500`: 服务器内部错误

---

### 3.4 删除购物车商品

**接口地址**: `DELETE /api/cart/:id`

**接口描述**: 从购物车中删除商品

**请求头**:
```
Authorization: Bearer {token}
```

**路径参数**:
- `id`: 购物车项ID（整数）

**请求示例**:
```bash
curl -X DELETE http://localhost:8080/api/cart/1 \
  -H "Authorization: Bearer user_1"
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
- `200`: 成功
- `401`: 未授权，需要登录
- `404`: 购物车项不存在
- `500`: 服务器内部错误

---

## 4. 订单相关接口

> ⚠️ **注意**: 以下接口需要认证，请在请求头中添加 `Authorization: Bearer {token}`

### 4.1 创建订单

**接口地址**: `POST /api/orders`

**接口描述**: 根据购物车项创建订单

**请求头**:
```
Content-Type: application/json
Authorization: Bearer {token}
```

**请求参数**:
```json
{
  "cart_item_ids": [1, 2, 3]  // 购物车项ID数组，必填
}
```

**请求示例**:
```bash
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer user_1" \
  -d '{
    "cart_item_ids": [1, 2]
  }'
```

**响应示例**:
```json
{
  "message": "订单创建成功",
  "order_id": 1,
  "total_price": 177.00
}
```

**错误响应**:
```json
{
  "error": "商品库存不足: 拉布布 经典款"
}
```

**状态码**:
- `200`: 成功
- `400`: 请求参数错误或库存不足
- `401`: 未授权，需要登录
- `404`: 购物车项不存在
- `500`: 服务器内部错误

---

### 4.2 获取订单列表

**接口地址**: `GET /api/orders`

**接口描述**: 获取当前用户的所有订单

**请求头**:
```
Authorization: Bearer {token}
```

**请求参数**: 无

**请求示例**:
```bash
curl -X GET http://localhost:8080/api/orders \
  -H "Authorization: Bearer user_1"
```

**响应示例**:
```json
{
  "orders": [
    {
      "id": 1,
      "user_id": 1,
      "total_price": 177.00,
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
            "name": "拉布布 经典款",
            "description": "经典拉布布盲盒，随机款式",
            "price": 59.00,
            "image": "https://via.placeholder.com/300x300?text=拉布布经典款",
            "stock": 98,
            "series": "拉布布",
            "created_at": "2025-11-30T13:00:00Z",
            "updated_at": "2025-11-30T13:00:00Z"
          },
          "created_at": "2025-11-30T13:00:00Z",
          "updated_at": "2025-11-30T13:00:00Z"
        }
      ],
      "created_at": "2025-11-30T13:00:00Z",
      "updated_at": "2025-11-30T13:00:00Z"
    }
  ]
}
```

**错误响应**:
```json
{
  "error": "未授权"
}
```

**状态码**:
- `200`: 成功
- `401`: 未授权，需要登录
- `500`: 服务器内部错误

---

### 4.3 获取订单详情

**接口地址**: `GET /api/orders/:id`

**接口描述**: 根据订单ID获取订单详情

**请求头**:
```
Authorization: Bearer {token}
```

**路径参数**:
- `id`: 订单ID（整数）

**请求示例**:
```bash
curl -X GET http://localhost:8080/api/orders/1 \
  -H "Authorization: Bearer user_1"
```

**响应示例**:
```json
{
  "id": 1,
  "user_id": 1,
  "total_price": 177.00,
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
        "name": "拉布布 经典款",
        "description": "经典拉布布盲盒，随机款式",
        "price": 59.00,
        "image": "https://via.placeholder.com/300x300?text=拉布布经典款",
        "stock": 98,
        "series": "拉布布",
        "created_at": "2025-11-30T13:00:00Z",
        "updated_at": "2025-11-30T13:00:00Z"
      },
      "created_at": "2025-11-30T13:00:00Z",
      "updated_at": "2025-11-30T13:00:00Z"
    }
  ],
  "created_at": "2025-11-30T13:00:00Z",
  "updated_at": "2025-11-30T13:00:00Z"
}
```

**错误响应**:
```json
{
  "error": "订单不存在"
}
```

**状态码**:
- `200`: 成功
- `401`: 未授权，需要登录
- `404`: 订单不存在
- `500`: 服务器内部错误

---

## 5. 错误码说明

| 状态码 | 说明 |
|--------|------|
| 200 | 请求成功 |
| 400 | 请求参数错误 |
| 401 | 未授权，需要登录或token无效 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

## 6. 数据模型

### User (用户)
```json
{
  "id": 1,
  "username": "string",
  "email": "string",
  "created_at": "2025-11-30T13:00:00Z",
  "updated_at": "2025-11-30T13:00:00Z"
}
```

### Product (商品)
```json
{
  "id": 1,
  "name": "string",
  "description": "string",
  "price": 59.00,
  "image": "string",
  "stock": 100,
  "series": "拉布布",
  "created_at": "2025-11-30T13:00:00Z",
  "updated_at": "2025-11-30T13:00:00Z"
}
```

### CartItem (购物车项)
```json
{
  "id": 1,
  "user_id": 1,
  "product_id": 1,
  "quantity": 2,
  "product": { /* Product对象 */ },
  "created_at": "2025-11-30T13:00:00Z",
  "updated_at": "2025-11-30T13:00:00Z"
}
```

### Order (订单)
```json
{
  "id": 1,
  "user_id": 1,
  "total_price": 177.00,
  "status": "pending",
  "items": [ /* OrderItem数组 */ ],
  "created_at": "2025-11-30T13:00:00Z",
  "updated_at": "2025-11-30T13:00:00Z"
}
```

### OrderItem (订单项)
```json
{
  "id": 1,
  "order_id": 1,
  "product_id": 1,
  "quantity": 2,
  "price": 59.00,
  "product": { /* Product对象 */ },
  "created_at": "2025-11-30T13:00:00Z",
  "updated_at": "2025-11-30T13:00:00Z"
}
```

## 7. 注意事项

1. **Token 格式**: 当前实现为简化版，token 格式为 `user_{id}`，生产环境建议使用 JWT
2. **时间格式**: 所有时间字段使用 ISO 8601 格式（UTC 时间）
3. **价格精度**: 价格使用 `decimal(10,2)` 类型，保留两位小数
4. **库存检查**: 添加购物车和创建订单时会检查库存，库存不足会返回错误
5. **事务处理**: 创建订单使用数据库事务，确保数据一致性

## 8. 测试建议

1. 使用 Postman 或类似工具测试 API
2. 先注册用户，然后登录获取 token
3. 浏览商品列表，选择商品加入购物车
4. 查看购物车，修改商品数量
5. 创建订单，查看订单历史

## 9. 更新日志

- **v1.0.0** (2025-11-30)
  - 初始版本
  - 实现用户注册、登录
  - 实现商品展示
  - 实现购物车管理
  - 实现订单管理

