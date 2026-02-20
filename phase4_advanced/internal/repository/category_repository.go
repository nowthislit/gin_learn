package repository

import (
	"gin-learn/phase4/internal/model"

	"gorm.io/gorm"
)

// CategoryRepository 分类仓库接口
type CategoryRepository interface {
	Create(category *model.Category) error
	GetByID(id uint) (*model.Category, error)
	List() ([]model.Category, error)
}

// categoryRepository 分类仓库实现
type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(category *model.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) GetByID(id uint) (*model.Category, error) {
	var category model.Category
	if err := r.db.Preload("Products").First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) List() ([]model.Category, error) {
	var categories []model.Category
	if err := r.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}
