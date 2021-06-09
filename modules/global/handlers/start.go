package handlers

import (
	. "HouseKeeperBot/common"
	"HouseKeeperBot/modules/global/methods"
	"fmt"
	"os"
	"strconv"

	tgbotapi "github.com/mukeran/telegram-bot-api"
)

func Start() CommandHandlerFunc {
	return func(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, from *tgbotapi.User) {
		adminId, _ := strconv.Atoi(os.Getenv("ADMIN_ID"))
		if adminId != 0 && msg.From.ID == adminId {
			methods.SetAdmin(msg.From.ID, true)
		}
		QuickSendTextMessage(bot, msg.Chat.ID, fmt.Sprintf("Welcome to HouseKeeperBot! Your Telegram ID is %v", msg.From.ID))
	}
}
