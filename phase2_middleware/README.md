# Phase 2: 中间件与验证

## 学习目标
- 理解中间件的概念和执行流程
- 学会使用内置中间件（Logger、Recovery）
- 掌握自定义中间件的开发
- 学会在中间件间传递数据
- 掌握请求数据的验证
- 了解自定义验证规则
- 学会全局错误处理

## 预计学习时间
5-7天

## 目录结构
```
phase2_middleware/
├── go.mod       # Go模块配置
├── go.sum       # 依赖锁定
├── main.go      # 完整示例代码
└── README.md    # 本文件
```

## 运行方式

```bash
cd phase2_middleware
go run main.go
```

## 知识点详解

### 1. 中间件基础

中间件是处理HTTP请求的函数，在请求到达处理器之前或之后执行。

```go
// 中间件签名
func MiddlewareName() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 请求前处理
        
        c.Next() // 继续执行后续中间件和处理器
        
        // 请求后处理（响应返回前）
    }
}
```

**执行流程：**
```
请求 → 中间件1（前） → 中间件2（前） → 处理器 → 中间件2（后） → 中间件1（后） → 响应
```

### 2. 使用中间件的三种方式

#### 全局中间件
```go
// 应用于所有路由
r := gin.New()
r.Use(LoggerMiddleware())
r.Use(CORSMiddleware())
```

#### 路由组中间件
```go
// 仅应用于特定分组
authorized := r.Group("/admin")
authorized.Use(AuthMiddleware())
{
    authorized.GET("/users", getUsers)
}
```

#### 单个路由中间件
```go
// 仅应用于单个路由
r.GET("/admin", AdminOnly(), adminHandler)
```

### 3. 内置中间件

```go
// Recovery - 捕获panic，防止程序崩溃
r.Use(gin.Recovery())

// Logger - 日志记录
r.Use(gin.Logger())

// LoggerWithFormatter - 自定义格式
r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
    return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
        param.ClientIP,
        param.TimeStamp.Format(time.RFC1123),
        param.Method,
        param.Path,
        param.Request.Proto,
        param.StatusCode,
        param.Latency,
        param.Request.UserAgent(),
        param.ErrorMessage,
    )
}))
```

### 4. 常用自定义中间件

#### 日志中间件
```go
func LoggerMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        
        // 处理请求
        c.Next()
        
        // 记录日志
        latency := time.Since(start)
        status := c.Writer.Status()
        
        log.Printf("[%d] %s %s %v", status, c.Request.Method, path, latency)
    }
}
```

#### CORS中间件
```go
func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        
        c.Next()
    }
}
```

#### 认证中间件
```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        
        if token == "" {
            c.JSON(401, gin.H{"error": "未授权"})
            c.Abort() // 终止后续处理
            return
        }
        
        // 验证token...
        
        c.Next()
    }
}
```

### 5. 中间件数据传递

```go
func SetUserMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 设置数据
        c.Set("user_id", "12345")
        c.Set("username", "admin")
        
        c.Next()
    }
}

func Handler(c *gin.Context) {
    // 获取数据
    userID, exists := c.Get("user_id")
    if !exists {
        // 数据不存在
    }
}
```

### 6. 数据验证

#### 基本验证标签
```go
type User struct {
    Username string `binding:"required,min=3,max=20"`  // 必填，长度3-20
    Password string `binding:"required,min=6"`          // 必填，最小6位
    Email    string `binding:"required,email"`          // 必填，邮箱格式
    Age      int    `binding:"gte=0,lte=150"`           // 范围0-150
    Website  string `binding:"url"`                     // URL格式
    Phone    string `binding:"e164"`                    // 国际电话格式
}
```

#### 常用验证标签
- `required` - 必填
- `min`, `max` - 最小/最大长度或值
- `len` - 固定长度
- `email` - 邮箱格式
- `url` - URL格式
- `uuid` - UUID格式
- `oneof=a b c` - 枚举值
- `gt`, `gte`, `lt`, `lte` - 比较大小
- `numeric`, `alpha`, `alphanumeric` - 字符类型

#### 验证错误处理
```go
func handler(c *gin.Context) {
    var req User
    if convErr := c.ShouldBindJSON(&req); convErr != nil {
        // 解析验证错误
        if errs, ok := convErr.(validator.ValidationErrors); ok {
            for _, e := range errs {
                fmt.Printf("%s: %s\n", e.Field(), e.Tag())
            }
        }
        
        c.JSON(400, gin.H{"error": convErr.Error()})
        return
    }
}
```

### 7. 自定义验证规则

```go
import "github.com/go-playground/validator/v10"

// 自定义验证函数
func customValidator(fl validator.FieldLevel) bool {
    value := fl.Field().String()
    // 验证逻辑...
    return true
}

// 注册验证器
func init() {
    if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
        v.RegisterValidation("custom", customValidator)
    }
}

// 使用自定义验证
type Request struct {
    Code string `json:"code" binding:"required,custom"`
}
```

### 8. 错误处理

#### 全局错误恢复
```go
r.Use(gin.CustomRecovery(func(c *gin.Context, convErr interface{}) {
    log.Println("Panic recovered:", convErr)
    c.JSON(500, gin.H{
        "error": "服务器内部错误",
    })
}))
```

#### 业务错误处理
```go
func handler(c *gin.Context) {
    if convErr := doSomething(); convErr != nil {
        c.JSON(500, gin.H{
            "error": "操作失败",
            "details": convErr.Error(),
        })
        return
    }
    
    c.JSON(200, gin.H{"message": "成功"})
}
```

## 实战练习

### 练习1：完整的认证系统

实现一个包含以下功能的认证中间件：
1. JWT Token验证
2. Token刷新机制
3. 权限检查（角色：user/admin/superadmin）
4. 黑名单（登出用户）

### 练习2：请求限流中间件

实现基于IP的请求限流中间件：
1. 每IP每分钟最多100次请求
2. 超出限制返回429状态码
3. 支持动态调整限制值
4. 记录限流日志

### 练习3：完整的数据验证

为电商系统设计数据验证：
1. 用户注册（用户名、密码、邮箱、手机号）
2. 商品创建（名称、价格、库存、分类）
3. 订单创建（商品列表、地址、支付方式）
4. 自定义验证规则（手机号格式、商品价格范围等）

### 练习4：链路追踪中间件

实现分布式链路追踪：
1. 生成/传播Trace ID
2. 记录请求耗时
3. 记录关键节点时间
4. 集成日志系统

## 测试命令

```bash
# 测试公开接口（无需认证）
curl http://localhost:8080/public/info

# 测试受保护接口（无token）
curl http://localhost:8080/private/profile
# 预期返回401

# 测试受保护接口（有token）
curl http://localhost:8080/private/profile \
  -H "Authorization: Bearer valid-token"

# 测试管理员接口（无权限）
curl http://localhost:8080/admin
# 预期返回403

# 测试管理员接口（有权限）
curl http://localhost:8080/admin \
  -H "X-User-Role: admin"

# 测试注册验证（成功）
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"username":"john","password":"123456","email":"john@example.com","age":25}'

# 测试注册验证（失败-参数错误）
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"username":"ab","password":"123","email":"invalid","age":200}'

# 测试panic恢复
curl http://localhost:8080/panic
# 程序不会崩溃，返回500错误
```

## 常见问题

### Q1: 中间件执行顺序是怎样的？
```go
r.Use(M1, M2, M3)
// 执行顺序: M1(前) → M2(前) → M3(前) → Handler → M3(后) → M2(后) → M1(后)
```

### Q2: 如何在中间件中提前返回？
```go
func middleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        if !check() {
            c.JSON(401, gin.H{"error": "未授权"})
            c.Abort() // 终止后续执行
            return
        }
        c.Next()
    }
}
```

### Q3: 验证错误信息如何自定义？
```go
// 注册翻译器
import "github.com/go-playground/validator/v10/translations/zh"

// 翻译验证错误
for _, e := range errs {
    translatedErr := e.Translate(trans)
    fmt.Println(translatedErr)
}
```

## 下一步

完成第二阶段后，进入**Phase 3: 数据库集成**，学习：
- GORM基础使用
- 数据库连接池配置
- CRUD操作封装
- 事务处理
- 关联关系处理
