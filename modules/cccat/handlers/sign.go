package handlers

import (
	"HouseKeeperBot/cache"
	. "HouseKeeperBot/common"
	"HouseKeeperBot/database"
	"HouseKeeperBot/modules/cccat/methods"
	"HouseKeeperBot/modules/cccat/models"
	"fmt"
	"github.com/mukeran/telegram-bot-api"
	"log"
)

const (
	tplSignSucceeded        = `Successfully signed account %v. Got %v MB.`
	tplSignSigned           = `Account %v has been signed today.`
	tplWrongEmailOrPassword = `CCCAT reports "wrong email or password" for account %v. Please check your settings.`
	tplInvalidCookie        = `Account %v has an invalid cookie.`
	tplSignFailed           = `Failed to sign account %v.`
)

func generateSignReport(accountID, got uint, err error) string {
	switch err {
	case nil:
		return fmt.Sprintf(tplSignSucceeded, accountID, got)
	case methods.ErrSigned:
		return fmt.Sprintf(tplSignSigned, accountID)
	case methods.ErrWrongAccountEmailOrPassword:
		return fmt.Sprintf(tplWrongEmailOrPassword, accountID)
	case methods.ErrInvalidCookie:
		return fmt.Sprintf(tplInvalidCookie, accountID)
	default:
		log.Printf("Failed to sign account %v. Error: %v", accountID, err)
		return fmt.Sprintf(tplSignFailed, accountID)
	}
}

type paramSign struct {
	All       bool `json:"all,omitempty"`
	AccountID uint `json:"accountID,omitempty"`
}

func Sign() CommandHandlerFunc {
	return func(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, from *tgbotapi.User) {
		resp := tgbotapi.NewMessage(msg.Chat.ID, "")
		defer MustSend(bot, &resp)
		resp.Text = "Please select an account (or all account) to sign"
		var buttons [][]tgbotapi.InlineKeyboardButton
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("All account",
				cache.RecordCallback(CallbackCccatSign, EncodeParam(paramSign{All: true})),
			),
		))
		buttons = append(buttons, generateAccountListInlineKeyboardButtons(from.ID, CallbackCccatSign)...)
		keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons...)
		resp.ReplyMarkup = keyboard
	}
}

func OnSignButtonClick() CallbackQueryHandlerFunc {
	return func(bot *tgbotapi.BotAPI, lastMsg *tgbotapi.Message, from *tgbotapi.User,
		callbackQueryID string, param string) {
		var params paramSign
		DecodeParam(param, &params)
		if params.All {
			QuickAnswerCallbackQuery(bot, callbackQueryID, "Signing procedure started")
			var accountIDs []uint
			tx := database.Db
			DatabasePanicError(tx.Table(models.TableCCCATAccount).Where("created_by = ?", from.ID).
				Select("id").Pluck("id", &accountIDs))
			for _, accountID := range accountIDs {
				got, err := methods.SignWithAccountID(accountID)
				QuickSendTextMessage(bot, lastMsg.Chat.ID, generateSignReport(accountID, got, err))
			}
			QuickSendTextMessage(bot, lastMsg.Chat.ID, "Sign completed")
		} else {
			resp := tgbotapi.NewMessage(lastMsg.Chat.ID, "")
			defer MustSend(bot, &resp)
			defer func(resp *tgbotapi.MessageConfig) {
				QuickAnswerCallbackQueryWithAlert(bot, callbackQueryID, resp.Text)
			}(&resp)
			account, err := getAccountByIDWithSecurityCheck(params.AccountID, from.ID)
			if err != nil {
				resp.Text = getRespText(err)
				return
			}
			QuickSendTextMessage(bot, lastMsg.Chat.ID, fmt.Sprintf("Signing account %v...", account.ID))
			got, err := methods.Sign(account)
			resp.Text = generateSignReport(account.ID, got, err)
		}
	}
}
