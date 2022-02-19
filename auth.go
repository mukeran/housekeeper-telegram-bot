package main

import (
	. "github.com/mukeran/housekeeper-telegram-bot/common"
	"github.com/mukeran/housekeeper-telegram-bot/modules/global/methods"
	globalModels "github.com/mukeran/housekeeper-telegram-bot/modules/global/models"
	tgbotapi "github.com/mukeran/telegram-bot-api"
)

func auth(bot *tgbotapi.BotAPI, chatID int64, telegramUserID int64) bool {
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
