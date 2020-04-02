package main

import (
	. "HouseKeeperBot/common"
	"HouseKeeperBot/modules/global/methods"
	globalModels "HouseKeeperBot/modules/global/models"
	tgbotapi "github.com/mukeran/telegram-bot-api"
)

func auth(bot *tgbotapi.BotAPI, chatID int64, telegramUserID int) bool {
	if methods.IsBlacklisted(telegramUserID) {
		return false
	}
	whitelistMode := methods.GetConfig(globalModels.ConfigWhitelistMode)
	if whitelistMode == "on" {
		if !methods.IsWhitelisted(telegramUserID) {
			QuickSendTextMessage(bot, chatID, "You are not on the whitelist!")
			return false
		}
	}
	return true
}
