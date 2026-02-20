# Phase 1: 基础入门

## 学习目标
- 掌握Gin框架的安装和基本使用
- 理解Gin引擎的创建和配置
- 学会定义各种HTTP路由
- 掌握路由参数和路由分组
- 学会处理Query、Form、JSON等请求数据
- 学会返回各种响应格式

## 预计学习时间
3-5天

## 目录结构
```
phase1_basic/
├── go.mod       # Go模块配置
├── go.sum       # 依赖锁定
├── main.go      # 完整示例代码
└── README.md    # 本文件
```

## 运行方式

```bash
cd phase1_basic
go run main.go

# 或者编译后运行
go build -o phase1.exe
./phase1.exe
```

## 知识点详解

### 1. 创建Gin引擎

```go
// 方式1: 默认引擎（推荐）
// 自带Logger和Recovery中间件
r := gin.Default()

// 方式2: 纯净引擎
// 不带任何中间件
r := gin.New()
```

### 2. 基本路由

```go
// GET请求
r.GET("/path", func(c *gin.Context) {
    c.JSON(200, gin.H{"key": "value"})
})

// POST/PUT/DELETE
r.POST("/path", handler)
r.PUT("/path", handler)
r.DELETE("/path", handler)

// 匹配所有方法
r.Any("/path", handler)

// 静态文件
r.Static("/static", "./static")
r.StaticFile("/favicon.ico", "./favicon.ico")
```

### 3. 路由参数

```go
// :param - 必选参数
// URL: /user/john
r.GET("/user/:name", func(c *gin.Context) {
    name := c.Param("name")  // john
})

// *param - 通配符参数
// URL: /files/images/photo.jpg
r.GET("/files/*filepath", func(c *gin.Context) {
    path := c.Param("filepath")  // /images/photo.jpg
})
```

### 4. 路由分组

```go
// 创建分组
v1 := r.Group("/api/v1")
{
    v1.GET("/users", getUsers)
    v1.POST("/users", createUser)
}

// 分组使用中间件
authorized := r.Group("/admin", authMiddleware())
{
    authorized.GET("/dashboard", dashboard)
}
```

### 5. 请求数据处理

#### Query参数
```go
// URL: /welcome?name=john&age=20
func handler(c *gin.Context) {
    // 获取参数，不存在返回空字符串
    name := c.Query("name")
    
    // 获取参数，不存在返回默认值
    age := c.DefaultQuery("age", "18")
    
    // 获取参数，返回(value, exists)
    city, ok := c.GetQuery("city")
}
```

#### Form参数
```go
// Content-Type: application/x-www-form-urlencoded
func handler(c *gin.Context) {
    username := c.PostForm("username")
    password := c.PostForm("password")
    remember := c.DefaultPostForm("remember", "false")
}
```

#### JSON参数
```go
type LoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

func handler(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    // 处理req...
}
```

### 6. 响应数据

```go
// JSON
c.JSON(200, gin.H{"message": "success"})

// String
c.String(200, "Hello %s", name)

// XML
c.XML(200, gin.H{"message": "success"})

// HTML
c.HTML(200, "template.html", gin.H{})

// 重定向
c.Redirect(301, "/new-path")

// 文件
c.File("/path/to/file")
```

### 7. Context上下文

`gin.Context`是Gin的核心，包含了请求和响应的所有信息：

```go
func handler(c *gin.Context) {
    // 请求信息
    c.Request        // *http.Request
    c.ClientIP()     // 客户端IP
    c.ContentType()  // Content-Type
    
    // 路径和参数
    c.FullPath()     // 完整路径
    c.Param("id")    // 路由参数
    c.Query("name")  // Query参数
    
    // 响应操作
    c.Status(200)           // 设置状态码
    c.Header("X-Key", "v")  // 设置响应头
    c.SetCookie(...)        // 设置Cookie
    
    // 数据存储（在同一次请求中共享数据）
    c.Set("user_id", 123)
    userID, exists := c.Get("user_id")
}
```

## 实战练习

### 练习1：用户管理系统

创建一个简单的用户管理API，包含以下功能：
1. 获取所有用户列表（GET /users）
2. 获取单个用户（GET /users/:id）
3. 创建用户（POST /users）
4. 更新用户（PUT /users/:id）
5. 删除用户（DELETE /users/:id）

提示：可以使用内存map存储用户数据。

### 练习2：带分页和搜索的API

扩展用户列表接口，支持：
1. 分页参数（page, page_size）
2. 搜索参数（name）
3. 排序参数（sort_by, order）

### 练习3：多版本API

创建v1和v2版本的用户API，v2版本返回更详细的信息。

## 测试命令

```bash
# 测试Hello World
curl http://localhost:8080/

# 测试HTTP方法
curl -X GET http://localhost:8080/get
curl -X POST http://localhost:8080/post

# 测试路由参数
curl http://localhost:8080/user/john
curl http://localhost:8080/files/images/photo.jpg

# 测试路由分组
curl http://localhost:8080/api/v1/users
curl http://localhost:8080/api/v2/users

# 测试Query参数
curl "http://localhost:8080/welcome?name=john&age=20"

# 测试Form参数
curl -X POST http://localhost:8080/form \
  -d "username=admin&password=123456"

# 测试JSON登录
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"123456"}'

# 测试不同响应格式
curl http://localhost:8080/response
curl -H "Accept: application/xml" http://localhost:8080/response
```

## 常见问题

### Q1: 如何获取多个同名参数？
```go
// URL: /search?tag=go&tag=gin
func handler(c *gin.Context) {
    tags := c.QueryArray("tag")  // ["go", "gin"]
}
```

### Q2: 如何处理multipart/form-data？
```go
// 单文件
file, _ := c.FormFile("file")
c.SaveUploadedFile(file, dst)

// 多文件
form, _ := c.MultipartForm()
files := form.File["files"]
```

### Q3: 如何设置全局状态码？
```go
// 设置状态码但不返回
c.Status(404)
// 然后返回数据
c.JSON(404, gin.H{"error": "not found"})
```

## 下一步

完成第一阶段后，进入**Phase 2: 中间件与验证**，学习：
- 内置中间件使用
- 自定义中间件开发
- 请求数据验证
- 全局错误处理
