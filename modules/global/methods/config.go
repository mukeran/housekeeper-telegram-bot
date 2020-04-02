package methods

import (
	. "HouseKeeperBot/common"
	"HouseKeeperBot/database"
	"HouseKeeperBot/modules/global/models"
	"github.com/jinzhu/gorm"
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
