package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TODO: 按照分层架构实现电商系统
// 参考主目录的 phase4_advanced 完整项目结构

func main() {
	// 1. 初始化配置
	// TODO: 使用 Viper 加载配置文件

	// 2. 初始化日志
	// TODO: 使用 Zap 初始化日志系统

	// 3. 初始化数据库
	db, err := gorm.Open(sqlite.Open("shop.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	// 4. 自动迁移
	// TODO: 迁移所有模型

	// 5. 初始化各层
	// TODO: 初始化 Repository、Service、Handler

	// 6. 设置路由
	r := gin.Default()

	// TODO: 注册路由
	// 用户模块 /api/v1/auth/*
	// 商品模块 /api/v1/products/*
	// 购物车模块 /api/v1/cart/*
	// 订单模块 /api/v1/orders/*

	// 7. 启动服务
	r.Run(":8080")
}

// ==================== 练习清单 ====================

// TODO 1: 定义所有数据模型（model/）
// - User（用户）
// - Product（商品）
// - Category（分类）
// - Cart（购物车）
// - Order（订单）
// - OrderItem（订单商品）
// - Address（收货地址）

// TODO 2: 实现配置管理（config/）
// - 使用 Viper 读取配置文件
// - 支持数据库、Redis、日志等配置

// TODO 3: 实现日志系统（pkg/logger/）
// - 使用 Zap 封装日志工具
// - 支持不同日志级别

// TODO 4: 实现数据访问层（repository/）
// - UserRepository
// - ProductRepository
// - CartRepository
// - OrderRepository
// - 使用接口定义，便于测试

// TODO 5: 实现业务逻辑层（service/）
// - UserService（用户相关）
// - ProductService（商品相关）
// - CartService（购物车相关）
// - OrderService（订单相关，包含事务）

// TODO 6: 实现API层（api/）
// - AuthHandler（登录注册）
// - UserHandler（用户管理）
// - ProductHandler（商品管理）
// - CartHandler（购物车）
// - OrderHandler（订单管理）

// TODO 7: 实现中间件（api/middleware/）
// - JWTAuth（认证中间件）
// - Logger（日志中间件）
// - RateLimit（限流中间件）
// - Recovery（错误恢复）

// TODO 8: 实现缓存（可选）
// - 使用 Redis 缓存商品信息
// - 使用 Redis 缓存购物车

// TODO 9: 编写单元测试
// - 对 Service 层编写测试
// - 使用 mock 数据

// ==================== 核心功能实现提示 ====================

// 1. 订单创建流程（带事务）
// func CreateOrder(userID uint, cartItemIDs []uint, addressID uint) (*Order, error) {
//     return db.Transaction(func(tx *gorm.DB) error {
//         // 1. 创建订单主表
//         // 2. 查询购物车商品
//         // 3. 扣减库存（乐观锁）
//         // 4. 创建订单商品
//         // 5. 清空购物车
//         // 6. 返回订单
//     })
// }

// 2. 库存扣减（乐观锁）
// func DeductStock(productID uint, quantity int) error {
//     result := db.Model(&Product{}).
//         Where("id = ? AND stock >= ?", productID, quantity).
//         Update("stock", gorm.Expr("stock - ?", quantity))
//
//     if result.RowsAffected == 0 {
//         return errors.New("库存不足")
//     }
//     return nil
// }

// 3. JWT 认证中间件
// func JWTAuth() gin.HandlerFunc {
//     return func(c *gin.Context) {
//         // 1. 获取 Authorization 头
//         // 2. 解析 JWT Token
//         // 3. 将用户信息存入 Context
//         // 4. 验证失败返回 401
//         c.Next()
//     }
// }
