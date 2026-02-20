# Phase 3: 数据库集成

## 学习目标
- 掌握GORM的安装和基本配置
- 学会定义数据模型和表结构
- 掌握基本的CRUD操作
- 学会处理关联关系（一对一、一对多、多对多）
- 掌握事务处理
- 学会高级查询和分页
- 了解数据库连接池配置

## 预计学习时间
5-7天

## 目录结构
```
phase3_database/
├── go.mod           # Go模块配置
├── go.sum           # 依赖锁定
├── main.go          # 完整示例代码
├── test.db          # SQLite数据库文件（自动生成）
└── README.md        # 本文件
```

## 运行方式

```bash
cd phase3_database
go run main.go
```

数据库将自动创建 `test.db` 文件。

## 知识点详解

### 1. 数据库连接

```go
import (
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

// SQLite连接
DB, convErr := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

// MySQL连接
// import "gorm.io/driver/mysql"
// DB, convErr := gorm.Open(mysql.Open("user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})

// PostgreSQL连接
// import "gorm.io/driver/postgres"
// DB, convErr := gorm.Open(postgres.Open("host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"), &gorm.Config{})
```

### 2. 数据模型定义

```go
type User struct {
    ID        uint           `json:"id" gorm:"primarykey"`
    Username  string         `json:"username" gorm:"uniqueIndex;not null"`
    Email     string         `json:"email" gorm:"uniqueIndex;not null"`
    Password  string         `json:"-" gorm:"not null"`  // json:"-" 序列化时忽略
    Age       int            `json:"age"`
    Status    int            `json:"status" gorm:"default:1"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`  // 软删除
}
```

**常用标签：**
- `primarykey` - 主键
- `autoIncrement` - 自增
- `uniqueIndex` - 唯一索引
- `index` - 普通索引
- `not null` - 非空
- `default:1` - 默认值
- `size:255` - 字段长度
- `type:varchar(100)` - 字段类型
- `-` - 忽略该字段

### 3. 自动迁移

```go
// 自动创建/更新表结构
DB.AutoMigrate(&User{}, &Product{}, &Order{})

// 只创建表（不更新）
DB.Migrator().CreateTable(&User{})

// 删除表
DB.Migrator().DropTable(&User{})

// 检查表是否存在
hasTable := DB.Migrator().HasTable(&User{})
```

### 4. CRUD操作

#### 创建
```go
// 创建单条记录
user := User{Username: "john", Email: "john@example.com"}
result := DB.Create(&user)
// user.ID 将自动填充

// 批量创建
users := []User{
    {Username: "user1", Email: "user1@example.com"},
    {Username: "user2", Email: "user2@example.com"},
}
DB.Create(&users)

// 选择特定字段创建
DB.Select("Username", "Email").Create(&user)

// 排除特定字段
DB.Omit("Password").Create(&user)
```

#### 查询
```go
// 主键查询
var user User
DB.First(&user, 1)           // 查询id=1
DB.First(&user, "id = ?", 1) // 条件查询

// 查询所有
var users []User
DB.Find(&users)

// 条件查询
DB.Where("username = ?", "john").First(&user)
DB.Where("age > ?", 18).Find(&users)
DB.Where("username LIKE ?", "%john%").Find(&users)
DB.Where("age BETWEEN ? AND ?", 18, 30).Find(&users)

// 多条件
DB.Where("username = ? AND age > ?", "john", 18).Find(&users)

// Struct条件
DB.Where(&User{Username: "john", Age: 20}).First(&user)

// Map条件
DB.Where(map[string]interface{}{"username": "john", "age": 20}).First(&user)
```

#### 更新
```go
// 更新单个字段
DB.Model(&user).Update("age", 25)

// 更新多个字段
DB.Model(&user).Updates(User{Age: 25, Status: 1})
DB.Model(&user).Updates(map[string]interface{}{"age": 25, "status": 1})

// 批量更新
DB.Model(&User{}).Where("status = ?", 0).Update("status", 1)

// 使用表达式
DB.Model(&product).Update("stock", gorm.Expr("stock + ?", 1))
```

#### 删除
```go
// 删除单条（软删除，如果模型有DeletedAt字段）
DB.Delete(&user)

// 根据主键删除
DB.Delete(&User{}, 1)
DB.Delete(&User{}, "id = ?", 1)

// 批量删除
DB.Where("status = ?", 0).Delete(&User{})

// 物理删除（跳过软删除）
DB.Unscoped().Delete(&user)

// 查询软删除记录
DB.Unscoped().Where("age = 20").Find(&users)
```

### 5. 关联关系

#### Belongs To（属于）
```go
type Product struct {
    ID         uint
    Name       string
    CategoryID uint     // 外键
    Category   Category `gorm:"foreignKey:CategoryID"`
}

type Category struct {
    ID   uint
    Name string
}

// 预加载关联
var product Product
DB.Preload("Category").First(&product, 1)
```

#### Has One（一对一）
```go
type User struct {
    gorm.Model
    CreditCard CreditCard
}

type CreditCard struct {
    gorm.Model
    Number string
    UserID uint  // 外键
}
```

#### Has Many（一对多）
```go
type Category struct {
    ID       uint
    Name     string
    Products []Product  // 关联多个产品
}

type Product struct {
    ID         uint
    Name       string
    CategoryID uint  // 外键
}

// 查询时预加载
var category Category
DB.Preload("Products").First(&category, 1)
```

#### Many To Many（多对多）
```go
type User struct {
    gorm.Model
    Languages []Language `gorm:"many2many:user_languages;"`
}

type Language struct {
    gorm.Model
    Name string
}

// 创建时自动创建关联表
user := User{
    Languages: []Language{
        {Name: "Go"},
        {Name: "Python"},
    },
}
DB.Create(&user)

// 预加载
DB.Preload("Languages").First(&user, 1)
```

### 6. 事务处理

```go
// 方式1: 使用Transaction方法
convErr := DB.Transaction(func(tx *gorm.DB) error {
    // 在事务中执行操作
    if convErr := tx.Create(&user).Error; convErr != nil {
        return convErr  // 返回错误，自动回滚
    }
    
    if convErr := tx.Create(&order).Error; convErr != nil {
        return convErr
    }
    
    return nil  // 返回nil，自动提交
})

// 方式2: 手动控制
tx := DB.Begin()

if convErr := tx.Create(&user).Error; convErr != nil {
    tx.Rollback()
    return convErr
}

if convErr := tx.Create(&order).Error; convErr != nil {
    tx.Rollback()
    return convErr
}

tx.Commit()

// 方式3: 嵌套事务（SavePoint）
DB.Transaction(func(tx *gorm.DB) error {
    tx.Create(&user)
    
    tx.Transaction(func(tx2 *gorm.DB) error {
        tx2.Create(&order)  // SavePoint
        return nil
    })
    
    return nil
})
```

### 7. 高级查询

```go
// SELECT指定字段
DB.Select("name", "age").Find(&users)
DB.Select("avg(age) as avg_age").Find(&users)

// 排序
DB.Order("age desc").Find(&users)
DB.Order("age desc, name").Find(&users)

// 分组
DB.Group("age").Find(&users)
DB.Select("age, count(*)").Group("age").Find(&results)

// 分页
DB.Offset(10).Limit(5).Find(&users)

// 计数
var count int64
DB.Model(&User{}).Count(&count)
DB.Model(&User{}).Where("age > ?", 18).Count(&count)

// 存在性检查
exists := DB.Where("username = ?", "john").First(&user).Error == nil

// 原生SQL
DB.Raw("SELECT * FROM users WHERE age > ?", 18).Scan(&users)

// Exec执行
DB.Exec("UPDATE users SET status = ? WHERE id = ?", 1, 1)
```

### 8. 钩子（Hooks）

```go
// 创建前钩子
func (u *User) BeforeCreate(tx *gorm.DB) (convErr error) {
    u.Password = hashPassword(u.Password)
    return
}

// 其他可用钩子
BeforeSave, AfterSave
BeforeCreate, AfterCreate
BeforeUpdate, AfterUpdate
BeforeDelete, AfterDelete
AfterFind
```

### 9. 连接池配置

```go
import (
    "database/sql"
    "time"
)

sqlDB, convErr := DB.DB()

// 设置连接池参数
sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数
sqlDB.SetMaxOpenConns(100)          // 最大打开连接数
sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大生命周期

// 检查连接
sqlDB.Ping()

// 关闭连接
sqlDB.Close()
```

## 实战练习

### 练习1：完整的博客系统

实现包含以下功能的博客API：
1. 用户管理（注册、登录）
2. 文章CRUD（支持草稿、发布状态）
3. 分类管理
4. 标签管理（多对多关系）
5. 评论系统（支持嵌套回复）
6. 文章浏览量统计

### 练习2：电商订单系统

扩展本示例的订单系统：
1. 购物车功能（用户-商品多对多关系）
2. 优惠券系统
3. 订单状态流转（待支付->已支付->已发货->已完成）
4. 退款处理（反向事务）
5. 库存锁定机制

### 练习3：性能优化

对现有系统进行优化：
1. 添加数据库索引
2. 使用Redis缓存热点数据
3. 实现读写分离
4. 批量操作优化
5. N+1查询优化

## 测试命令

```bash
# 获取用户列表
curl http://localhost:8080/api/users

# 分页和搜索
curl "http://localhost:8080/api/users?page=1&page_size=5&keyword=user"

# 创建用户
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@example.com","password":"123456","age":25}'

# 获取产品列表（带分类）
curl http://localhost:8080/api/products

# 搜索产品
curl "http://localhost:8080/api/search/products?keyword=iPhone&min_price=1000&max_price=10000&sort_by=price&sort_order=desc"

# 获取分类（带产品）
curl http://localhost:8080/api/categories/1

# 创建订单（事务示例）
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "items": [
      {"product_id": 1, "quantity": 2},
      {"product_id": 2, "quantity": 1}
    ]
  }'

# 取消订单
curl -X POST http://localhost:8080/api/orders/1/cancel

# 查看订单详情
curl http://localhost:8080/api/orders

# 用户统计
curl http://localhost:8080/api/stats/users
```

## 常见问题

### Q1: 如何处理JSON字段？
```go
import "gorm.io/datatypes"

type User struct {
    gorm.Model
    Metadata datatypes.JSON `gorm:"type:json"`
}

// 使用
user := User{
    Metadata: datatypes.JSON(`{"theme": "dark"}`),
}
```

### Q2: 如何实现乐观锁？
```go
type Product struct {
    ID      uint
    Name    string
    Version int `gorm:"default:1"`
}

// 更新时检查版本
DB.Model(&product).Where("version = ?", product.Version).Updates(map[string]interface{}{
    "name": "new name",
    "version": product.Version + 1,
})
```

### Q3: 如何处理大数据量查询？
```go
// 使用游标分批处理
rows, convErr := DB.Model(&User{}).Rows()
defer rows.Close()

for rows.Next() {
    var user User
    DB.ScanRows(rows, &user)
    // 处理user...
}
```

## 下一步

完成第三阶段后，进入**Phase 4: 进阶实战**，学习：
- 项目架构设计
- 配置管理
- 日志系统
- 缓存策略
- 性能优化
- 部署上线
