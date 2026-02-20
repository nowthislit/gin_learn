package service

import (
	"errors"

	"gin-learn/phase4/internal/model"
	"gin-learn/phase4/internal/repository"
)

// UserService 用户服务接口
type UserService interface {
	Register(username, email, password string, age int) (*model.User, error)
	GetUser(id uint) (*model.User, error)
	ListUsers(page, pageSize int, keyword string) ([]model.User, int64, error)
	UpdateUser(id uint, updates map[string]interface{}) error
	DeleteUser(id uint) error
}

// userService 用户服务实现
type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Register(username, email, password string, age int) (*model.User, error) {
	// 检查用户名是否已存在
	if _, err := s.repo.GetByUsername(username); err == nil {
		return nil, errors.New("用户名已存在")
	}

	user := &model.User{
		Username: username,
		Email:    email,
		Password: password, // 实际应该加密
		Age:      age,
		Status:   1,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetUser(id uint) (*model.User, error) {
	return s.repo.GetByID(id)
}

func (s *userService) ListUsers(page, pageSize int, keyword string) ([]model.User, int64, error) {
	return s.repo.List(page, pageSize, keyword)
}

func (s *userService) UpdateUser(id uint, updates map[string]interface{}) error {
	// 不允许更新密码（简化处理）
	delete(updates, "password")

	user, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	// 更新字段
	if username, ok := updates["username"].(string); ok {
		user.Username = username
	}
	if email, ok := updates["email"].(string); ok {
		user.Email = email
	}
	if age, ok := updates["age"].(float64); ok {
		user.Age = int(age)
	}

	return s.repo.Update(user)
}

func (s *userService) DeleteUser(id uint) error {
	return s.repo.Delete(id)
}
