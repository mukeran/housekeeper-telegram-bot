package models

import "time"

const (
	TableConfig         = "config"
	ConfigWhitelistMode = "global.whitelistMode"
)

type Config struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	Field     string
	Value     string
}

func (*Config) TableName() string {
	return TableConfig
}
