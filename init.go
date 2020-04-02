package main

import (
	. "HouseKeeperBot/common"
	"HouseKeeperBot/database"
	cccatModels "HouseKeeperBot/modules/cccat/models"
	"HouseKeeperBot/modules/global/methods"
	globalModels "HouseKeeperBot/modules/global/models"
)

func initializeDatabase() {
	tx := database.Db.Begin()
	defer tx.RollbackUnlessCommitted()
	DatabasePanicError(tx.AutoMigrate(&globalModels.User{}, &globalModels.Config{}))
	DatabasePanicError(tx.AutoMigrate(&cccatModels.Account{}, &cccatModels.SignLog{}))
	DatabasePanicError(tx.Commit())
	methods.SetConfig(globalModels.ConfigWhitelistMode, "off")
}
