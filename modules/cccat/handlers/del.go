package handlers

import (
	"fmt"
	. "github.com/mukeran/housekeeper-telegram-bot/common"
	"github.com/mukeran/housekeeper-telegram-bot/modules/cccat/methods"
	"github.com/mukeran/telegram-bot-api"
)

const (
	tplDelNoAccount = `You have no account yet.`
)

func Del() CommandHandlerFunc {
	return func(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, from *tgbotapi.User) {
		resp := tgbotapi.NewMessage(msg.Chat.ID, "")
		defer MustSend(bot, &resp)
		resp.Text = "Please select an account to delete:"
		buttons := generateAccountListInlineKeyboardButtons(from.ID, CallbackCccatDel)
		if buttons == nil {
			resp.Text = tplDelNoAccount
		} else {
			keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons...)
			resp.ReplyMarkup = keyboard
		}
	}
}

func generateEditDelList(chatID int64, messageID int, fromID int64) (resp tgbotapi.EditMessageTextConfig) {
	resp = tgbotapi.NewEditMessageText(chatID, messageID, "Please select an account to delete:")
	buttons := generateAccountListInlineKeyboardButtons(fromID, CallbackCccatDel)
	if buttons == nil {
		resp.Text = tplDelNoAccount
	} else {
		keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons...)
		resp.ReplyMarkup = &keyboard
	}
	return
}

func OnDelButtonClick() CallbackQueryHandlerFunc {
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
				MustSend(bot, generateEditDelList(lastMsg.Chat.ID, lastMsg.MessageID, from.ID))
			}
			return
		}
		methods.DeleteAccount(account)
		resp.Text = fmt.Sprintf("Successfully deleted account %v", account.ID)
		MustSend(bot, generateEditDelList(lastMsg.Chat.ID, lastMsg.MessageID, from.ID))
	}
}
