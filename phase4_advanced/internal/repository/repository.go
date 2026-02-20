package repository

import (
	"fmt"

	"gin-learn/phase4/config"
	"gin-learn/phase4/internal/model"
	"gin-learn/phase4/pkg/logger"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var db *gorm.DB

// InitDB 初始化数据库连接
func InitDB() (*gorm.DB, error) {
	var err error

	cfg := config.C.DB

	switch cfg.Type {
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(cfg.Database), &gorm.Config{
			Logger: gormlogger.Default.LogMode(getLogMode(cfg.LogMode)),
		})
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.Type)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// 自动迁移
	if err := db.AutoMigrate(
		&model.User{},
		&model.Category{},
		&model.Product{},
		&model.Order{},
		&model.OrderItem{},
	); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	logger.Info("Database initialized successfully")
	return db, nil
}

func getLogMode(logMode bool) gormlogger.LogLevel {
	if logMode {
		return gormlogger.Info
	}
	return gormlogger.Silent
}

// Repository 仓库接口
type Repository struct {
	User     UserRepository
	Product  ProductRepository
	Category CategoryRepository
	Order    OrderRepository
}

// NewRepository 创建仓库实例
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		User:     NewUserRepository(db),
		Product:  NewProductRepository(db),
		Category: NewCategoryRepository(db),
		Order:    NewOrderRepository(db),
	}
}
