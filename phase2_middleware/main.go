package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func main() {
	// ========== 1. 中间件基础 ==========
	// 创建纯净引擎，手动添加中间件
	r := gin.New()

	// 使用Recovery中间件（捕获panic）
	r.Use(gin.Recovery())

	// ========== 2. 全局中间件 ==========
	// 自定义日志中间件
	r.Use(LoggerMiddleware())

	// CORS中间件
	r.Use(CORSMiddleware())

	// ========== 3. 路由组中间件 ==========
	// 公开API（无需认证）
	public := r.Group("/public")
	{
		public.GET("/info", publicInfo)
	}

	// 受保护API（需要认证）
	private := r.Group("/private")
	private.Use(AuthMiddleware())
	{
		private.GET("/profile", userProfile)
		private.POST("/data", updateData)
	}

	// ========== 4. 单个路由中间件 ==========
	r.GET("/admin", AdminOnly(), adminDashboard)

	// ========== 5. 数据验证示例 ==========
	r.POST("/register", handleRegister)
	r.POST("/create-product", handleCreateProduct)
	r.POST("/custom-validation", handleCustomValidation)

	// ========== 6. 错误处理示例 ==========
	r.GET("/error", handleError)
	r.GET("/panic", handlePanic)

	// ========== 7. 中间件数据传递 ==========
	r.GET("/trace", TraceMiddleware(), handleTrace)

	r.Run(":8080")
}

// ========== 中间件实现 ==========

// LoggerMiddleware 自定义日志中间件
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 请求前
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		fmt.Printf("[REQUEST] %s %s\n", c.Request.Method, path)

		// 处理请求
		c.Next()

		// 请求后
		latency := time.Since(start)
		clientIP := c.ClientIP()
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		fmt.Printf("[RESPONSE] %d | %v | %s | %s\n",
			statusCode, latency, clientIP, path)
	}
}

// CORSMiddleware 跨域中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers",
			"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "缺少认证令牌",
			})
			c.Abort()
			return
		}

		// 验证token（简化示例）
		if token != "Bearer valid-token" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "无效的认证令牌",
			})
			c.Abort()
			return
		}

		// 将用户信息存入context
		c.Set("user_id", "12345")
		c.Set("username", "admin")

		c.Next()
	}
}

// AdminOnly 管理员权限中间件
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetHeader("X-User-Role")
		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "需要管理员权限",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// TraceMiddleware 链路追踪中间件
func TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 生成或获取trace ID
		traceID := c.GetHeader("X-Trace-ID")
		if traceID == "" {
			traceID = fmt.Sprintf("trace-%d", time.Now().UnixNano())
		}

		// 存入context，供后续使用
		c.Set("trace_id", traceID)
		c.Header("X-Trace-ID", traceID)

		c.Next()
	}
}

// ========== 处理器函数 ==========

func publicInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "公开信息",
		"data":    "任何人都可以访问",
	})
}

func userProfile(c *gin.Context) {
	// 从context获取用户信息
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")

	c.JSON(http.StatusOK, gin.H{
		"user_id":  userID,
		"username": username,
		"profile":  "用户信息详情",
	})
}

func updateData(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "数据更新成功",
	})
}

func adminDashboard(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "管理员仪表板",
		"stats": gin.H{
			"users":   1000,
			"orders":  500,
			"revenue": 99999,
		},
	})
}

func handleTrace(c *gin.Context) {
	traceID, _ := c.Get("trace_id")
	c.JSON(http.StatusOK, gin.H{
		"trace_id": traceID,
		"message":  "请求已追踪",
	})
}

// ========== 数据验证 ==========

// 注册请求结构体
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
	Age      int    `json:"age" binding:"gte=0,lte=150"`
}

func handleRegister(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 解析验证错误
		var errMsgs []string
		if errs, ok := err.(validator.ValidationErrors); ok {
			for _, e := range errs {
				errMsgs = append(errMsgs, fmt.Sprintf("%s: %s", e.Field(), e.Tag()))
			}
		} else {
			errMsgs = append(errMsgs, err.Error())
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "参数验证失败",
			"details": errMsgs,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "注册成功",
		"username": req.Username,
		"email":    req.Email,
	})
}

// 产品请求结构体（自定义验证标签）
type ProductRequest struct {
	Name     string  `json:"name" binding:"required"`
	Price    float64 `json:"price" binding:"required,gt=0"`
	Category string  `json:"category" binding:"required,oneof=electronics clothing food"`
	InStock  bool    `json:"in_stock"`
}

func handleCreateProduct(c *gin.Context) {
	var req ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "产品创建成功",
		"product": req,
	})
}

// 自定义验证结构体
type CustomValidationRequest struct {
	Username string `json:"username" binding:"required"`
	// 自定义验证: 必须以字母开头，只能包含字母数字下划线
	Code string `json:"code" binding:"required,startswithletter,alphanumeric"`
}

// 自定义验证函数
func startswithletter(fl validator.FieldLevel) bool {
	code := fl.Field().String()
	if len(code) == 0 {
		return false
	}
	firstChar := code[0]
	return (firstChar >= 'a' && firstChar <= 'z') || (firstChar >= 'A' && firstChar <= 'Z')
}

func alphanumeric(fl validator.FieldLevel) bool {
	code := fl.Field().String()
	for _, char := range code {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '_') {
			return false
		}
	}
	return true
}

func handleCustomValidation(c *gin.Context) {
	// 注册自定义验证器（实际项目中应在初始化时注册一次）
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("startswithletter", startswithletter)
		v.RegisterValidation("alphanumeric", alphanumeric)
	}

	var req CustomValidationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "验证通过",
		"code":    req.Code,
	})
}

// ========== 错误处理 ==========

func handleError(c *gin.Context) {
	// 模拟业务错误
	if err := doBusinessLogic(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "业务逻辑错误",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "操作成功",
	})
}

func doBusinessLogic() error {
	// 模拟错误
	return fmt.Errorf("数据库连接失败")
}

func handlePanic(c *gin.Context) {
	// 这将触发panic，但会被Recovery中间件捕获
	panic("这是一个故意的panic！")
}

// 需要导入binding包 - 已移到文件开头
