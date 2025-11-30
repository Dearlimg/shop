# 泡泡玛特 - 拉布布电商网站

基于 Hertz 框架的泡泡玛特拉布布系列商品销售网站。

## 功能特性

- ✅ **商品展示**: 展示拉布布系列商品
- ✅ **用户系统**: 用户注册、登录
- ✅ **购物车**: 添加、删除、修改购物车商品
- ✅ **订单管理**: 创建订单、查看订单历史

## 技术栈

- **后端**: Hertz (Go 语言) - RESTful API
- **ORM**: GORM - Go 语言的 ORM 框架
- **前端**: HTML + JavaScript
- **数据库**: MySQL

## 项目结构

```
shop/
├── main.go              # 主程序入口
├── config.yaml          # 配置文件
├── config/              # 配置相关
│   └── config.go
├── model/               # 数据模型
│   ├── user.go
│   ├── product.go
│   ├── cart.go
│   └── order.go
├── dao/                 # 数据访问层
│   ├── user_dao.go
│   ├── product_dao.go
│   ├── cart_dao.go
│   └── order_dao.go
├── logic/               # 业务逻辑层
│   ├── user_logic.go
│   ├── product_logic.go
│   ├── cart_logic.go
│   └── order_logic.go
├── controller/         # 控制器层
│   └── api/
│       ├── user_controller.go
│       ├── product_controller.go
│       ├── cart_controller.go
│       └── order_controller.go
├── routers/            # 路由
│   └── router.go
├── middleware/         # 中间件
│   └── auth.go
├── global/             # 全局变量
│   └── db/
│       ├── db.go
│       └── migrations.go
└── static/            # 前端静态文件
    └── index.html
```

## 环境要求

- Go 1.21+
- MySQL 5.7+

## 安装和运行

### 1. 安装依赖

```bash
go mod tidy
```

### 2. 配置数据库

编辑 `config.yaml` 文件，配置数据库连接信息：

```yaml
database:
  host: 47.118.19.28      # 数据库地址
  port: 3306               # 数据库端口
  user: root               # 数据库用户名
  password: sta_go         # 数据库密码
  database: blog           # 数据库名称
  charset: utf8mb4         # 字符集

server:
  port: 8080               # 服务器端口
  host: "0.0.0.0"          # 服务器地址
```

**注意**: 请确保数据库已创建，程序会自动创建所需的表结构。

### 3. 运行程序

默认使用项目根目录下的 `config.yaml` 文件：

```bash
go run main.go
```

或者指定配置文件路径（通过环境变量）：

```bash
export CONFIG_PATH="/path/to/config.yaml"
go run main.go
```

服务器将根据配置文件中的设置启动（默认 `http://0.0.0.0:8080`）。

### 4. 访问网站

在浏览器中打开 `http://localhost:8080`

## API 接口

详细的 API 接口文档请查看：[API 接口文档](./docs/API.md)

### 快速参考

**公开接口**:
- `POST /api/register` - 用户注册
- `POST /api/login` - 用户登录
- `GET /api/products` - 获取商品列表
- `GET /api/products/:id` - 获取商品详情

**需要认证的接口**（需在 Header 中添加 `Authorization: Bearer {token}`）:
- `GET /api/cart` - 获取购物车
- `POST /api/cart` - 添加到购物车
- `PUT /api/cart/:id` - 更新购物车商品数量
- `DELETE /api/cart/:id` - 删除购物车商品
- `POST /api/orders` - 创建订单
- `GET /api/orders` - 获取订单列表
- `GET /api/orders/:id` - 获取订单详情

## 使用说明

1. **注册/登录**: 点击右上角的"注册"或"登录"按钮
2. **浏览商品**: 在商品页面浏览拉布布系列商品
3. **添加到购物车**: 登录后可以点击"加入购物车"按钮
4. **管理购物车**: 在购物车标签页可以修改数量或删除商品
5. **下单**: 在购物车页面点击"结算"按钮创建订单
6. **查看订单**: 在订单标签页查看订单历史

## 配置文件说明

项目使用 `config.yaml` 文件进行配置管理。配置文件包含以下部分：

- **database**: 数据库连接配置
  - `host`: 数据库服务器地址
  - `port`: 数据库端口（默认 3306）
  - `user`: 数据库用户名
  - `password`: 数据库密码
  - `database`: 数据库名称
  - `charset`: 字符集（默认 utf8mb4）

- **server**: 服务器配置
  - `host`: 服务器监听地址（默认 0.0.0.0）
  - `port`: 服务器端口（默认 8080）

可以通过环境变量 `CONFIG_PATH` 指定配置文件路径：
```bash
export CONFIG_PATH="/path/to/your/config.yaml"
go run main.go
```

## 注意事项

- 首次运行会自动创建数据库表和初始化商品数据
- Token 认证为简化实现，生产环境建议使用 JWT
- 前端使用 localStorage 存储 token，刷新页面后仍保持登录状态
- **重要**: 请妥善保管 `config.yaml` 文件，不要将包含敏感信息的配置文件提交到版本控制系统
