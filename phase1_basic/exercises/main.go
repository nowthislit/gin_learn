package main

import (
	"fmt"
	"net/http"
	"net/mail"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// User 用户模型
type User struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // 不返回给前端
	Age       int       `json:"age"`
	Status    int       `json:"status"` // 1:正常 0:已删除
	CreatedAt time.Time `json:"created_at"`
}

// 全局数据存储（使用内存 map）
var (
	users                = make(map[uint]*User)
	nextID  uint         = 1
	usersMu sync.RWMutex // 用于并发安全
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type Request struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
}

func main() {
	r := gin.Default()

	// TODO: 在这里添加路由
	// 提示：使用路由分组 /api/v1/users

	// 示例：健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, Response{
			Code:    200,
			Message: "success",
			Data:    gin.H{"status": "running"},
		})
	})

	// ==================== 练习区域开始 ====================

	// TODO 1: 创建路由组
	// api := r.Group("/api/v1")
	// usersGroup := api.Group("/users")
	api := r.Group("/api/v1")
	userApi := api.Group("/users")

	// TODO 2: 实现用户注册
	// POST /api/v1/users/register
	// 接收：username, email, password, age
	// 验证：username 不能为空且长度 3-20，email 格式，password 至少 6 位
	// 创建用户，ID 自增，Status=1，CreatedAt=time.Now()
	userApi.POST("/register", func(ctx *gin.Context) {
		req := &Request{}
		jErr := ctx.ShouldBindJSON(req)
		if jErr != nil {
			resp := Response{
				Code:    -1,
				Message: jErr.Error(),
			}
			ctx.JSON(http.StatusInternalServerError, resp)
			return
		}
		addr, mErr := mail.ParseAddress(req.Email)
		if mErr != nil {
			ctx.JSON(http.StatusInternalServerError, Response{
				Code:    -1,
				Message: mErr.Error(),
			})
			return
		}
		if req.Username == "" || addr.Address == "" || len(req.Password) < 6 {
			ctx.JSON(http.StatusOK, Response{
				Code:    -1,
				Message: "the request params is invalid",
			})
			return
		}

		// 用户名长度验证
		if len(req.Username) < 3 || len(req.Username) > 20 {
			ctx.JSON(http.StatusBadRequest, Response{
				Code:    -1,
				Message: "username length must be 3-20",
			})
			return
		}

		SaveUser(&User{
			ID:        GenerateID(),
			Username:  req.Username,
			Email:     req.Email,
			Password:  req.Password,
			Status:    1,
			Age:       req.Age,
			CreatedAt: time.Now(),
		})

		ctx.JSON(http.StatusOK, Response{
			Code:    0,
			Message: "register success",
		})
		for _, user := range GetAllUsers() {
			fmt.Println(*user)
		}

	})

	// TODO 3: 实现用户登录
	// POST /api/v1/users/login
	// 接收：username, password
	// 验证用户名密码是否匹配
	userApi.POST("/login", func(ctx *gin.Context) {
		req := &Request{}
		jErr := ctx.ShouldBindJSON(req)
		if jErr != nil {
			resp := Response{
				Code:    -1,
				Message: jErr.Error(),
			}
			ctx.JSON(http.StatusInternalServerError, resp)
			return
		}
		if req.Username == "" || req.Password == "" {
			ctx.JSON(http.StatusOK, Response{
				Code:    -1,
				Message: "the request params is invalid",
			})
			return
		}
		user := GetUserByUsername(req.Username)
		if user == nil {
			ctx.JSON(http.StatusOK, Response{
				Code:    -1,
				Message: "the user does not exist",
			})
			return
		}
		if user.Password != req.Password {
			ctx.JSON(http.StatusOK, Response{
				Code:    -1,
				Message: "the password does not match",
			})
			return
		}
		ctx.JSON(http.StatusOK, Response{
			Code:    0,
			Message: "login success",
		})
	})

	// TODO 4: 获取用户列表
	// GET /api/v1/users
	// 支持 Query 参数：page, page_size, keyword
	// 返回：用户列表（不包含密码）、总数、分页信息
	userApi.GET("/", func(ctx *gin.Context) {
		page, size := ParsePageParams(ctx)

		allUsers := GetAllUsers()
		total := len(allUsers)

		// 计算分页
		start := (page - 1) * size
		if start >= total {
			start = total
		}
		end := start + size
		if end > total {
			end = total
		}

		allUsers = allUsers[start:end]
		ctx.JSON(http.StatusOK, Response{
			Code:    0,
			Message: "success",
			Data: gin.H{
				"users":      allUsers,
				"total":      total,
				"page":       page,
				"page_size":  size,
				"total_page": (total + size - 1) / size,
			},
		})
	})

	// TODO 5: 获取用户详情
	// GET /api/v1/users/:id
	// 根据 ID 查询用户
	// 如果用户不存在或已删除，返回 404
	userApi.GET("/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		if id == "" {
			ctx.JSON(http.StatusOK, Response{
				Code:    -1,
				Message: "the id param is invalid",
			})
		}
		idx, convErr := strconv.Atoi(id)
		if convErr != nil {
			ctx.JSON(http.StatusInternalServerError, Response{
				Code:    -1,
				Message: convErr.Error(),
			})
			return
		}
		user := GetUserByID(uint(idx))
		if user == nil {
			ctx.JSON(http.StatusNotFound, nil)
			return
		}
		ctx.JSON(http.StatusOK, Response{
			Code:    0,
			Message: "success",
			Data:    user,
		})

	})

	// TODO 6: 更新用户信息
	// PUT /api/v1/users/:id
	// 可更新：email, age
	// 用户不存在返回 404
	userApi.PUT("/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		if id == "" {
			ctx.JSON(http.StatusOK, Response{
				Code:    -1,
				Message: "the id param is invalid",
			})
			return
		}
		idx, convErr := strconv.Atoi(id)
		if convErr != nil {
			ctx.JSON(http.StatusInternalServerError, Response{
				Code:    -1,
				Message: convErr.Error(),
			})
			return
		}
		user := GetUserByID(uint(idx))
		if user == nil {
			ctx.JSON(http.StatusOK, Response{
				Code:    -1,
				Message: "the user does not exist",
			})
			return
		}
		req := &Request{}
		bindErr := ctx.ShouldBindJSON(req)
		if bindErr != nil {
			ctx.JSON(http.StatusInternalServerError, Response{
				Code:    -1,
				Message: bindErr.Error(),
			})
			return
		}
		address, parseErr := mail.ParseAddress(req.Email)
		if parseErr != nil {
			ctx.JSON(http.StatusInternalServerError, Response{
				Code:    -1,
				Message: parseErr.Error(),
			})
			return
		}
		if address != nil && address.Address != "" {
			user.Email = address.Address
		}
		if req.Age != user.Age {
			user.Age = req.Age
		}
		ctx.JSON(http.StatusOK, Response{
			Code:    0,
			Message: "update success",
		})

	})

	// TODO 7: 删除用户（软删除）
	// DELETE /api/v1/users/:id
	// 将 Status 设为 0
	// 用户不存在返回 404
	userApi.DELETE("/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		if id == "" {
			ctx.JSON(http.StatusOK, Response{
				Code:    -1,
				Message: "the id param is invalid",
			})
			return
		}
		idx, convErr := strconv.Atoi(id)
		if convErr != nil {
			ctx.JSON(http.StatusInternalServerError, Response{
				Code:    -1,
				Message: convErr.Error(),
			})
			return
		}
		user := GetUserByID(uint(idx))
		if user == nil {
			ctx.JSON(http.StatusNotFound, nil)
			return
		}
		user.Status = 0
		ctx.JSON(http.StatusOK, Response{
			Code:    0,
			Message: "delete success",
		})
	})

	// ==================== 练习区域结束 ====================

	r.Run(":8080")
}

// ==================== 辅助函数（可自行修改）====================

// GenerateID 生成自增ID
func GenerateID() uint {
	usersMu.Lock()
	defer usersMu.Unlock()
	id := nextID
	nextID++
	return id
}

// GetUserByUsername 根据用户名查找用户
func GetUserByUsername(username string) *User {
	usersMu.RLock()
	defer usersMu.RUnlock()

	for _, user := range users {
		if user.Username == username && user.Status == 1 {
			return user
		}
	}
	return nil
}

// GetUserByID 根据ID查找用户
func GetUserByID(id uint) *User {
	usersMu.RLock()
	defer usersMu.RUnlock()

	if user, exists := users[id]; exists && user.Status == 1 {
		return user
	}
	return nil
}

// SaveUser 保存用户
func SaveUser(user *User) {
	usersMu.Lock()
	defer usersMu.Unlock()
	users[user.ID] = user
}

// GetAllUsers 获取所有未删除的用户
func GetAllUsers() []*User {
	usersMu.RLock()
	defer usersMu.RUnlock()

	result := make([]*User, 0)
	for _, user := range users {
		if user.Status == 1 {
			result = append(result, user)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})
	return result
}

// ParsePageParams 解析分页参数
func ParsePageParams(c *gin.Context) (page, pageSize int) {
	page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ = strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return page, pageSize
}
