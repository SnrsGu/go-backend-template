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

func (r *Repository) Create(ctx context.Context, user *User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("创建用户记录失败: %w", err)
	}
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id uint64) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).Clauses(dbresolver.Write).First(&user, "id = ?", id).Error

	if err != nil {
		return nil, fmt.Errorf("查询用户 %d 失败: %w", id, err)
	}
	return &user, nil
}

func (r *Repository) List(ctx context.Context, page int, pageSize int) ([]User, int64, error) {
	var users []User
	var total int64

	db := r.db.WithContext(ctx).Model(&User{})

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("统计用户数量失败: %w", err)
	}

	offset := (page - 1) * pageSize

	if err := db.
		Order("id DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&users).
		Error; err != nil {
		return nil, 0, fmt.Errorf("查询用户列表失败: %w", err)
	}

	return users, total, nil
}

func (r *Repository) UpdateUsername(ctx context.Context, id uint64, username string) (bool, error) {
	result := r.db.WithContext(ctx).Model(&User{}).Where("id = ?", id).Update("username", username)

	if result.Error != nil {
		return false, fmt.Errorf("修改用户 %d 失败: %w", id, result.Error)
	}

	return result.RowsAffected > 0, nil
}

func (r *Repository) Delete(ctx context.Context, id uint64) (bool, error) {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&User{})

	if result.Error != nil {
		return false, fmt.Errorf("删除用户 %d 失败: %w", id, result.Error)
	}
	return result.RowsAffected > 0, nil
}
