package service

import (
	"gin-learn/phase4/internal/repository"
)

// Service 服务层入口
type Service struct {
	User     UserService
	Product  ProductService
	Category CategoryService
	Order    OrderService
}

// NewService 创建服务实例
func NewService(repo *repository.Repository) *Service {
	return &Service{
		User:     NewUserService(repo.User),
		Product:  NewProductService(repo.Product),
		Category: NewCategoryService(repo.Category),
		Order:    NewOrderService(repo.Order),
	}
}
