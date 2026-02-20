package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ==================== 数据模型 ====================

// User 用户模型
type User struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	Username  string    `json:"username" gorm:"uniqueIndex;not null"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null"`
	Password  string    `json:"-" gorm:"not null"`
	Avatar    string    `json:"avatar"`
	Bio       string    `json:"bio"` // 简介
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 关联
	Articles  []Article `json:"articles,omitempty" gorm:"foreignKey:UserID"`
	Comments  []Comment `json:"comments,omitempty" gorm:"foreignKey:UserID"`
	Following []User    `json:"following,omitempty" gorm:"many2many:user_follows;"`
	Followers []User    `json:"followers,omitempty" gorm:"many2many:user_follows;"`
}

// Article 文章模型
type Article struct {
	ID           uint           `json:"id" gorm:"primarykey"`
	Title        string         `json:"title" gorm:"not null;index"`
	Content      string         `json:"content" gorm:"type:text;not null"`
	Summary      string         `json:"summary"`
	CoverImage   string         `json:"cover_image"`
	Status       int            `json:"status" gorm:"default:0"` // 0:草稿 1:已发布
	ViewCount    int            `json:"view_count" gorm:"default:0"`
	LikeCount    int            `json:"like_count" gorm:"default:0"`
	CommentCount int            `json:"comment_count" gorm:"default:0"`
	UserID       uint           `json:"user_id"`
	CategoryID   uint           `json:"category_id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联
	User     User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Category Category  `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	Tags     []Tag     `json:"tags,omitempty" gorm:"many2many:article_tags;"`
	Comments []Comment `json:"comments,omitempty" gorm:"foreignKey:ArticleID"`
}

// Category 分类模型
type Category struct {
	ID        uint   `json:"id" gorm:"primarykey"`
	Name      string `json:"name" gorm:"uniqueIndex;not null"`
	Slug      string `json:"slug" gorm:"uniqueIndex;not null"`
	ParentID  *uint  `json:"parent_id"`
	Level     int    `json:"level" gorm:"default:1"`
	SortOrder int    `json:"sort_order" gorm:"default:0"`
	CreatedAt time.Time

	// 关联
	Parent   *Category  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children []Category `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	Articles []Article  `json:"articles,omitempty" gorm:"foreignKey:CategoryID"`
}

// Tag 标签模型
type Tag struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	Name      string    `json:"name" gorm:"uniqueIndex;not null"`
	CreatedAt time.Time `json:"created_at"`

	// 关联
	Articles []Article `json:"articles,omitempty" gorm:"many2many:article_tags;"`
}

// Comment 评论模型
type Comment struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	Content   string         `json:"content" gorm:"type:text;not null"`
	UserID    uint           `json:"user_id"`
	ArticleID uint           `json:"article_id"`
	ParentID  *uint          `json:"parent_id"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联
	User    User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Article Article   `json:"article,omitempty" gorm:"foreignKey:ArticleID"`
	Parent  *Comment  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Replies []Comment `json:"replies,omitempty" gorm:"foreignKey:ParentID"`
}

// Response 统一响应
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func main() {
	// 初始化数据库
	initDB()

	// 设置路由
	r := gin.Default()

	// TODO: 实现以下 API

	// ==================== 认证相关 ====================
	// POST /api/v1/auth/register - 注册
	// POST /api/v1/auth/login - 登录

	// ==================== 文章相关 ====================
	// GET /api/v1/articles - 文章列表（支持分页、分类、标签筛选）
	// GET /api/v1/articles/:id - 文章详情
	// POST /api/v1/articles - 发布文章（需要登录）
	// PUT /api/v1/articles/:id - 编辑文章（需要登录）
	// DELETE /api/v1/articles/:id - 删除文章（需要登录）
	// GET /api/v1/articles/search - 搜索文章

	// ==================== 分类相关 ====================
	// GET /api/v1/categories - 分类列表（树形结构）
	// GET /api/v1/categories/:id/articles - 分类下的文章
	// POST /api/v1/categories - 创建分类（需要管理员权限）

	// ==================== 标签相关 ====================
	// GET /api/v1/tags - 标签列表
	// GET /api/v1/tags/:id/articles - 标签下的文章
	// GET /api/v1/tags/popular - 热门标签

	// ==================== 评论相关 ====================
	// GET /api/v1/articles/:id/comments - 文章评论列表
	// POST /api/v1/articles/:id/comments - 发表评论（需要登录）
	// DELETE /api/v1/comments/:id - 删除评论（需要登录）

	// ==================== 用户相关 ====================
	// GET /api/v1/users/:id - 用户主页
	// GET /api/v1/users/:id/articles - 用户的文章
	// POST /api/v1/users/:id/follow - 关注用户（需要登录）

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, Response{Code: 200, Message: "ok"})
	})

	r.Run(":8080")
}

// initDB 初始化数据库
func initDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("blog.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 自动迁移
	DB.AutoMigrate(
		&User{},
		&Article{},
		&Category{},
		&Tag{},
		&Comment{},
	)

	// TODO: 插入测试数据
}

// ==================== 辅助函数 ====================

// TODO 1: 实现获取分类树
// func GetCategoryTree(parentID *uint) ([]Category, error) {
//     // 递归获取分类树
// }

// TODO 2: 实现获取文章评论树
// func GetArticleComments(articleID uint) ([]Comment, error) {
//     // 获取文章的评论，包含嵌套回复
// }

// TODO 3: 实现发布文章（带事务）
// func CreateArticle(article *Article, tagIDs []uint) error {
//     // 使用事务创建文章并关联标签
// }

// TODO 4: 实现搜索文章
// func SearchArticles(keyword string, page, pageSize int) ([]Article, int64, error) {
//     // 按标题和内容搜索文章
// }
