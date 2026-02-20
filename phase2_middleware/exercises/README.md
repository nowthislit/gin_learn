# Phase 2 练习：中间件与验证系统

## 练习目标
通过实现一个带认证和限流的 API 系统，巩固中间件、验证和错误处理的知识点。

## 功能需求

### 基础功能（必做）

#### 1. JWT 认证中间件
实现基于 JWT 的认证系统：
- **登录接口** - POST /api/v1/auth/login
  - 验证用户名密码，返回 JWT Token
- **注册接口** - POST /api/v1/auth/register
  - 创建新用户
- **刷新 Token** - POST /api/v1/auth/refresh
  - 使用 Refresh Token 获取新的 Access Token
- **JWT 中间件**
  - 验证请求头中的 Authorization: Bearer <token>
  - 无效 Token 返回 401
  - 将用户信息（user_id, username）存入 Context

#### 2. 请求限流中间件
实现基于 IP 的请求限流：
- 每分钟最多 60 次请求
- 超出限制返回 429 (Too Many Requests)
- 在响应头中返回剩余请求次数

#### 3. 日志中间件
实现详细的请求日志：
- 请求方法、路径、IP
- 请求耗时
- 响应状态码
- 用户 ID（如果已登录）

#### 4. 统一错误处理
实现全局错误处理中间件：
- 捕获 panic，防止程序崩溃
- 统一错误响应格式
- 记录错误日志

### 进阶功能（选做）

#### 5. RBAC 权限控制
实现基于角色的权限控制：
- 角色：admin, user, guest
- 资源权限控制（如只有 admin 可以删除用户）

#### 6. 请求验证中间件
实现通用的请求验证：
- 统一验证请求参数
- 支持自定义验证规则

## API 列表

### 公开 API（无需认证）
- `POST /api/v1/auth/register` - 注册
- `POST /api/v1/auth/login` - 登录
- `GET /api/v1/public/info` - 公开信息

### 受保护 API（需要认证）
- `GET /api/v1/user/profile` - 获取当前用户信息
- `PUT /api/v1/user/profile` - 更新用户信息

### 管理员 API（需要 admin 角色）
- `GET /api/v1/admin/users` - 获取所有用户列表
- `DELETE /api/v1/admin/users/:id` - 删除用户

## 数据存储
使用内存 map 存储用户数据。

## 用户结构
```go
type User struct {
    ID        uint      `json:"id"`
    Username  string    `json:"username"`
    Password  string    `json:"-"`
    Email     string    `json:"email"`
    Role      string    `json:"role"` // admin, user, guest
    CreatedAt time.Time `json:"created_at"`
}
```

## JWT 结构
```go
type Claims struct {
    UserID   uint   `json:"user_id"`
    Username string `json:"username"`
    Role     string `json:"role"`
    jwt.StandardClaims
}
```

## 响应格式
```json
{
    "code": 200,
    "message": "success",
    "data": { }
}
```

## 错误码定义
- 200 - 成功
- 400 - 参数错误
- 401 - 未授权
- 403 - 无权限
- 404 - 资源不存在
- 429 - 请求过多
- 500 - 服务器错误

## 测试示例

```bash
# 1. 注册用户
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"123456","email":"admin@example.com","role":"admin"}'

# 2. 登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"123456"}'
# 返回：{"data":{"access_token":"xxx","refresh_token":"yyy"}}

# 3. 访问受保护接口（带 Token）
curl http://localhost:8080/api/v1/user/profile \
  -H "Authorization: Bearer <access_token>"

# 4. 访问管理员接口
curl http://localhost:8080/api/v1/admin/users \
  -H "Authorization: Bearer <access_token>"

# 5. 测试限流（快速请求）
for i in {1..70}; do 
  curl -s http://localhost:8080/api/v1/public/info | head -1
done
```

## 评分标准

| 功能 | 分值 |
|------|------|
| JWT 生成和验证 | 25分 |
| 认证中间件 | 20分 |
| 限流中间件 | 20分 |
| 日志中间件 | 15分 |
| 错误处理 | 10分 |
| RBAC 权限控制 | 10分 |
| **总分** | **100分** |

## 依赖包
```bash
go get github.com/golang-jwt/jwt/v5
```

## 提示

### 1. JWT 生成示例
```go
func GenerateToken(userID uint, username, role string) (string, error) {
    claims := Claims{
        UserID:   userID,
        Username: username,
        Role:     role,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Add(time.Hour * 2).Unix(),
            IssuedAt:  time.Now().Unix(),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte("your-secret-key"))
}
```

### 2. 限流器实现思路
```go
type RateLimiter struct {
    requests map[string][]time.Time // IP -> 请求时间列表
    mu       sync.RWMutex
    limit    int           // 限制次数
    window   time.Duration // 时间窗口
}
```

### 3. 中间件执行顺序
```go
r.Use(Recovery())      // 1. 恢复 panic
r.Use(Logger())        // 2. 记录日志
r.Use(RateLimit())     // 3. 限流检查
r.Use(JWTAuth())       // 4. 认证检查
```

## 参考答案位置
完成练习后，可以查看 `solution/main.go` 参考实现。

Good luck! 🔒
