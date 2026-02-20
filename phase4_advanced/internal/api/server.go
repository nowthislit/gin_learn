package api

import (
	"fmt"
	"gin-learn/phase4/internal/service"
	"gin-learn/phase4/pkg/logger"

	"github.com/gin-gonic/gin"
)

// Server HTTP服务器
type Server struct {
	router  *gin.Engine
	service *service.Service
}

// NewServer 创建服务器实例
func NewServer(svc *service.Service) *Server {
	// 设置Gin模式
	gin.SetMode(gin.DebugMode)

	r := gin.New()

	// 全局中间件
	r.Use(gin.Recovery())
	r.Use(LoggerMiddleware())
	r.Use(CORSMiddleware())

	server := &Server{
		router:  r,
		service: svc,
	}

	server.setupRoutes()
	return server
}

// setupRoutes 设置路由
func (s *Server) setupRoutes() {
	// 健康检查
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1
	v1 := s.router.Group("/api/v1")
	{
		// 用户路由
		users := v1.Group("/users")
		{
			users.GET("", s.ListUsers)
			users.GET("/:id", s.GetUser)
			users.POST("", s.Register)
			users.PUT("/:id", s.UpdateUser)
			users.DELETE("/:id", s.DeleteUser)
		}

		// 产品路由
		products := v1.Group("/products")
		{
			products.GET("", s.ListProducts)
			products.GET("/:id", s.GetProduct)
			products.POST("", s.CreateProduct)
			products.PUT("/:id", s.UpdateProduct)
			products.DELETE("/:id", s.DeleteProduct)
		}

		// 分类路由
		categories := v1.Group("/categories")
		{
			categories.GET("", s.ListCategories)
			categories.GET("/:id", s.GetCategory)
			categories.POST("", s.CreateCategory)
		}

		// 订单路由
		orders := v1.Group("/orders")
		{
			orders.GET("", s.ListOrders)
			orders.POST("", s.CreateOrder)
			orders.GET("/:id", s.GetOrder)
			orders.POST("/:id/cancel", s.CancelOrder)
		}

		// 搜索路由
		v1.GET("/search/products", s.SearchProducts)
	}
}

// Run 启动服务器
func (s *Server) Run(port int) error {
	addr := fmt.Sprintf(":%d", port)
	logger.Info("Starting HTTP server", logger.String("address", addr))
	return s.router.Run(addr)
}
