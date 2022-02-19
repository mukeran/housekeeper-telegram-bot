package handlers

import (
	. "github.com/mukeran/housekeeper-telegram-bot/common"
	"github.com/mukeran/housekeeper-telegram-bot/modules/global/methods"
	"github.com/mukeran/telegram-bot-api"
)

const (
	tplHelp = `HouseKeeper Telegram Bot Help
Global:
/start \- Obtain bot's main menu \(not functional yet\)
/help \- This help information

Module CCCAT:
/cccat\_add \- Add a CCCAT account
/cccat\_del \- Delete a CCCAT account
/cccat\_list \- List all of your CCCAT account
/cccat\_sign \- Start a CCCAT sign procedure
`
	tplHelpForAdmin = tplHelp + `
For administrator:
/manage \- Manage this bot
`
)

func Help() CommandHandlerFunc {
	return func(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, from *tgbotapi.User) {
		resp := tgbotapi.NewMessage(msg.Chat.ID, "")
		if methods.IsAdmin(from.ID) {
			resp.Text = tplHelpForAdmin
		} else {
			resp.Text = tplHelp
		}
		resp.ParseMode = tgbotapi.ModeMarkdownV2
		MustSend(bot, resp)
	}
}
