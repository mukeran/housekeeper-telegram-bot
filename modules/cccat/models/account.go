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
	CookieUserAuth      string
	CookieUserPwd       string
	HasLoginCredentials bool
	Email               string
	Password            string
	AutoSign            bool
	CreatedBy           int64
	ResultChatID        int64
}

func (*Account) TableName() string {
	return TableCCCATAccount
}
