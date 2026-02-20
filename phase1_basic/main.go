package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 创建默认的Gin引擎（带Logger和Recovery中间件）
	r := gin.Default()

	// ========== 1. Hello World ==========
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, Gin!",
		})
	})

	// ========== 2. HTTP方法路由 ==========
	r.GET("/get", handleGet)
	r.POST("/post", handlePost)
	r.PUT("/put", handlePut)
	r.DELETE("/delete", handleDelete)

	// ========== 3. 路由参数 ==========
	// :param 必选参数
	r.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.JSON(http.StatusOK, gin.H{
			"user": name,
		})
	})

	// *path 通配符参数
	r.GET("/files/*filepath", func(c *gin.Context) {
		path := c.Param("filepath")
		c.JSON(http.StatusOK, gin.H{
			"path": path,
		})
	})

	// ========== 4. 路由分组 ==========
	v1 := r.Group("/api/v1")
	{
		v1.GET("/users", getUsersV1)
		v1.GET("/users/:id", getUserByIDV1)
		v1.POST("/users", createUserV1)
	}

	v2 := r.Group("/api/v2")
	{
		v2.GET("/users", getUsersV2)
	}

	// ========== 5. Query参数 ==========
	r.GET("/welcome", handleQuery)

	// ========== 6. Form参数 ==========
	r.POST("/form", handleForm)

	// ========== 7. JSON绑定 ==========
	r.POST("/login", handleLogin)

	// ========== 8. 响应类型 ==========
	r.GET("/response", handleResponse)

	// 启动服务器
	r.Run(":8080")
}

// HTTP方法处理器
func handleGet(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"method": "GET"})
}

func handlePost(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"method": "POST"})
}

func handlePut(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"method": "PUT"})
}

func handleDelete(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"method": "DELETE"})
}

// API v1 处理器
func getUsersV1(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"api":   "v1",
		"users": []string{"user1", "user2"},
	})
}

func getUserByIDV1(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"api":     "v1",
		"user_id": id,
	})
}

func createUserV1(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"api":     "v1",
		"message": "用户创建成功",
	})
}

// API v2 处理器
func getUsersV2(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"api":         "v2",
		"users":       []string{"user1", "user2", "user3"},
		"new_feature": true,
	})
}

// Query参数处理器
func handleQuery(c *gin.Context) {
	name := c.Query("name")
	age := c.DefaultQuery("age", "18")
	city, exists := c.GetQuery("city")
	if !exists {
		city = "未知"
	}

	c.JSON(http.StatusOK, gin.H{
		"name": name,
		"age":  age,
		"city": city,
	})
}

// Form参数处理器
func handleForm(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	remember := c.DefaultPostForm("remember", "false")

	c.JSON(http.StatusOK, gin.H{
		"username": username,
		"password": password,
		"remember": remember,
	})
}

// 登录请求结构体
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// JSON登录处理器
func handleLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if req.Username == "admin" && req.Password == "123456" {
		c.JSON(http.StatusOK, gin.H{
			"message":  "登录成功",
			"username": req.Username,
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "用户名或密码错误",
		})
	}
}

// 响应类型处理器
func handleResponse(c *gin.Context) {
	// 获取Accept头部判断响应格式
	accept := c.GetHeader("Accept")

	if accept == "application/xml" {
		c.XML(http.StatusOK, gin.H{
			"message": "XML响应",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "JSON响应",
		"support": []string{"JSON", "XML", "String", "HTML", "File"},
	})
}
