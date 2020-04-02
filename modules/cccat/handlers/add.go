package handlers

import (
	"HouseKeeperBot/cache"
	. "HouseKeeperBot/common"
	"HouseKeeperBot/database"
	"HouseKeeperBot/modules/cccat/methods"
	"HouseKeeperBot/modules/cccat/models"
	"fmt"
	"github.com/mukeran/telegram-bot-api"
)

const (
	modeCredential          = "Using email and password"
	modeCookie              = "Using cookie"
	stepSelectMode          = "selectMode"
	stepInputEmail          = "inputEmail"
	stepInputPassword       = "inputPassword"
	stepInputCookieUid      = "inputCookieUid"
	stepInputCookieUserPwd  = "inputCookieUserPwd"
	tplSuccessAddCredential = `Successfully added account\!
*ID*: %v
*Mode*: Credential
*Email*: %v
*Password*: %v`
	tplSuccessAddCookie = `Successfully added account\!
*ID*: %v
*Mode*: Cookie
*uid*: %v
*user\_pwd*: %v`
)

type paramAdding struct {
	Step      string `json:"step"`
	Email     string `json:"username,omitempty"`
	CookieUid string `json:"cookieUid,omitempty"`
}

func generateSuccessfulAddKeyboard(account *models.Account) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Sign now",
				cache.RecordCallback(CallbackCccatSign,
					EncodeParam(paramSign{ID: account.ID})),
			),
			tgbotapi.NewInlineKeyboardButtonData(func() string {
				if account.AutoSign {
					return "Auto sign on"
				} else {
					return "Auto sign off"
				}
			}(), cache.RecordCallback(CallbackCccatAddResultToggleAutoSign,
				EncodeParam(ParamID{ID: account.ID})),
			),
			tgbotapi.NewInlineKeyboardButtonData("Delete",
				cache.RecordCallback(CallbackCccatAddResultDelete,
					EncodeParam(ParamID{ID: account.ID})),
			),
		),
	)
}

func Add() CommandHandlerFunc {
	return func(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, from *tgbotapi.User) {
		resp := tgbotapi.NewMessage(msg.Chat.ID, "")
		defer MustSend(bot, &resp)
		keyboard := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(modeCredential)),
			tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(modeCookie)),
		)
		resp.Text = "Please select a way to add your account:"
		keyboard.OneTimeKeyboard = true
		resp.ReplyMarkup = keyboard
		cache.RecordProcedure(msg.Chat.ID, from.ID, ProcedureCccatAdd, EncodeParam(paramAdding{Step: stepSelectMode}))
	}
}

func deleteAddSuccessfulMessage(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	e := tgbotapi.NewEditMessageText(msg.Chat.ID, msg.MessageID, "~"+EscapeMarkdownV2(msg.Text)+"~")
	e.ParseMode = tgbotapi.ModeMarkdownV2
	MustSend(bot, e)
}

func ProcedureAdd() ProcedureHandlerFunc {
	return func(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, from *tgbotapi.User, param string) {
		var params paramAdding
		DecodeParam(param, &params)
		resp := tgbotapi.NewMessage(msg.Chat.ID, "")
		resp.ParseMode = tgbotapi.ModeMarkdown
		defer MustSend(bot, &resp)
		switch params.Step {
		case stepSelectMode:
			resp.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
			switch msg.Text {
			case modeCredential:
				resp.Text = "Please input your CCCAT account *email*:"
				params.Step = stepInputEmail
			case modeCookie:
				resp.Text = "Please input your active CCCAT session's Cookie *uid*:"
				params.Step = stepInputCookieUid
			default:
				resp.Text = `Invalid mode\! Please start over.`
				cache.ClearProcedure(msg.Chat.ID, from.ID)
				return
			}
		case stepInputEmail:
			resp.Text = "Please input your CCCAT account *password*:"
			params.Step = stepInputPassword
			params.Email = msg.Text
		case stepInputCookieUid:
			resp.Text = "Please input your active CCCAT session's Cookie *user_pwd*:"
			params.Step = stepInputCookieUserPwd
			params.CookieUid = msg.Text
		case stepInputPassword, stepInputCookieUserPwd:
			var account models.Account
			if params.Step == stepInputPassword {
				account = models.Account{
					HasLoginCredentials: true,
					Email:               params.Email,
					Password:            msg.Text,
				}
			} else {
				account = models.Account{
					CookieUID:     params.CookieUid,
					CookieUserPwd: msg.Text,
				}
			}
			account.AutoSign = true
			account.CreatedBy = from.ID
			account.ResultChatID = msg.Chat.ID
			tx := database.Db.Begin()
			defer tx.RollbackUnlessCommitted()
			DatabasePanicError(tx.Create(&account))
			DatabasePanicError(tx.Commit())
			if params.Step == stepInputPassword {
				resp.Text = fmt.Sprintf(tplSuccessAddCredential,
					account.ID, EscapeMarkdownV2(account.Email), EscapeMarkdownV2(account.Password))
			} else {
				resp.Text = fmt.Sprintf(tplSuccessAddCookie,
					account.ID, EscapeMarkdownV2(account.CookieUID), EscapeMarkdownV2(account.CookieUserPwd))
			}
			resp.ParseMode = tgbotapi.ModeMarkdownV2
			resp.ReplyMarkup = generateSuccessfulAddKeyboard(&account)
			cache.ClearProcedure(msg.Chat.ID, from.ID)
			return
		}
		cache.RecordProcedure(msg.Chat.ID, from.ID, ProcedureCccatAdd, EncodeParam(params))
	}
}

func OnAutoSignToggle() CallbackQueryHandlerFunc {
	return func(bot *tgbotapi.BotAPI, lastMsg *tgbotapi.Message, from *tgbotapi.User,
		callbackQueryID string, param string) {
		var params ParamID
		DecodeParam(param, &params)
		account, err := getAccountByIDWithSecurityCheck(params.ID, from.ID)
		if err != nil {
			QuickSendTextMessage(bot, lastMsg.Chat.ID, getRespText(err))
			return
		}
		account.AutoSign = !account.AutoSign
		tx := database.Db.Begin()
		defer tx.RollbackUnlessCommitted()
		DatabasePanicError(tx.Save(account))
		DatabasePanicError(tx.Commit())
		editMsg := tgbotapi.NewEditMessageReplyMarkup(lastMsg.Chat.ID, lastMsg.MessageID,
			generateSuccessfulAddKeyboard(account))
		MustSend(bot, &editMsg)
		QuickAnswerCallbackQuery(bot, callbackQueryID,
			fmt.Sprintf("Successfully set account %v's auto sign to %v", account.ID, func() string {
				if account.AutoSign {
					return "on"
				} else {
					return "off"
				}
			}()))
	}
}

func OnAddResultDeleteButtonClick() CallbackQueryHandlerFunc {
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
				deleteAddSuccessfulMessage(bot, lastMsg)
			}
			return
		}
		methods.DeleteAccount(account)
		resp.Text = fmt.Sprintf("Successfully deleted account %v", account.ID)
		deleteAddSuccessfulMessage(bot, lastMsg)
	}
}
