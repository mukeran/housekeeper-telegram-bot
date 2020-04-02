package models

import "time"

const (
	TableCCCATSignLog    = "cccat_sign_log"
	SignStatusSuccessful = iota
	SignStatusSigned
	SignStatusFailed
)

type SignLog struct {
	ID          uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Account     Account `foreignkey:"AccountID"`
	AccountID   uint
	Status      uint
	GotTransfer uint
	Raw         string
}

func (*SignLog) TableName() string {
	return TableCCCATSignLog
}
