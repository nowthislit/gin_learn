package main

import (
	"net/http"
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

	// TODO 2: 实现用户注册
	// POST /api/v1/users/register
	// 接收：username, email, password, age
	// 验证：username 不能为空且长度 3-20，email 格式，password 至少 6 位
	// 创建用户，ID 自增，Status=1，CreatedAt=time.Now()

	// TODO 3: 实现用户登录
	// POST /api/v1/users/login
	// 接收：username, password
	// 验证用户名密码是否匹配

	// TODO 4: 获取用户列表
	// GET /api/v1/users
	// 支持 Query 参数：page, page_size, keyword
	// 返回：用户列表（不包含密码）、总数、分页信息

	// TODO 5: 获取用户详情
	// GET /api/v1/users/:id
	// 根据 ID 查询用户
	// 如果用户不存在或已删除，返回 404

	// TODO 6: 更新用户信息
	// PUT /api/v1/users/:id
	// 可更新：email, age
	// 用户不存在返回 404

	// TODO 7: 删除用户（软删除）
	// DELETE /api/v1/users/:id
	// 将 Status 设为 0
	// 用户不存在返回 404

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
