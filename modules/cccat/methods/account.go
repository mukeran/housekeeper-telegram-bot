package methods

import (
	. "HouseKeeperBot/common"
	"HouseKeeperBot/database"
	"HouseKeeperBot/modules/cccat/models"
	"github.com/jinzhu/gorm"
	"log"
)

func GetAccountByID(accountID uint) *models.Account {
	var account models.Account
	tx := database.Db
	if v := tx.Where("id = ?", accountID).First(&account); gorm.IsRecordNotFoundError(v.Error) {
		return nil
	} else if v.Error != nil {
		log.Panic(v.Error)
	}
	return &account
}

func DeleteAccount(account *models.Account) {
	tx := database.Db.Begin()
	defer tx.RollbackUnlessCommitted()
	DatabasePanicError(tx.Delete(account))
	DatabasePanicError(tx.Commit())
}
