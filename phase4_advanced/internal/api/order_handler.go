package api

import (
	"net/http"
	"strconv"

	"gin-learn/phase4/internal/service"

	"github.com/gin-gonic/gin"
)

// 订单请求结构体
type CreateOrderRequest struct {
	UserID uint                     `json:"user_id" binding:"required"`
	Items  []service.OrderItemInput `json:"items" binding:"required,min=1"`
}

// CreateOrder 创建订单
func (s *Server) CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := s.service.Order.CreateOrder(req.UserID, req.Items)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

// GetOrder 获取订单详情
func (s *Server) GetOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的订单ID"})
		return
	}

	order, err := s.service.Order.GetOrder(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// ListOrders 获取订单列表
func (s *Server) ListOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	orders, total, err := s.service.Order.ListOrders(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      orders,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// CancelOrder 取消订单
func (s *Server) CancelOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的订单ID"})
		return
	}

	if err := s.service.Order.CancelOrder(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "订单已取消"})
}
