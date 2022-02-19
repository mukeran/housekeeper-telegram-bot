package handlers

import (
	"fmt"
	. "github.com/mukeran/housekeeper-telegram-bot/common"
	"github.com/mukeran/housekeeper-telegram-bot/modules/global/methods"
	"os"
	"strconv"

	tgbotapi "github.com/mukeran/telegram-bot-api"
)

func Start() CommandHandlerFunc {
	return func(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, from *tgbotapi.User) {
		adminId, _ := strconv.ParseInt(os.Getenv("ADMIN_ID"), 10, 64)
		if adminId != 0 && msg.From.ID == adminId {
			methods.SetAdmin(msg.From.ID, true)
		}
		QuickSendTextMessage(bot, msg.Chat.ID, fmt.Sprintf("Welcome to HouseKeeper Telegram Bot! Your Telegram ID is %v", msg.From.ID))
	}
}
