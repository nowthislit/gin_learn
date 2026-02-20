package repository

import (
	"fmt"
	"gin-learn/phase4/internal/model"

	"gorm.io/gorm"
)

// OrderRepository 订单仓库接口
type OrderRepository interface {
	CreateOrder(userID uint, items []OrderItemInput) (*model.Order, error)
	GetByID(id uint) (*model.Order, error)
	List(page, pageSize int) ([]model.Order, int64, error)
	CancelOrder(id uint) error
}

// OrderItemInput 订单项输入
type OrderItemInput struct {
	ProductID uint
	Quantity  int
}

// orderRepository 订单仓库实现
type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) CreateOrder(userID uint, items []OrderItemInput) (*model.Order, error) {
	var order model.Order

	err := r.db.Transaction(func(tx *gorm.DB) error {
		// 验证用户
		var user model.User
		if err := tx.First(&user, userID).Error; err != nil {
			return fmt.Errorf("用户不存在")
		}

		// 创建订单
		order = model.Order{
			UserID: userID,
			Status: "pending",
			Total:  0,
		}
		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		// 处理订单项
		var total float64
		for _, item := range items {
			var product model.Product
			if err := tx.First(&product, item.ProductID).Error; err != nil {
				return fmt.Errorf("产品不存在: %d", item.ProductID)
			}

			if product.Stock < item.Quantity {
				return fmt.Errorf("库存不足: %s", product.Name)
			}

			// 扣减库存
			if err := tx.Model(&product).Update("stock", gorm.Expr("stock - ?", item.Quantity)).Error; err != nil {
				return err
			}

			// 创建订单项
			orderItem := model.OrderItem{
				OrderID:   order.ID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:     product.Price,
			}
			if err := tx.Create(&orderItem).Error; err != nil {
				return err
			}

			total += product.Price * float64(item.Quantity)
		}

		// 更新订单总价
		if err := tx.Model(&order).Update("total", total).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 加载关联数据
	r.db.Preload("User").Preload("Items").Preload("Items.Product").First(&order, order.ID)

	return &order, nil
}

func (r *orderRepository) GetByID(id uint) (*model.Order, error) {
	var order model.Order
	if err := r.db.Preload("User").Preload("Items").Preload("Items.Product").First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) List(page, pageSize int) ([]model.Order, int64, error) {
	var total int64
	r.db.Model(&model.Order{}).Count(&total)

	var orders []model.Order
	offset := (page - 1) * pageSize
	if err := r.db.Preload("User").Preload("Items").Preload("Items.Product").Offset(offset).Limit(pageSize).Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *orderRepository) CancelOrder(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var order model.Order
		if err := tx.First(&order, id).Error; err != nil {
			return fmt.Errorf("订单不存在")
		}

		if order.Status == "cancelled" {
			return fmt.Errorf("订单已取消")
		}

		if order.Status == "completed" {
			return fmt.Errorf("订单已完成，无法取消")
		}

		var items []model.OrderItem
		if err := tx.Where("order_id = ?", id).Find(&items).Error; err != nil {
			return err
		}

		// 恢复库存
		for _, item := range items {
			if err := tx.Model(&model.Product{}).Where("id = ?", item.ProductID).Update("stock", gorm.Expr("stock + ?", item.Quantity)).Error; err != nil {
				return err
			}
		}

		// 更新订单状态
		if err := tx.Model(&order).Update("status", "cancelled").Error; err != nil {
			return err
		}

		return nil
	})
}
