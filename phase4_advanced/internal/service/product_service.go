package service

import (
	"gin-learn/phase4/internal/model"
	"gin-learn/phase4/internal/repository"
)

// ProductService 产品服务接口
type ProductService interface {
	CreateProduct(name, description string, price float64, stock int, categoryID uint) (*model.Product, error)
	GetProduct(id uint) (*model.Product, error)
	ListProducts(page, pageSize int, categoryID uint, keyword string) ([]model.Product, int64, error)
	UpdateProduct(id uint, updates map[string]interface{}) error
	DeleteProduct(id uint) error
	SearchProducts(minPrice, maxPrice float64, categoryID uint, keyword, sortBy, sortOrder string, page, pageSize int) ([]model.Product, int64, error)
}

// productService 产品服务实现
type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) CreateProduct(name, description string, price float64, stock int, categoryID uint) (*model.Product, error) {
	product := &model.Product{
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       stock,
		CategoryID:  categoryID,
	}

	if err := s.repo.Create(product); err != nil {
		return nil, err
	}

	return s.repo.GetByID(product.ID)
}

func (s *productService) GetProduct(id uint) (*model.Product, error) {
	return s.repo.GetByID(id)
}

func (s *productService) ListProducts(page, pageSize int, categoryID uint, keyword string) ([]model.Product, int64, error) {
	return s.repo.List(page, pageSize, categoryID, keyword)
}

func (s *productService) UpdateProduct(id uint, updates map[string]interface{}) error {
	product, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	// 更新字段
	if name, ok := updates["name"].(string); ok {
		product.Name = name
	}
	if desc, ok := updates["description"].(string); ok {
		product.Description = desc
	}
	if price, ok := updates["price"].(float64); ok {
		product.Price = price
	}
	if stock, ok := updates["stock"].(float64); ok {
		product.Stock = int(stock)
	}
	if catID, ok := updates["category_id"].(float64); ok {
		product.CategoryID = uint(catID)
	}

	return s.repo.Update(product)
}

func (s *productService) DeleteProduct(id uint) error {
	return s.repo.Delete(id)
}

func (s *productService) SearchProducts(minPrice, maxPrice float64, categoryID uint, keyword, sortBy, sortOrder string, page, pageSize int) ([]model.Product, int64, error) {
	return s.repo.Search(minPrice, maxPrice, categoryID, keyword, sortBy, sortOrder, page, pageSize)
}
