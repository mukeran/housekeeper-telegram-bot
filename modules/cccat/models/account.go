package models

import (
	"time"
)

const (
	TableCCCATAccount = "cccat_account"
)

type Account struct {
	ID                  uint
	CreatedAt           time.Time
	UpdatedAt           time.Time
	CookieUID           string
	CookieUserPwd       string
	HasLoginCredentials bool
	Email               string
	Password            string
	AutoSign            bool
	CreatedBy           int
	ResultChatID        int64
}

func (*Account) TableName() string {
	return TableCCCATAccount
}
