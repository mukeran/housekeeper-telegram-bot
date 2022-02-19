package main

import (
	. "github.com/mukeran/housekeeper-telegram-bot/common"
	"github.com/mukeran/housekeeper-telegram-bot/database"
	cccatModels "github.com/mukeran/housekeeper-telegram-bot/modules/cccat/models"
	"github.com/mukeran/housekeeper-telegram-bot/modules/global/methods"
	globalModels "github.com/mukeran/housekeeper-telegram-bot/modules/global/models"
)

func initializeDatabase() {
	tx := database.Db.Begin()
	defer tx.RollbackUnlessCommitted()
	DatabasePanicError(tx.AutoMigrate(&globalModels.User{}, &globalModels.Config{}))
	DatabasePanicError(tx.AutoMigrate(&cccatModels.Account{}, &cccatModels.SignLog{}))
	DatabasePanicError(tx.Commit())
	methods.SetConfig(globalModels.ConfigWhitelistMode, "off")
}
