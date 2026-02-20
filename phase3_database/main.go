package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 数据库连接
var DB *gorm.DB

// ========== 数据模型 ==========

// User 用户模型
type User struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	Username  string         `json:"username" gorm:"uniqueIndex;not null"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"not null"` // json:"-" 不返回给前端
	Age       int            `json:"age"`
	Status    int            `json:"status" gorm:"default:1"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"` // 软删除
}

// Product 产品模型
type Product struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	Price       float64   `json:"price" gorm:"not null"`
	Stock       int       `json:"stock" gorm:"default:0"`
	CategoryID  uint      `json:"category_id"`
	Category    Category  `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Category 分类模型
type Category struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	Name        string    `json:"name" gorm:"uniqueIndex;not null"`
	Description string    `json:"description"`
	Products    []Product `json:"products,omitempty" gorm:"foreignKey:CategoryID"`
	CreatedAt   time.Time `json:"created_at"`
}

// Order 订单模型
type Order struct {
	ID        uint        `json:"id" gorm:"primarykey"`
	UserID    uint        `json:"user_id"`
	User      User        `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Total     float64     `json:"total"`
	Status    string      `json:"status" gorm:"default:'pending'"`
	Items     []OrderItem `json:"items,omitempty" gorm:"foreignKey:OrderID"`
	CreatedAt time.Time   `json:"created_at"`
}

// OrderItem 订单项
type OrderItem struct {
	ID        uint    `json:"id" gorm:"primarykey"`
	OrderID   uint    `json:"order_id"`
	ProductID uint    `json:"product_id"`
	Product   Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"` // 下单时的价格
}

func main() {
	// 初始化数据库
	initDB()

	// 设置Gin
	r := gin.Default()

	// 用户API
	userGroup := r.Group("/api/users")
	{
		userGroup.GET("", listUsers)
		userGroup.GET("/:id", getUser)
		userGroup.POST("", createUser)
		userGroup.PUT("/:id", updateUser)
		userGroup.DELETE("/:id", deleteUser)
	}

	// 产品API
	productGroup := r.Group("/api/products")
	{
		productGroup.GET("", listProducts)
		productGroup.GET("/:id", getProduct)
		productGroup.POST("", createProduct)
		productGroup.PUT("/:id", updateProduct)
		productGroup.DELETE("/:id", deleteProduct)
	}

	// 分类API
	categoryGroup := r.Group("/api/categories")
	{
		categoryGroup.GET("", listCategories)
		categoryGroup.GET("/:id", getCategory)
		categoryGroup.POST("", createCategory)
	}

	// 订单API（包含事务）
	orderGroup := r.Group("/api/orders")
	{
		orderGroup.GET("", listOrders)
		orderGroup.POST("", createOrder)
		orderGroup.POST("/:id/cancel", cancelOrder)
	}

	// 高级查询API
	r.GET("/api/search/products", searchProducts)
	r.GET("/api/stats/users", userStats)

	log.Println("Server running on http://localhost:8080")
	r.Run(":8080")
}

// 初始化数据库
func initDB() {
	var err error

	// 使用SQLite（生产环境可替换为MySQL/PostgreSQL）
	DB, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 打印SQL日志
	})
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	// 自动迁移
	err = DB.AutoMigrate(
		&User{},
		&Category{},
		&Product{},
		&Order{},
		&OrderItem{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// 插入测试数据
	seedData()
}

// 插入测试数据
func seedData() {
	// 清空旧数据
	DB.Exec("DELETE FROM order_items")
	DB.Exec("DELETE FROM orders")
	DB.Exec("DELETE FROM products")
	DB.Exec("DELETE FROM categories")
	DB.Exec("DELETE FROM users")

	// 创建用户
	users := []User{
		{Username: "user1", Email: "user1@example.com", Password: "pass123", Age: 25},
		{Username: "user2", Email: "user2@example.com", Password: "pass123", Age: 30},
	}
	for i := range users {
		DB.Create(&users[i])
	}

	// 创建分类
	categories := []Category{
		{Name: "Electronics", Description: "电子产品"},
		{Name: "Clothing", Description: "服装"},
		{Name: "Food", Description: "食品"},
	}
	for i := range categories {
		DB.Create(&categories[i])
	}

	// 创建产品
	products := []Product{
		{Name: "iPhone 15", Description: "苹果手机", Price: 6999, Stock: 100, CategoryID: 1},
		{Name: "MacBook Pro", Description: "苹果笔记本", Price: 14999, Stock: 50, CategoryID: 1},
		{Name: "T-Shirt", Description: "纯棉T恤", Price: 99, Stock: 200, CategoryID: 2},
		{Name: "Coffee", Description: "咖啡豆", Price: 89, Stock: 500, CategoryID: 3},
	}
	for i := range products {
		DB.Create(&products[i])
	}

	log.Println("Test data seeded successfully")
}

// ========== 用户API处理器 ==========

func listUsers(c *gin.Context) {
	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// 搜索参数
	keyword := c.Query("keyword")

	// 构建查询
	query := DB.Model(&User{})
	if keyword != "" {
		query = query.Where("username LIKE ? OR email LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%")
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 分页查询
	var users []User
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Find(&users)

	c.JSON(http.StatusOK, gin.H{
		"data":      users,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func getUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var user User
	if err := DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func createUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func updateUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var user User
	if err := DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 不允许更新ID和密码（这里简化处理）
	delete(input, "id")
	delete(input, "password")

	if err := DB.Model(&user).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func deleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := DB.Delete(&User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户已删除"})
}

// ========== 产品API处理器 ==========

func listProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	categoryID := c.Query("category_id")

	query := DB.Model(&Product{})
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	var total int64
	query.Count(&total)

	var products []Product
	offset := (page - 1) * pageSize
	query.Preload("Category").Offset(offset).Limit(pageSize).Find(&products)

	c.JSON(http.StatusOK, gin.H{
		"data":      products,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func getProduct(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var product Product
	if err := DB.Preload("Category").First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "产品不存在"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func createProduct(c *gin.Context) {
	var product Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回时包含分类信息
	DB.Preload("Category").First(&product, product.ID)

	c.JSON(http.StatusCreated, product)
}

func updateProduct(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var product Product
	if err := DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "产品不存在"})
		return
	}

	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	delete(input, "id")

	if err := DB.Model(&product).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	DB.Preload("Category").First(&product, product.ID)
	c.JSON(http.StatusOK, product)
}

func deleteProduct(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := DB.Delete(&Product{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "产品已删除"})
}

// ========== 分类API处理器 ==========

func listCategories(c *gin.Context) {
	var categories []Category
	DB.Find(&categories)
	c.JSON(http.StatusOK, categories)
}

func getCategory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var category Category
	if err := DB.Preload("Products").First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "分类不存在"})
		return
	}

	c.JSON(http.StatusOK, category)
}

func createCategory(c *gin.Context) {
	var category Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, category)
}

// ========== 订单API处理器（事务示例） ==========

func listOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var total int64
	DB.Model(&Order{}).Count(&total)

	var orders []Order
	offset := (page - 1) * pageSize
	DB.Preload("User").
		Preload("Items").
		Preload("Items.Product").
		Offset(offset).
		Limit(pageSize).
		Find(&orders)

	c.JSON(http.StatusOK, gin.H{
		"data":      orders,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func createOrder(c *gin.Context) {
	type OrderItemInput struct {
		ProductID uint `json:"product_id" binding:"required"`
		Quantity  int  `json:"quantity" binding:"required,min=1"`
	}

	type CreateOrderInput struct {
		UserID uint             `json:"user_id" binding:"required"`
		Items  []OrderItemInput `json:"items" binding:"required,min=1"`
	}

	var input CreateOrderInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 使用事务
	var order Order
	err := DB.Transaction(func(tx *gorm.DB) error {
		// 1. 验证用户
		var user User
		if err := tx.First(&user, input.UserID).Error; err != nil {
			return fmt.Errorf("用户不存在")
		}

		// 2. 创建订单
		order = Order{
			UserID: input.UserID,
			Status: "pending",
			Total:  0,
		}
		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		// 3. 处理订单项
		var total float64
		for _, itemInput := range input.Items {
			// 获取产品
			var product Product
			if err := tx.First(&product, itemInput.ProductID).Error; err != nil {
				return fmt.Errorf("产品不存在: %d", itemInput.ProductID)
			}

			// 检查库存
			if product.Stock < itemInput.Quantity {
				return fmt.Errorf("库存不足: %s", product.Name)
			}

			// 扣减库存
			if err := tx.Model(&product).Update("stock", gorm.Expr("stock - ?", itemInput.Quantity)).Error; err != nil {
				return err
			}

			// 创建订单项
			orderItem := OrderItem{
				OrderID:   order.ID,
				ProductID: itemInput.ProductID,
				Quantity:  itemInput.Quantity,
				Price:     product.Price,
			}
			if err := tx.Create(&orderItem).Error; err != nil {
				return err
			}

			total += product.Price * float64(itemInput.Quantity)
		}

		// 4. 更新订单总价
		if err := tx.Model(&order).Update("total", total).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 加载关联数据返回
	DB.Preload("User").
		Preload("Items").
		Preload("Items.Product").
		First(&order, order.ID)

	c.JSON(http.StatusCreated, order)
}

func cancelOrder(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	err := DB.Transaction(func(tx *gorm.DB) error {
		// 1. 获取订单
		var order Order
		if err := tx.First(&order, id).Error; err != nil {
			return fmt.Errorf("订单不存在")
		}

		if order.Status == "cancelled" {
			return fmt.Errorf("订单已取消")
		}

		if order.Status == "completed" {
			return fmt.Errorf("订单已完成，无法取消")
		}

		// 2. 获取订单项
		var items []OrderItem
		if err := tx.Where("order_id = ?", id).Find(&items).Error; err != nil {
			return err
		}

		// 3. 恢复库存
		for _, item := range items {
			if err := tx.Model(&Product{}).
				Where("id = ?", item.ProductID).
				Update("stock", gorm.Expr("stock + ?", item.Quantity)).Error; err != nil {
				return err
			}
		}

		// 4. 更新订单状态
		if err := tx.Model(&order).Update("status", "cancelled").Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "订单已取消"})
}

// ========== 高级查询处理器 ==========

func searchProducts(c *gin.Context) {
	keyword := c.Query("keyword")
	minPrice, _ := strconv.ParseFloat(c.Query("min_price"), 64)
	maxPrice, _ := strconv.ParseFloat(c.Query("max_price"), 64)
	categoryID := c.Query("category_id")
	sortBy := c.DefaultQuery("sort_by", "id")
	sortOrder := c.DefaultQuery("sort_order", "asc")

	query := DB.Model(&Product{})

	// 关键词搜索
	if keyword != "" {
		query = query.Where("name LIKE ? OR description LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%")
	}

	// 价格范围
	if minPrice > 0 {
		query = query.Where("price >= ?", minPrice)
	}
	if maxPrice > 0 {
		query = query.Where("price <= ?", maxPrice)
	}

	// 分类过滤
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	// 排序
	orderStr := sortBy + " " + sortOrder
	query = query.Order(orderStr)

	// 分页
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	var total int64
	query.Count(&total)

	var products []Product
	query.Preload("Category").Offset(offset).Limit(pageSize).Find(&products)

	c.JSON(http.StatusOK, gin.H{
		"data":      products,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func userStats(c *gin.Context) {
	var totalUsers int64
	var avgAge float64
	var activeUsers int64

	DB.Model(&User{}).Count(&totalUsers)
	DB.Model(&User{}).Select("AVG(age)").Scan(&avgAge)
	DB.Model(&User{}).Where("status = ?", 1).Count(&activeUsers)

	c.JSON(http.StatusOK, gin.H{
		"total_users":  totalUsers,
		"average_age":  avgAge,
		"active_users": activeUsers,
	})
}
