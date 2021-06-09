package handlers

import (
	"HouseKeeperBot/cache"
	. "HouseKeeperBot/common"
	"HouseKeeperBot/database"
	"fmt"

	tgbotapi "github.com/mukeran/telegram-bot-api"
)

const (
	tplUpdateNoAccount = `You have no account yet.`
	tplUserAuthUpdated = `Successfully updated user %v's Cookie user\_auth to %v.`
)

type paramUpdating struct {
	AccountID uint
}

func Update() CommandHandlerFunc {
	return func(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, from *tgbotapi.User) {
		resp := tgbotapi.NewMessage(msg.Chat.ID, "")
		defer MustSend(bot, &resp)
		resp.Text = "Please select an account to update Cookie user_auth:"
		buttons := generateAccountListInlineKeyboardButtons(from.ID, CallbackCccatUpdate)
		if buttons == nil {
			resp.Text = tplUpdateNoAccount
		} else {
			keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons...)
			resp.ReplyMarkup = keyboard
		}
	}
}

func generateEditUpdateList(chatID int64, messageID int, fromID int) (resp tgbotapi.EditMessageTextConfig) {
	resp = tgbotapi.NewEditMessageText(chatID, messageID, "Please select an account to update Cookie user_auth:")
	buttons := generateAccountListInlineKeyboardButtons(fromID, CallbackCccatUpdate)
	if buttons == nil {
		resp.Text = tplUpdateNoAccount
	} else {
		keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons...)
		resp.ReplyMarkup = &keyboard
	}
	return
}

func OnUpdateButtonClick() CallbackQueryHandlerFunc {
	return func(bot *tgbotapi.BotAPI, lastMsg *tgbotapi.Message, from *tgbotapi.User,
		callbackQueryID string, param string) {
		var params ParamID
		DecodeParam(param, &params)
		resp := tgbotapi.NewMessage(lastMsg.Chat.ID, "")
		defer MustSend(bot, &resp)
		defer func(resp *tgbotapi.MessageConfig) {
			QuickAnswerCallbackQuery(bot, callbackQueryID, resp.Text)
		}(&resp)
		account, err := getAccountByIDWithSecurityCheck(params.ID, from.ID)
		if err != nil {
			resp.Text = getRespText(err)
			if err == errAccountNotFound {
				MustSend(bot, generateEditUpdateList(lastMsg.Chat.ID, lastMsg.MessageID, from.ID))
			}
			return
		}
		resp.Text = fmt.Sprintf(`Please input new Cookie user_auth for account %v:`, params.ID)
		cache.RecordProcedure(lastMsg.Chat.ID, from.ID, ProcedureCccatUpdate, EncodeParam(paramUpdating{AccountID: account.ID}))
	}
}

func ProcedureUpdate() ProcedureHandlerFunc {
	return func(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, from *tgbotapi.User, param string) {
		var params paramUpdating
		DecodeParam(param, &params)
		resp := tgbotapi.NewMessage(msg.Chat.ID, "")
		resp.ParseMode = tgbotapi.ModeMarkdown
		defer MustSend(bot, &resp)
		account, err := getAccountByIDWithSecurityCheck(params.AccountID, from.ID)
		if err != nil {
			resp.Text = getRespText(err)
			if err == errAccountNotFound {
				MustSend(bot, generateEditUpdateList(msg.Chat.ID, msg.MessageID, from.ID))
			}
			return
		}
		tx := database.Db.Begin()
		defer tx.RollbackUnlessCommitted()
		account.CookieUserAuth = msg.Text
		DatabasePanicError(tx.Save(&account))
		DatabasePanicError(tx.Commit())
		resp.Text = fmt.Sprintf(tplUserAuthUpdated, params.AccountID, msg.Text)
		cache.ClearProcedure(msg.Chat.ID, from.ID)
	}
}
