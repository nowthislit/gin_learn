package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	Username  string         `json:"username" gorm:"uniqueIndex;not null;size:50"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null;size:100"`
	Password  string         `json:"-" gorm:"not null;size:255"`
	Age       int            `json:"age"`
	Status    int            `json:"status" gorm:"default:1;index"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Product 产品模型
type Product struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	Name        string    `json:"name" gorm:"not null;size:200;index"`
	Description string    `json:"description" gorm:"size:500"`
	Price       float64   `json:"price" gorm:"not null;index"`
	Stock       int       `json:"stock" gorm:"default:0"`
	CategoryID  uint      `json:"category_id"`
	Category    Category  `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Category 分类模型
type Category struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	Name        string    `json:"name" gorm:"uniqueIndex;not null;size:100"`
	Description string    `json:"description" gorm:"size:500"`
	Products    []Product `json:"products,omitempty" gorm:"foreignKey:CategoryID"`
	CreatedAt   time.Time `json:"created_at"`
}

// Order 订单模型
type Order struct {
	ID        uint        `json:"id" gorm:"primarykey"`
	UserID    uint        `json:"user_id" gorm:"index"`
	User      User        `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Total     float64     `json:"total"`
	Status    string      `json:"status" gorm:"default:'pending';index"`
	Items     []OrderItem `json:"items,omitempty" gorm:"foreignKey:OrderID"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// OrderItem 订单项
type OrderItem struct {
	ID        uint    `json:"id" gorm:"primarykey"`
	OrderID   uint    `json:"order_id" gorm:"index"`
	ProductID uint    `json:"product_id"`
	Product   Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}
