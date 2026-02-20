package service

import (
	"gin-learn/phase4/internal/model"
	"gin-learn/phase4/internal/repository"
)

// OrderItemInput 订单项输入
type OrderItemInput struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,min=1"`
}

// OrderService 订单服务接口
type OrderService interface {
	CreateOrder(userID uint, items []OrderItemInput) (*model.Order, error)
	GetOrder(id uint) (*model.Order, error)
	ListOrders(page, pageSize int) ([]model.Order, int64, error)
	CancelOrder(id uint) error
}

// orderService 订单服务实现
type orderService struct {
	repo repository.OrderRepository
}

func NewOrderService(repo repository.OrderRepository) OrderService {
	return &orderService{repo: repo}
}

func (s *orderService) CreateOrder(userID uint, items []OrderItemInput) (*model.Order, error) {
	// 转换输入
	repoItems := make([]repository.OrderItemInput, len(items))
	for i, item := range items {
		repoItems[i] = repository.OrderItemInput{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
	}

	return s.repo.CreateOrder(userID, repoItems)
}

func (s *orderService) GetOrder(id uint) (*model.Order, error) {
	return s.repo.GetByID(id)
}

func (s *orderService) ListOrders(page, pageSize int) ([]model.Order, int64, error) {
	return s.repo.List(page, pageSize)
}

func (s *orderService) CancelOrder(id uint) error {
	return s.repo.CancelOrder(id)
}
