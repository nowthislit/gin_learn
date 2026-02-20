# Phase 3 练习：博客系统

## 练习目标
通过实现一个完整的博客系统，巩固 GORM 数据库操作的知识点。

## 系统功能

### 基础功能（必做）

#### 1. 用户管理
- 用户注册/登录
- 用户信息管理
- 用户关注功能（多对多）

#### 2. 文章管理
- 发布文章（标题、内容、分类、标签）
- 编辑文章
- 删除文章（软删除）
- 文章列表（支持分页、分类筛选、标签筛选）
- 文章详情

#### 3. 分类管理
- 创建分类
- 分类树形结构（支持父子分类）
- 文章分类统计

#### 4. 标签系统
- 创建标签
- 文章标签关联（多对多）
- 热门标签统计

#### 5. 评论系统
- 发表评论（支持嵌套回复）
- 评论列表
- 删除评论

### 进阶功能（选做）

#### 6. 文章收藏
- 用户收藏文章
- 收藏列表

#### 7. 文章点赞
- 点赞/取消点赞
- 点赞数统计

#### 8. 文章统计
- 浏览量统计
- 点赞数统计
- 评论数统计

## 数据模型

### User（用户）
```go
type User struct {
    ID        uint      `gorm:"primarykey"`
    Username  string    `gorm:"uniqueIndex;not null"`
    Email     string    `gorm:"uniqueIndex;not null"`
    Password  string    `gorm:"not null"`
    Avatar    string
    Bio       string    // 简介
    CreatedAt time.Time
    UpdatedAt time.Time
    
    Articles  []Article  // 用户的文章
    Comments  []Comment  // 用户的评论
    Following []User     `gorm:"many2many:user_follows;"` // 关注的人
    Followers []User     `gorm:"many2many:user_follows;"` // 粉丝
}
```

### Article（文章）
```go
type Article struct {
    ID          uint      `gorm:"primarykey"`
    Title       string    `gorm:"not null;index"`
    Content     string    `gorm:"type:text;not null"`
    Summary     string    // 摘要
    CoverImage  string    // 封面图
    Status      int       // 0:草稿 1:已发布
    ViewCount   int       // 浏览量
    LikeCount   int       // 点赞数
    CommentCount int      // 评论数
    UserID      uint      // 作者ID
    CategoryID  uint      // 分类ID
    CreatedAt   time.Time
    UpdatedAt   time.Time
    DeletedAt   gorm.DeletedAt `gorm:"index"`
    
    User     User      // 作者
    Category Category  // 分类
    Tags     []Tag     `gorm:"many2many:article_tags;"` // 标签
    Comments []Comment // 评论
}
```

### Category（分类）
```go
type Category struct {
    ID        uint   `gorm:"primarykey"`
    Name      string `gorm:"uniqueIndex;not null"`
    Slug      string `gorm:"uniqueIndex;not null"` // URL友好的名称
    ParentID  *uint  // 父分类ID（支持树形结构）
    Level     int    // 分类层级
    SortOrder int    // 排序
    CreatedAt time.Time
    
    Parent   *Category  // 父分类
    Children []Category `gorm:"foreignKey:ParentID"` // 子分类
    Articles []Article  // 分类下的文章
}
```

### Tag（标签）
```go
type Tag struct {
    ID        uint      `gorm:"primarykey"`
    Name      string    `gorm:"uniqueIndex;not null"`
    CreatedAt time.Time
    
    Articles []Article `gorm:"many2many:article_tags;"`
}
```

### Comment（评论）
```go
type Comment struct {
    ID        uint      `gorm:"primarykey"`
    Content   string    `gorm:"type:text;not null"`
    UserID    uint      // 评论者ID
    ArticleID uint      // 文章ID
    ParentID  *uint     // 父评论ID（支持嵌套）
    CreatedAt time.Time
    
    User     User     // 评论者
    Article  Article  // 所属文章
    Parent   *Comment // 父评论
    Replies  []Comment `gorm:"foreignKey:ParentID"` // 回复
}
```

## API 列表

### 认证相关
- `POST /api/v1/auth/register` - 注册
- `POST /api/v1/auth/login` - 登录

### 文章相关
- `GET /api/v1/articles` - 文章列表（支持分页、分类、标签筛选）
- `GET /api/v1/articles/:id` - 文章详情
- `POST /api/v1/articles` - 发布文章（需要登录）
- `PUT /api/v1/articles/:id` - 编辑文章（需要登录，只能编辑自己的）
- `DELETE /api/v1/articles/:id` - 删除文章（需要登录）
- `GET /api/v1/articles/search` - 搜索文章（按标题和内容）

### 分类相关
- `GET /api/v1/categories` - 分类列表（树形结构）
- `GET /api/v1/categories/:id/articles` - 分类下的文章
- `POST /api/v1/categories` - 创建分类（需要管理员权限）

### 标签相关
- `GET /api/v1/tags` - 标签列表
- `GET /api/v1/tags/:id/articles` - 标签下的文章
- `GET /api/v1/tags/popular` - 热门标签

### 评论相关
- `GET /api/v1/articles/:id/comments` - 文章评论列表
- `POST /api/v1/articles/:id/comments` - 发表评论（需要登录）
- `DELETE /api/v1/comments/:id` - 删除评论（需要登录）

### 用户相关
- `GET /api/v1/users/:id` - 用户主页
- `GET /api/v1/users/:id/articles` - 用户的文章
- `POST /api/v1/users/:id/follow` - 关注用户（需要登录）
- `GET /api/v1/users/following` - 关注的用户列表（需要登录）

## 响应格式
```json
{
    "code": 200,
    "message": "success",
    "data": { }
}
```

## 练习要求

### 1. 数据库设计
- 正确使用 GORM 标签定义模型
- 建立正确的关联关系
- 使用数据库迁移自动创建表

### 2. 查询优化
- 使用 Preload 预加载关联数据
- 避免 N+1 查询问题
- 合理使用索引

### 3. 事务处理
- 文章发布时使用事务（创建文章 + 更新标签关联）
- 评论发布时使用事务

### 4. 软删除
- 文章和评论使用软删除
- 查询时过滤已删除的数据

## 测试示例

```bash
# 1. 注册
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"author1","email":"author1@example.com","password":"123456"}'

# 2. 登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"author1","password":"123456"}'

# 3. 创建分类
curl -X POST http://localhost:8080/api/v1/categories \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"name":"技术","slug":"tech"}'

# 4. 发布文章
curl -X POST http://localhost:8080/api/v1/articles \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "title":"Go语言入门",
    "content":"这是一篇关于Go语言的文章...",
    "summary":"Go语言入门指南",
    "category_id":1,
    "tag_ids":[1,2,3]
  }'

# 5. 获取文章列表
curl "http://localhost:8080/api/v1/articles?page=1&page_size=10&category_id=1"

# 6. 发表评论
curl -X POST http://localhost:8080/api/v1/articles/1/comments \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"content":"写得真好！","parent_id":null}'
```

## 评分标准

| 功能 | 分值 |
|------|------|
| 数据模型设计 | 20分 |
| 用户管理 | 10分 |
| 文章 CRUD | 20分 |
| 分类和标签 | 15分 |
| 评论系统 | 15分 |
| 关联查询 | 10分 |
| 事务处理 | 10分 |
| **总分** | **100分** |

## 提示

### 1. 树形分类查询
```go
// 递归获取分类树
func GetCategoryTree(parentID *uint, level int) ([]Category, error) {
    var categories []Category
    query := db.Where("level = ?", level)
    if parentID != nil {
        query = query.Where("parent_id = ?", *parentID)
    } else {
        query = query.Where("parent_id IS NULL")
    }
    convErr := query.Find(&categories).Error
    return categories, convErr
}
```

### 2. 多对多关联创建
```go
// 创建文章时关联标签
article := Article{
    Title:    "文章标题",
    Content:  "文章内容",
    CategoryID: 1,
    Tags: []Tag{
        {ID: 1},
        {ID: 2},
    },
}
db.Create(&article)
```

### 3. 嵌套评论查询
```go
// 获取文章的评论树
var comments []Comment
db.Where("article_id = ? AND parent_id IS NULL", articleID).
    Preload("User").
    Preload("Replies.User").
    Find(&comments)
```

Good luck! 📝
