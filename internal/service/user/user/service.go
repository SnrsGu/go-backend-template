package user

import (
	"context"
	"errors"
	"fmt"
	"go-backend-template/internal/platform/response"
	"strings"
	"unicode/utf8"

	"gorm.io/gorm"
)

type userRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uint64) (*User, error)
	List(ctx context.Context, page, pageSize int) ([]User, int64, error)
	UpdateUsername(ctx context.Context, id uint64, username string) (bool, error)
	Delete(ctx context.Context, id uint64) (bool, error)
}

type Service struct {
	repository userRepository
}

func NewService(repository userRepository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Create(ctx context.Context, req CreateUserRequest) (*UserResponse, error) {
	username, err := validateUsername(req.Username)
	if err != nil {
		return nil, err
	}

	user := &User{
		Username: username,
	}

	if err = s.repository.Create(ctx, user); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, CodeUsernameExists
		}
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	result := &UserResponse{
		ID:       user.ID,
		Username: user.Username,
	}
	return result, nil
}

func (s *Service) GetByID(ctx context.Context, id uint64) (*UserResponse, error) {
	user, err := s.repository.GetByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, CodeUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("获取用户 %d 失败: %w", id, err)
	}
	result := &UserResponse{
		ID:       user.ID,
		Username: user.Username,
	}
	return result, nil
}

func (s *Service) List(ctx context.Context, page int, pageSize int) (*UserListResponse, error) {
	users, total, err := s.repository.List(ctx, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("获取用户列表失败: %w", err)
	}

	items := make([]UserResponse, 0, len(users))

	for i := range users {
		items = append(items, UserResponse{ID: users[i].ID, Username: users[i].Username})
	}

	return &UserListResponse{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (s *Service) UpdateUsername(ctx context.Context, id uint64, req UpdateUserRequest) error {
	username, err := validateUsername(req.Username)
	if err != nil {
		return err
	}

	found, err := s.repository.UpdateUsername(ctx, id, username)

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return CodeUsernameExists
	}

	if err != nil {
		return fmt.Errorf("修改用户 %d 失败: %w", id, err)
	}

	if !found {
		return CodeUserNotFound
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, id uint64) error {
	deleted, err := s.repository.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("删除用户 %d 失败: %w", id, err)
	}

	if !deleted {
		return CodeUserNotFound
	}

	return nil
}

func validateUsername(username string) (string, error) {
	username = strings.TrimSpace(username)

	length := utf8.RuneCountInString(username)

	if length < 3 || length > 64 {
		return "", response.CodeInvalidArgument
	}

	return username, nil
}
