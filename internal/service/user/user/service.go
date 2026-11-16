package user

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type userRepository interface {
	GetByID(ctx context.Context, id uint64) (*User, error)
}

type Service struct {
	repository userRepository
}

func NewService(repository userRepository) *Service {
	return &Service{repository: repository}
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
