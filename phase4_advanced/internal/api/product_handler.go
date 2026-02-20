package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 产品请求结构体
type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock" binding:"gte=0"`
	CategoryID  uint    `json:"category_id"`
}

// CreateProduct 创建产品
func (s *Server) CreateProduct(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := s.service.Product.CreateProduct(req.Name, req.Description, req.Price, req.Stock, req.CategoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// GetProduct 获取产品详情
func (s *Server) GetProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的产品ID"})
		return
	}

	product, err := s.service.Product.GetProduct(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "产品不存在"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// ListProducts 获取产品列表
func (s *Server) ListProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	categoryID, _ := strconv.ParseUint(c.DefaultQuery("category_id", "0"), 10, 32)
	keyword := c.Query("keyword")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	products, total, err := s.service.Product.ListProducts(page, pageSize, uint(categoryID), keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      products,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// UpdateProduct 更新产品
func (s *Server) UpdateProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的产品ID"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.service.Product.UpdateProduct(uint(id), updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "产品更新成功"})
}

// DeleteProduct 删除产品
func (s *Server) DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的产品ID"})
		return
	}

	if err := s.service.Product.DeleteProduct(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "产品已删除"})
}

// SearchProducts 搜索产品
func (s *Server) SearchProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	minPrice, _ := strconv.ParseFloat(c.DefaultQuery("min_price", "0"), 64)
	maxPrice, _ := strconv.ParseFloat(c.DefaultQuery("max_price", "0"), 64)
	categoryID, _ := strconv.ParseUint(c.DefaultQuery("category_id", "0"), 10, 32)
	keyword := c.Query("keyword")
	sortBy := c.DefaultQuery("sort_by", "id")
	sortOrder := c.DefaultQuery("sort_order", "asc")

	products, total, err := s.service.Product.SearchProducts(minPrice, maxPrice, uint(categoryID), keyword, sortBy, sortOrder, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      products,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
