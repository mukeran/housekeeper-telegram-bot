package models

import "time"

const (
	TableUser = "user"
)

type User struct {
	ID             uint
	CreatedAt      time.Time
	UpdatedAt      time.Time
	TelegramUserID int64
	IsAdmin        bool
	IsWhitelisted  bool
	IsBlacklisted  bool
}

func (*User) TableName() string {
	return TableUser
}
