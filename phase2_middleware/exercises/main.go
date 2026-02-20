package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// ==================== 数据模型 ====================

// User 用户模型
type User struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	Email     string    `json:"email"`
	Role      string    `json:"role"` // admin, user, guest
	CreatedAt time.Time `json:"created_at"`
}

// Claims JWT Claims
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// Response 统一响应
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ==================== 全局变量 ====================

var (
	users          = make(map[uint]*User)
	nextID    uint = 1
	usersMu   sync.RWMutex
	jwtSecret = []byte("your-secret-key-change-in-production")
)

// ==================== 主函数 ====================

func main() {
	r := gin.Default()

	// TODO: 在这里注册全局中间件
	// 建议顺序：Recovery -> Logger -> RateLimit

	// 公开路由组（无需认证）
	public := r.Group("/api/v1")
	{
		// TODO: 实现注册接口
		// POST /api/v1/auth/register

		// TODO: 实现登录接口
		// POST /api/v1/auth/login

		// TODO: 实现刷新 Token 接口
		// POST /api/v1/auth/refresh

		// 公开信息接口
		public.GET("/public/info", func(c *gin.Context) {
			c.JSON(http.StatusOK, Response{
				Code:    200,
				Message: "success",
				Data:    gin.H{"info": "This is public information"},
			})
		})
	}

	// 受保护路由组（需要认证）
	// TODO: 添加 JWTAuth 中间件
	protected := r.Group("/api/v1")
	// protected.Use(JWTAuth())
	{
		// 用户相关
		user := protected.Group("/user")
		{
			// TODO: GET /user/profile - 获取当前用户信息
			// 从 Context 中获取 user_id，返回用户信息

			// TODO: PUT /user/profile - 更新用户信息
			// 可更新 email
		}

		// 管理员相关（需要 admin 角色）
		// TODO: 添加 AdminOnly 中间件
		admin := protected.Group("/admin")
		// admin.Use(AdminOnly())
		{
			// TODO: GET /admin/users - 获取所有用户列表

			// TODO: DELETE /admin/users/:id - 删除用户
		}
	}

	r.Run(":8080")
}

// ==================== 中间件实现区域 ====================

// RecoveryMiddleware 错误恢复中间件
// TODO: 实现 panic 恢复，防止程序崩溃
// 捕获 panic 后返回 500 错误，并记录日志
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: 使用 defer 捕获 panic
		c.Next()
	}
}

// LoggerMiddleware 请求日志中间件
// TODO: 记录请求信息：方法、路径、IP、耗时、状态码
// 如果用户已登录，还要记录 user_id
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: 实现日志记录
		c.Next()
	}
}

// RateLimitMiddleware 限流中间件
// TODO: 实现基于 IP 的限流，每分钟最多 60 次请求
// 使用内存存储请求记录（可以使用 map[string][]time.Time）
// 超出限制返回 429 状态码
func RateLimitMiddleware() gin.HandlerFunc {
	// TODO: 初始化限流器数据结构
	return func(c *gin.Context) {
		// TODO: 检查并更新请求计数
		c.Next()
	}
}

// JWTAuth JWT认证中间件
// TODO: 验证 Authorization 头中的 Bearer Token
// 验证成功后，将 user_id, username, role 存入 Context
// 验证失败返回 401
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: 实现 JWT 验证
		// 1. 获取 Authorization 头
		// 2. 提取 Bearer Token
		// 3. 解析并验证 Token
		// 4. 将用户信息存入 Context
		// 5. 验证失败返回 401
		c.Next()
	}
}

// AdminOnly 管理员权限中间件
// TODO: 检查用户角色是否为 admin
// 如果不是 admin，返回 403
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: 从 Context 获取 role，检查是否为 admin
		c.Next()
	}
}

// ==================== 辅助函数 ====================

// GenerateToken 生成JWT Token
// TODO: 生成 Access Token（2小时有效期）和 Refresh Token（7天有效期）
func GenerateToken(userID uint, username, role string) (accessToken, refreshToken string, err error) {
	// TODO: 实现 Token 生成
	return "", "", nil
}

// ParseToken 解析JWT Token
// TODO: 解析 Token 并返回 Claims
func ParseToken(tokenString string) (*Claims, error) {
	// TODO: 实现 Token 解析
	return nil, nil
}

// GetUserByUsername 根据用户名查找用户
func GetUserByUsername(username string) *User {
	usersMu.RLock()
	defer usersMu.RUnlock()

	for _, user := range users {
		if user.Username == username {
			return user
		}
	}
	return nil
}

// GetUserByID 根据ID查找用户
func GetUserByID(id uint) *User {
	usersMu.RLock()
	defer usersMu.RUnlock()

	if user, exists := users[id]; exists {
		return user
	}
	return nil
}

// GenerateID 生成自增ID
func GenerateID() uint {
	usersMu.Lock()
	defer usersMu.Unlock()
	id := nextID
	nextID++
	return id
}
