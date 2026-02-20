# Phase 4 练习：电商系统

## 练习目标
通过实现一个完整的电商系统，掌握企业级项目开发的最佳实践。

## 系统功能

### 基础功能（必做）

#### 1. 用户系统
- 用户注册/登录（JWT认证）
- 用户信息管理
- 收货地址管理

#### 2. 商品系统
- 商品分类（树形结构）
- 商品管理（CRUD）
- 商品搜索（Elasticsearch 或简单搜索）
- 商品库存管理

#### 3. 购物车系统
- 添加商品到购物车
- 修改购物车商品数量
- 删除购物车商品
- 清空购物车
- 购物车结算

#### 4. 订单系统
- 创建订单（从事务）
- 订单状态流转（待支付->已支付->已发货->已完成）
- 订单取消（从事务，回滚库存）
- 订单列表和详情

#### 5. 支付系统（模拟）
- 模拟支付接口
- 支付回调处理

### 进阶功能（选做）

#### 6. 优惠券系统
- 优惠券类型（满减、折扣、直减）
- 优惠券领取和使用
- 优惠券过期处理

#### 7. 秒杀系统
- 秒杀活动管理
- 库存预热（Redis）
- 秒杀队列（防止超卖）
- 限流控制

#### 8. 消息通知
- 订单状态变更通知
- 短信/邮件通知（模拟）

#### 9. 数据统计
- 销售统计
- 商品热度排行
- 用户行为分析

## 项目架构

采用分层架构：

```
phase4_advanced/exercises/
├── main.go                  # 入口
├── config.yaml             # 配置文件
├── config/
│   └── config.go          # 配置管理
├── internal/
│   ├── api/               # API层
│   │   ├── server.go
│   │   ├── handler/       # 处理器
│   │   └── middleware/    # 中间件
│   ├── service/           # 业务层
│   ├── repository/        # 数据层
│   ├── model/             # 模型
│   └── utils/             # 工具函数
├── pkg/
│   ├── logger/            # 日志
│   ├── cache/             # 缓存
│   └── queue/             # 队列
└── README.md
```

## 数据模型

### User（用户）
```go
type User struct {
    ID        uint      `gorm:"primarykey"`
    Username  string    `gorm:"uniqueIndex;not null"`
    Email     string    `gorm:"uniqueIndex;not null"`
    Password  string    `gorm:"not null"`
    Phone     string
    Status    int       // 1:正常 0:禁用
    CreatedAt time.Time
}
```

### Product（商品）
```go
type Product struct {
    ID          uint      `gorm:"primarykey"`
    Name        string    `gorm:"not null;index"`
    Description string    `gorm:"type:text"`
    Price       float64   `gorm:"not null"`
    Stock       int       `gorm:"not null"`
    CategoryID  uint
    Status      int       // 1:上架 0:下架
    CreatedAt   time.Time
}
```

### Cart（购物车）
```go
type Cart struct {
    ID        uint    `gorm:"primarykey"`
    UserID    uint    `gorm:"index"`
    ProductID uint
    Quantity  int
    Selected  bool    // 是否选中结算
}
```

### Order（订单）
```go
type Order struct {
    ID            uint      `gorm:"primarykey"`
    OrderNo       string    `gorm:"uniqueIndex;not null"` // 订单号
    UserID        uint
    TotalAmount   float64   // 订单总金额
    PayAmount     float64   // 实付金额
    Status        string    // pending, paid, shipped, completed, cancelled
    AddressID     uint      // 收货地址
    Remark        string    // 订单备注
    PaidAt        *time.Time
    ShippedAt     *time.Time
    CompletedAt   *time.Time
    CreatedAt     time.Time
    
    Items []OrderItem // 订单商品
}
```

### OrderItem（订单商品）
```go
type OrderItem struct {
    ID          uint    `gorm:"primarykey"`
    OrderID     uint
    ProductID   uint
    ProductName string  // 快照：商品名称
    Price       float64 // 快照：商品价格
    Quantity    int
    TotalPrice  float64
}
```

## API 列表

### 用户模块
- `POST /api/v1/auth/register` - 注册
- `POST /api/v1/auth/login` - 登录
- `GET /api/v1/user/profile` - 个人信息
- `PUT /api/v1/user/profile` - 修改信息

### 地址模块
- `GET /api/v1/addresses` - 地址列表
- `POST /api/v1/addresses` - 添加地址
- `PUT /api/v1/addresses/:id` - 修改地址
- `DELETE /api/v1/addresses/:id` - 删除地址

### 商品模块
- `GET /api/v1/products` - 商品列表（分页、分类、搜索）
- `GET /api/v1/products/:id` - 商品详情
- `GET /api/v1/categories` - 分类列表（树形）
- `GET /api/v1/products/search` - 商品搜索

### 购物车模块
- `GET /api/v1/cart` - 购物车列表
- `POST /api/v1/cart` - 添加商品
- `PUT /api/v1/cart/:id` - 修改数量
- `DELETE /api/v1/cart/:id` - 删除商品
- `DELETE /api/v1/cart` - 清空购物车
- `POST /api/v1/cart/checkout` - 结算

### 订单模块
- `POST /api/v1/orders` - 创建订单
- `GET /api/v1/orders` - 订单列表
- `GET /api/v1/orders/:id` - 订单详情
- `POST /api/v1/orders/:id/cancel` - 取消订单
- `POST /api/v1/orders/:id/pay` - 支付订单（模拟）

## 练习要求

### 1. 项目结构
- 严格按照分层架构组织代码
- 每层通过接口进行依赖注入
- 配置文件使用 Viper 管理

### 2. 日志系统
- 使用 Zap 记录日志
- 区分不同级别的日志
- 记录请求追踪 ID

### 3. 缓存系统
- 使用 Redis 缓存热点数据
- 商品详情缓存
- 购物车缓存

### 4. 事务管理
- 订单创建使用事务
- 支付回调使用事务
- 库存扣减使用乐观锁

### 5. 错误处理
- 统一的错误码定义
- 全局错误处理中间件
- 错误日志记录

### 6. 接口安全
- JWT 认证
- 接口限流
- 参数验证

## 测试示例

```bash
# 1. 注册
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"buyer","email":"buyer@example.com","password":"123456"}'

# 2. 登录获取 Token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"buyer","password":"123456"}'

# 3. 添加收货地址
curl -X POST http://localhost:8080/api/v1/addresses \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "name":"张三",
    "phone":"13800138000",
    "province":"北京市",
    "city":"北京市",
    "district":"朝阳区",
    "detail":"xxx街道xxx号"
  }'

# 4. 添加商品到购物车
curl -X POST http://localhost:8080/api/v1/cart \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"product_id":1,"quantity":2}'

# 5. 查看购物车
curl http://localhost:8080/api/v1/cart \
  -H "Authorization: Bearer <token>"

# 6. 创建订单
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "cart_item_ids":[1,2],
    "address_id":1,
    "remark":"请尽快发货"
  }'

# 7. 支付订单（模拟）
curl -X POST http://localhost:8080/api/v1/orders/1/pay \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"payment_method":"alipay"}'
```

## 评分标准

| 功能 | 分值 |
|------|------|
| 项目架构 | 15分 |
| 用户系统 | 10分 |
| 商品系统 | 15分 |
| 购物车系统 | 15分 |
| 订单系统（含事务） | 20分 |
| 缓存使用 | 10分 |
| 日志系统 | 10分 |
| 代码规范 | 5分 |
| **总分** | **100分** |

## 技术栈

- **Web 框架**: Gin
- **数据库**: MySQL / PostgreSQL / SQLite
- **ORM**: GORM
- **缓存**: Redis
- **配置**: Viper
- **日志**: Zap
- **验证**: go-playground/validator

## 提示

### 1. 订单号生成
```go
func GenerateOrderNo() string {
    // 格式：年月日时分秒 + 6位随机数
    return time.Now().Format("20060102150405") + randomString(6)
}
```

### 2. 库存扣减（乐观锁）
```go
// 使用 version 字段实现乐观锁
result := db.Model(&Product{}).
    Where("id = ? AND stock >= ?", productID, quantity).
    Updates(map[string]interface{}{
        "stock": gorm.Expr("stock - ?", quantity),
    })

if result.RowsAffected == 0 {
    return errors.New("库存不足")
}
```

### 3. Redis 缓存商品
```go
// 获取商品时先查缓存
func (s *ProductService) GetProduct(id uint) (*Product, error) {
    // 1. 查 Redis
    key := fmt.Sprintf("product:%d", id)
    data, err := redis.Get(key)
    if err == nil {
        var product Product
        json.Unmarshal(data, &product)
        return &product, nil
    }
    
    // 2. 查数据库
    product, err := s.repo.GetByID(id)
    if err != nil {
        return nil, err
    }
    
    // 3. 写入 Redis
    data, _ = json.Marshal(product)
    redis.Set(key, data, time.Hour)
    
    return product, nil
}
```

Good luck! 🛒
