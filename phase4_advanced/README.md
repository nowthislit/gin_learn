# Phase 4: 进阶实战

## 学习目标
- 掌握项目分层架构设计
- 学会使用配置管理工具（Viper）
- 学会使用日志系统（Zap）
- 理解依赖注入和接口设计
- 掌握完整的RESTful API开发
- 了解生产级项目最佳实践

## 预计学习时间
7-10天

## 项目架构

本项目采用经典的分层架构（Layered Architecture）：

```
phase4_advanced/
├── main.go              # 应用程序入口
├── config.yaml          # 配置文件
├── config/              # 配置管理
│   └── config.go       # Viper配置加载
├── internal/            # 内部代码（不对外暴露）
│   ├── api/            # API层（Handler/Controller）
│   │   ├── server.go   # HTTP服务器
│   │   ├── middleware.go # 中间件
│   │   ├── user_handler.go
│   │   ├── product_handler.go
│   │   ├── category_handler.go
│   │   └── order_handler.go
│   ├── service/        # 业务逻辑层
│   │   ├── service.go
│   │   ├── user_service.go
│   │   ├── product_service.go
│   │   ├── category_service.go
│   │   └── order_service.go
│   ├── repository/     # 数据访问层（DAO）
│   │   ├── repository.go
│   │   ├── user_repository.go
│   │   ├── product_repository.go
│   │   ├── category_repository.go
│   │   └── order_repository.go
│   ├── model/          # 数据模型（Entity）
│   │   └── model.go
│   └── middleware/     # 自定义中间件
├── pkg/                 # 公共包（可对外使用）
│   └── logger/         # 日志工具
└── README.md           # 本文件
```

### 分层职责

| 层级 | 职责 | 示例 |
|------|------|------|
| API层 | 接收请求，参数校验，调用Service，返回响应 | Handler |
| Service层 | 业务逻辑处理，事务控制，编排Repository | Service |
| Repository层 | 数据库操作，数据持久化 | Repository |
| Model层 | 数据结构定义，数据验证 | Model |

### 数据流向

```
HTTP Request → API(Handler) → Service → Repository → Database
     ↑              ↑            ↑           ↑            ↑
     |              |            |           |            |
   响应返回      参数校验     业务逻辑    数据查询     数据存储
```

## 运行方式

```bash
cd phase4_advanced

# 方式1: 直接运行
go run main.go

# 方式2: 编译后运行
go build -o app.exe
./app.exe

# 方式3: 指定配置文件
go run main.go -config=/path/to/config.yaml
```

服务启动后会监听 8080 端口。

## 核心知识点

### 1. 配置管理（Viper）

Viper 是 Go 语言中流行的配置管理库，支持多种配置格式。

```go
// 定义配置结构
type Config struct {
    App    AppConfig    `mapstructure:"app"`
    Server ServerConfig `mapstructure:"server"`
    DB     DBConfig     `mapstructure:"database"`
}

// 加载配置
viper.SetConfigName("config")
viper.SetConfigType("yaml")
viper.AddConfigPath(".")
viper.ReadInConfig()
viper.Unmarshal(&config)
```

**特性：**
- 支持 JSON/TOML/YAML/ENV 等多种格式
- 支持默认值设置
- 支持环境变量覆盖
- 支持配置热重载

### 2. 日志系统（Zap）

Zap 是 Uber 开源的高性能日志库。

```go
// 初始化
logger.Init("info")
defer logger.Sync()

// 使用
logger.Info("User created", 
    logger.String("username", "john"),
    logger.Int("id", 123))

logger.Error("Failed to create user", 
    logger.ErrorField(err))
```

**特性：**
- 高性能（结构化日志比标准库快几倍）
- 结构化日志（JSON格式）
- 日志级别控制（Debug/Info/Warn/Error/Fatal）
- 字段类型安全

### 3. 依赖注入

通过接口实现依赖注入，便于测试和替换实现。

```go
// 定义接口
type UserRepository interface {
    Create(user *model.User) error
    GetByID(id uint) (*model.User, error)
}

// 实现接口
type userRepository struct {
    db *gorm.DB
}

// 构造函数注入
type UserService struct {
    repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
    return &UserService{repo: repo}
}
```

### 4. 接口设计原则

**Repository 接口:**
```go
type UserRepository interface {
    Create(user *model.User) error
    GetByID(id uint) (*model.User, error)
    List(page, pageSize int, keyword string) ([]model.User, int64, error)
    Update(user *model.User) error
    Delete(id uint) error
}
```

**Service 接口:**
```go
type UserService interface {
    Register(username, email, password string, age int) (*model.User, error)
    GetUser(id uint) (*model.User, error)
    ListUsers(page, pageSize int, keyword string) ([]model.User, int64, error)
    UpdateUser(id uint, updates map[string]interface{}) error
    DeleteUser(id uint) error
}
```

## API列表

### 用户管理
- `GET    /api/v1/users` - 获取用户列表
- `GET    /api/v1/users/:id` - 获取用户详情
- `POST   /api/v1/users` - 创建用户
- `PUT    /api/v1/users/:id` - 更新用户
- `DELETE /api/v1/users/:id` - 删除用户

### 产品管理
- `GET    /api/v1/products` - 获取产品列表
- `GET    /api/v1/products/:id` - 获取产品详情
- `POST   /api/v1/products` - 创建产品
- `PUT    /api/v1/products/:id` - 更新产品
- `DELETE /api/v1/products/:id` - 删除产品

### 分类管理
- `GET    /api/v1/categories` - 获取分类列表
- `GET    /api/v1/categories/:id` - 获取分类详情
- `POST   /api/v1/categories` - 创建分类

### 订单管理
- `GET    /api/v1/orders` - 获取订单列表
- `GET    /api/v1/orders/:id` - 获取订单详情
- `POST   /api/v1/orders` - 创建订单（含事务）
- `POST   /api/v1/orders/:id/cancel` - 取消订单（含事务）

### 搜索
- `GET    /api/v1/search/products` - 高级搜索产品

## 测试命令

```bash
# 测试用户API
curl http://localhost:8080/api/v1/users
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@test.com","password":"123456","age":25}'

# 测试产品API
curl http://localhost:8080/api/v1/products
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{"name":"iPhone 16","price":7999,"stock":100,"category_id":1}'

# 测试搜索
curl "http://localhost:8080/api/v1/search/products?keyword=iPhone&min_price=1000&sort_by=price&sort_order=desc"

# 测试订单（事务）
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "items": [
      {"product_id": 1, "quantity": 2},
      {"product_id": 2, "quantity": 1}
    ]
  }'
```

## 生产部署建议

### 1. 配置文件管理
- 开发环境：`config.yaml`
- 测试环境：`config.test.yaml`
- 生产环境：使用环境变量或配置中心（Consul/ETCD）

### 2. 日志配置
- 开发：输出到控制台
- 生产：输出到文件，使用日志轮转
- 考虑集成 ELK（Elasticsearch + Logstash + Kibana）

### 3. 数据库配置
- 使用连接池
- 配置合理的最大连接数
- 使用读写分离（如果数据量大）

### 4. 监控和告警
- 集成 Prometheus 监控
- 配置关键指标告警
- 使用 Jaeger 做链路追踪

### 5. 容器化部署

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY config.yaml .
EXPOSE 8080
CMD ["./main"]
```

## 扩展方向

完成本阶段后，可以进一步学习：

1. **认证授权**
   - JWT认证
   - OAuth2.0
   - RBAC权限控制

2. **缓存**
   - Redis缓存
   - 本地缓存（BigCache）
   - 缓存一致性策略

3. **微服务**
   - gRPC通信
   - 服务发现（Consul/Nacos）
   - 熔断限流（Hystrix/Rate Limit）

4. **消息队列**
   - Kafka
   - RabbitMQ
   - NATS

5. **性能优化**
   - pprof性能分析
   - 压测和调优
   - 数据库索引优化

## 常见问题

### Q1: 为什么要分层？
分层可以让代码：
- 职责清晰，易于维护
- 便于单元测试
- 支持灵活替换实现
- 团队协作更高效

### Q2: 接口应该定义在哪一层？
- Repository 接口：由 Service 层定义，Repository 层实现
- Service 接口：由 API 层定义，Service 层实现
- 遵循依赖倒置原则（DIP）

### Q3: 如何处理数据库事务？
```go
// 在Repository层处理事务
func (r *orderRepository) CreateOrder(...) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        // 在事务中执行操作
        // ...
        return nil // 提交
        // return err // 回滚
    })
}
```

### Q4: 如何实现优雅的关机？
```go
// 监听系统信号
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit

// 优雅关闭
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
server.Shutdown(ctx)
```

## 下一步

恭喜！你已经完成了 Gin 框架的所有学习阶段。

建议你接下来：
1. 做一个完整项目实践（如博客系统、电商平台）
2. 学习 Go 语言的高级特性
3. 研究优秀开源项目的架构设计
4. 参与开源项目贡献

Good luck! 🎉
