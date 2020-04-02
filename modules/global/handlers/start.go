package handlers

import (
	. "HouseKeeperBot/common"
	"github.com/mukeran/telegram-bot-api"
)

func Start() CommandHandlerFunc {
	return func(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, from *tgbotapi.User) {
		QuickSendTextMessage(bot, msg.Chat.ID, "Welcome to HouseKeeperBot!")
	}
}
