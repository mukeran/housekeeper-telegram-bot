package methods

import (
	"github.com/jinzhu/gorm"
	. "github.com/mukeran/housekeeper-telegram-bot/common"
	"github.com/mukeran/housekeeper-telegram-bot/database"
	"github.com/mukeran/housekeeper-telegram-bot/modules/global/models"
	"log"
)

func SetConfig(field, value string) {
	tx := database.Db.Begin()
	defer tx.RollbackUnlessCommitted()
	var config models.Config
	DatabasePanicError(tx.FirstOrCreate(&config, models.Config{Field: field}))
	config.Value = value
	DatabasePanicError(tx.Save(&config))
	DatabasePanicError(tx.Commit())
}

func GetConfig(field string) string {
	tx := database.Db
	var config models.Config
	if v := tx.Where("field = ?", field).First(&config); gorm.IsRecordNotFoundError(v.Error) {
		return ""
	} else if v.Error != nil {
		log.Panic(v.Error)
	}
	return config.Value
}
