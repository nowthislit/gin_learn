package main

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
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
	r.Use(RecoveryMiddleware()).Use(LoggerMiddleware()).Use(RateLimitMiddleware())
	// 公开路由组（无需认证）
	public := r.Group("/api/v1")
	{
		// TODO: 实现注册接口
		// POST /api/v1/auth/register
		auth := public.Group("/auth")
		auth.POST("/register", func(c *gin.Context) {
			var user User
			if err := c.ShouldBindJSON(&user); err != nil {
				c.JSON(http.StatusBadRequest, Response{
					Code:    -1,
					Message: err.Error(),
				})
				return
			}
			user.ID = GenerateID()
			user.CreatedAt = time.Now()
			usersMu.Lock()
			users[user.ID] = &user
			usersMu.Unlock()
			c.JSON(http.StatusOK, Response{Code: 0, Message: "success", Data: user})
		})

		// TODO: 实现登录接口
		// POST /api/v1/auth/login
		auth.POST("/login", func(c *gin.Context) {
			var user User
			if err := c.ShouldBindJSON(&user); err != nil {
				c.JSON(http.StatusBadRequest, Response{
					Code:    -1,
					Message: err.Error(),
				})
				return
			}
			suser := GetUserByUsername(user.Username)
			if suser == nil {
				c.JSON(http.StatusBadRequest, Response{
					Code:    -1,
					Message: fmt.Sprintf("user %s not found", user.Username),
				})
				return
			}
			accessToken, refreshToken, err := GenerateToken(suser.ID, suser.Username, suser.Role)
			if err != nil {
				c.JSON(http.StatusBadRequest, Response{
					Code:    -1,
					Message: err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, Response{
				Code:    0,
				Message: "success",
				Data:    gin.H{"accessToken": accessToken, "refreshToken": refreshToken, "user": suser},
			})

		})

		// TODO: 实现刷新 Token 接口
		// POST /api/v1/auth/refresh
		auth.POST("/refresh", func(c *gin.Context) {
			var params User
			if err := c.ShouldBindJSON(&params); err != nil {
				c.JSON(http.StatusBadRequest, Response{
					Code:    -1,
					Message: err.Error(),
				})
				return
			}
			user := GetUserByUsername(params.Username)
			if user == nil {
				c.JSON(http.StatusBadRequest, Response{
					Code:    -1,
					Message: fmt.Sprintf("user %s not found", params.Username),
				})
				return
			}
			token, refreshToken, err := GenerateToken(user.ID, user.Username, user.Role)
			if err != nil {
				c.JSON(http.StatusBadRequest, Response{
					Code:    -1,
					Message: err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, Response{
				Code:    0,
				Message: "success",
				Data:    gin.H{"accessToken": token, "refreshToken": refreshToken, "user": user},
			})

		})

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
	protected.Use(JWTAuth())
	{
		// 用户相关
		user := protected.Group("/user")
		{
			// TODO: GET /user/profile - 获取当前用户信息
			// 从 Context 中获取 user_id，返回用户信息
			user.GET("/profile", func(c *gin.Context) {
				userID, exists := c.Get("userID")
				if !exists {
					c.JSON(http.StatusBadRequest, Response{
						Code:    -1,
						Message: "user id not found",
					})
					return
				}
				uid := userID.(uint)
				user := GetUserByID(uid)
				c.JSON(http.StatusOK, Response{
					Code:    0,
					Message: "success",
					Data:    gin.H{"user": user},
				})
			})
			// TODO: PUT /user/profile - 更新用户信息
			// 可更新 email
			user.PUT("/profile", func(c *gin.Context) {
				var params User
				if err := c.ShouldBindJSON(&params); err != nil {
					c.JSON(http.StatusBadRequest, Response{
						Code:    -1,
						Message: err.Error(),
					})
					return
				}
				userID, exists := c.Get("userID")
				if !exists {
					c.JSON(http.StatusBadRequest, Response{
						Code:    -1,
						Message: "user id not found",
					})
					return
				}
				uid := userID.(uint)
				user := GetUserByID(uid)
				if user == nil {
					c.JSON(http.StatusNotFound, Response{
						Code:    -1,
						Message: "user not found",
					})
					return
				}
				user.Email = params.Email
				c.JSON(http.StatusOK, Response{
					Code:    0,
					Message: "success",
					Data:    gin.H{"user": user},
				})
			})
		}

		// 管理员相关（需要 admin 角色）
		// TODO: 添加 AdminOnly 中间件
		admin := protected.Group("/admin")
		// admin.Use(AdminOnly())
		admin.Use(AdminOnly())
		{
			// TODO: GET /admin/users - 获取所有用户列表
			admin.GET("/users", func(c *gin.Context) {
				allUser := make([]*User, 0)
				for _, user := range users {
					allUser = append(allUser, user)
				}
				sort.Slice(allUser, func(i, j int) bool {
					return allUser[i].ID < allUser[j].ID
				})
				c.JSON(http.StatusOK, Response{
					Code:    0,
					Message: "success",
					Data:    allUser,
				})
			})

			// TODO: DELETE /admin/users/:id - 删除用户
			admin.DELETE("/users/:id", func(c *gin.Context) {
				userID := c.Param("id")
				if userID == "" {
					c.JSON(http.StatusOK, Response{
						Code:    -1,
						Message: fmt.Sprintf("param id not found"),
					})
					return
				}

				uid, err := strconv.Atoi(userID)
				if err != nil {
					c.JSON(http.StatusOK, Response{
						Code:    -1,
						Message: fmt.Sprintf("param id not a number"),
					})
					return
				}
				usersMu.Lock()
				delete(users, uint(uid))
				usersMu.Unlock()
				c.JSON(http.StatusOK, Response{
					Code:    0,
					Message: "success",
				})

			})
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
		defer func() {
			if err := recover(); err != nil {
				c.JSON(http.StatusInternalServerError, Response{
					Code:    http.StatusInternalServerError,
					Message: "服务器内部异常",
				})
			}
			c.Abort()
			return
		}()
		c.Next()
	}
}

// LoggerMiddleware 请求日志中间件
// TODO: 记录请求信息：方法、路径、IP、耗时、状态码
// 如果用户已登录，还要记录 user_id
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: 实现日志记录
		start := time.Now()
		c.Next()

		// 获取 user_id（如果已登录）
		userID, _ := c.Get("userID")

		fmt.Printf("[%s] %s %s user_id=%v status=%d duration=%v\n",
			time.Now().Format("2006-01-02 15:04:05"),
			c.Request.Method,
			c.Request.URL.Path,
			userID,
			c.Writer.Status(),
			time.Since(start),
		)
	}
}

// RateLimitMiddleware 限流中间件
// TODO: 实现基于 IP 的限流，每分钟最多 60 次请求
// 使用内存存储请求记录（可以使用 map[string][]time.Time）
// 超出限制返回 429 状态码
func RateLimitMiddleware() gin.HandlerFunc {
	// TODO: 初始化限流器数据结构
	var rateLimit = make(map[string][]time.Time)
	return func(c *gin.Context) {
		// TODO: 检查并更新请求计数
		ip := c.ClientIP()
		now := time.Now()
		beforeMin := now.Add(-time.Minute)

		// 清理过期记录并统计有效请求
		validTimes := make([]time.Time, 0)
		for _, t := range rateLimit[ip] {
			if t.After(beforeMin) {
				validTimes = append(validTimes, t)
			}
		}

		if len(validTimes) >= 60 {
			c.JSON(http.StatusTooManyRequests, Response{
				Code:    http.StatusTooManyRequests,
				Message: "请求频繁",
			})
			c.Abort()
			return
		}

		validTimes = append(validTimes, now)
		rateLimit[ip] = validTimes
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
		authorization := c.GetHeader("Authorization")

		// 2. 提取 Bearer Token
		if authorization == "" {
			c.JSON(http.StatusUnauthorized, Response{
				Code:    http.StatusUnauthorized,
				Message: "the authorization header is empty",
			})
			c.Abort()
			return
		}
		split := strings.Split(authorization, " ")
		if len(split) != 2 {
			c.JSON(http.StatusUnauthorized, Response{
				Code:    http.StatusUnauthorized,
				Message: "the authorization header is invalid",
			})
			return
		}
		token := split[1]
		// 3. 解析并验证 Token
		claims, err := ParseToken(token)
		if err != nil || claims == nil {
			c.JSON(http.StatusUnauthorized, Response{
				Code:    http.StatusUnauthorized,
				Message: "authorization header is invalid",
			})
			c.Abort()
			return
		}
		// 4. 将用户信息存入 Context
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
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
		if role, ok := c.Get("role"); !ok || role != "admin" {
			c.JSON(http.StatusOK, Response{
				Code:    http.StatusForbidden,
				Message: "user is not admin",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// ==================== 辅助函数 ====================

// GenerateToken 生成JWT Token
// TODO: 生成 Access Token（2小时有效期）和 Refresh Token（7天有效期）
func GenerateToken(userID uint, username, role string) (accessToken, refreshToken string, err error) {
	// TODO: 实现 Token 生成
	accessClaims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	refreshClaims := accessClaims
	refreshClaims.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	fmt.Println(accessClaims)
	fmt.Println(refreshClaims)
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenStr, err := at.SignedString(jwtSecret)
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenStr, err := rt.SignedString(jwtSecret)
	return accessTokenStr, refreshTokenStr, err
}

// ParseToken 解析JWT Token
// TODO: 解析 Token 并返回 Claims
func ParseToken(tokenString string) (*Claims, error) {
	// TODO: 实现 Token 解析
	var claims Claims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return &claims, nil
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
