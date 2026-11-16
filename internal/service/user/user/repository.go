package user

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetByID(ctx context.Context, id uint64) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).Clauses(dbresolver.Write).First(&user, "id = ?", id).Error
	if err != nil {
		return nil, fmt.Errorf("查询用户 %d 失败: %w", id, err)
	}
	return &user, nil
}
