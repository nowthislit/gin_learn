package repository

import (
	"gin-learn/phase4/internal/model"

	"gorm.io/gorm"
)

// ProductRepository 产品仓库接口
type ProductRepository interface {
	Create(product *model.Product) error
	GetByID(id uint) (*model.Product, error)
	List(page, pageSize int, categoryID uint, keyword string) ([]model.Product, int64, error)
	Update(product *model.Product) error
	Delete(id uint) error
	Search(minPrice, maxPrice float64, categoryID uint, keyword, sortBy, sortOrder string, page, pageSize int) ([]model.Product, int64, error)
}

// productRepository 产品仓库实现
type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *model.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) GetByID(id uint) (*model.Product, error) {
	var product model.Product
	if err := r.db.Preload("Category").First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) List(page, pageSize int, categoryID uint, keyword string) ([]model.Product, int64, error) {
	query := r.db.Model(&model.Product{})

	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}

	if keyword != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var products []model.Product
	offset := (page - 1) * pageSize
	if err := query.Preload("Category").Offset(offset).Limit(pageSize).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepository) Update(product *model.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(id uint) error {
	return r.db.Delete(&model.Product{}, id).Error
}

func (r *productRepository) Search(minPrice, maxPrice float64, categoryID uint, keyword, sortBy, sortOrder string, page, pageSize int) ([]model.Product, int64, error) {
	query := r.db.Model(&model.Product{})

	if keyword != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	if minPrice > 0 {
		query = query.Where("price >= ?", minPrice)
	}
	if maxPrice > 0 {
		query = query.Where("price <= ?", maxPrice)
	}
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}

	orderStr := sortBy + " " + sortOrder
	query = query.Order(orderStr)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	var products []model.Product
	if err := query.Preload("Category").Offset(offset).Limit(pageSize).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}
