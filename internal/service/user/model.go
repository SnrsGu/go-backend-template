package user

import "time"

type User struct {
	ID        uint64    `gorm:"column:id;primary_key;auto_increment"`
	Username  string    `gorm:"column:username"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (User) TableName() string {
	return "user"
}
