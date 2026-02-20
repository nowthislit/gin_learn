package service

import (
	"gin-learn/phase4/internal/model"
	"gin-learn/phase4/internal/repository"
)

// CategoryService 分类服务接口
type CategoryService interface {
	CreateCategory(name, description string) (*model.Category, error)
	GetCategory(id uint) (*model.Category, error)
	ListCategories() ([]model.Category, error)
}

// categoryService 分类服务实现
type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) CreateCategory(name, description string) (*model.Category, error) {
	category := &model.Category{
		Name:        name,
		Description: description,
	}

	if err := s.repo.Create(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) GetCategory(id uint) (*model.Category, error) {
	return s.repo.GetByID(id)
}

func (s *categoryService) ListCategories() ([]model.Category, error) {
	return s.repo.List()
}
