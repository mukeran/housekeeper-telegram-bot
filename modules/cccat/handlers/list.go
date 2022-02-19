package handlers

import (
	"fmt"
	"github.com/mukeran/housekeeper-telegram-bot/cache"
	. "github.com/mukeran/housekeeper-telegram-bot/common"
	"github.com/mukeran/housekeeper-telegram-bot/database"
	"github.com/mukeran/housekeeper-telegram-bot/modules/cccat/methods"
	"github.com/mukeran/housekeeper-telegram-bot/modules/cccat/models"
	"log"

	"github.com/jinzhu/gorm"
	tgbotapi "github.com/mukeran/telegram-bot-api"
)

const (
	tplListNoAccount           = `You haven't add any account. Add one first /cccat_add.`
	tplManageAccountCredential = `Account %v
*Mode*: Credential
*Email*: %v
*Auto sign*: %v
*Added at*: %v`
	tplManageAccountCookie = `Account %v
*Mode*: Cookie
*CookieUid*: %v
*Auto sign*: %v
*Added at*: %v`
	tplQueryRemainingTransferSucceeded = `Account %v remains transfer %v GB.`
	tplQueryRemainingTransferFailed    = `Failed to query remaining transfer for account %v.`
	tplNoSuccessfulSignResult          = `Account %v has no successful sign result yet.`
	tplLastSuccessfulSignResult        = `Account %v's last successful sign result:
Time: %v
Got transfer: %v MB`
)

type paramManageUpdating struct {
	LastMsgID int
	AccountID uint
}

func generateListMainMenu(chatID int64, fromID int64) (resp tgbotapi.MessageConfig) {
	resp = tgbotapi.NewMessage(chatID, "Please select an account to manage:")
	buttons := generateAccountListInlineKeyboardButtons(fromID, CallbackCccatList)
	if buttons == nil {
		resp.Text = tplListNoAccount
	} else {
		resp.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(buttons...)
	}
	return
}

func generateEditList(chatID int64, messageID int, fromID int64) (resp tgbotapi.EditMessageTextConfig) {
	resp = tgbotapi.NewEditMessageText(chatID, messageID, "Please select an account to manage:")
	buttons := generateAccountListInlineKeyboardButtons(fromID, CallbackCccatList)
	if buttons == nil {
		resp.Text = tplListNoAccount
	} else {
		keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons...)
		resp.ReplyMarkup = &keyboard
	}
	return
}

func List() CommandHandlerFunc {
	return func(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, from *tgbotapi.User) {
		resp := generateListMainMenu(msg.Chat.ID, from.ID)
		MustSend(bot, &resp)
	}
}

func generateEditManage(chatID int64, messageID int, account *models.Account) (resp tgbotapi.EditMessageTextConfig) {
	resp = tgbotapi.NewEditMessageText(chatID, messageID, func() string {
		if account.HasLoginCredentials {
			return fmt.Sprintf(tplManageAccountCredential, account.ID, EscapeMarkdownV2(account.Email), func() string {
				if account.AutoSign {
					return "on"
				} else {
					return "off"
				}
			}(), EscapeMarkdownV2(stringifyTime(account.CreatedAt)))
		} else {
			return fmt.Sprintf(tplManageAccountCookie, account.ID, EscapeMarkdownV2(account.CookieUID), func() string {
				if account.AutoSign {
					return "on"
				} else {
					return "off"
				}
			}(), EscapeMarkdownV2(stringifyTime(account.CreatedAt)))
		}
	}())
	resp.ParseMode = tgbotapi.ModeMarkdownV2
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Sign now",
				cache.RecordCallback(CallbackCccatSign,
					EncodeParam(paramSign{ID: account.ID}))),
			tgbotapi.NewInlineKeyboardButtonData("Toggle auto sign",
				cache.RecordCallback(CallbackCccatManageToggleAutoSign,
					EncodeParam(ParamID{ID: account.ID}))),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Remain. transfer",
				cache.RecordCallback(CallbackCccatManageQueryRemainingTransfer,
					EncodeParam(ParamID{ID: account.ID}))),
			tgbotapi.NewInlineKeyboardButtonData("Last successful",
				cache.RecordCallback(CallbackCccatManageGetLastSuccessfulSignResult,
					EncodeParam(ParamID{ID: account.ID}))),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Update user_auth",
				cache.RecordCallback(CallbackCccatManageUpdate,
					EncodeParam(ParamID{ID: account.ID}))),
			tgbotapi.NewInlineKeyboardButtonData("Delete",
				cache.RecordCallback(CallbackCccatManageDelete,
					EncodeParam(ParamID{ID: account.ID}))),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("< Back to list",
				cache.RecordCallback(CallbackCccatManageBackToList,
					EncodeParam(ParamID{ID: account.ID}))),
		),
	)
	resp.ReplyMarkup = &keyboard
	return
}

func OnListButtonClick() CallbackQueryHandlerFunc {
	return func(bot *tgbotapi.BotAPI, lastMsg *tgbotapi.Message, from *tgbotapi.User,
		callbackQueryID string, param string) {
		var params ParamID
		DecodeParam(param, &params)
		account, err := getAccountByIDWithSecurityCheck(params.ID, from.ID)
		if err != nil {
			QuickSendTextMessage(bot, lastMsg.Chat.ID, getRespText(err))
			if err == errAccountNotFound {
				MustSend(bot, generateEditList(lastMsg.Chat.ID, lastMsg.MessageID, from.ID))
			}
			return
		}
		MustSend(bot, generateEditManage(lastMsg.Chat.ID, lastMsg.MessageID, account))
	}
}

func OnManageToggleAutoSignButtonClick() CallbackQueryHandlerFunc {
	return func(bot *tgbotapi.BotAPI, lastMsg *tgbotapi.Message, from *tgbotapi.User,
		callbackQueryID string, param string) {
		var params ParamID
		DecodeParam(param, &params)
		account, err := getAccountByIDWithSecurityCheck(params.ID, from.ID)
		if err != nil {
			QuickSendTextMessage(bot, lastMsg.Chat.ID, getRespText(err))
			if err == errAccountNotFound {
				MustSend(bot, generateEditList(lastMsg.Chat.ID, lastMsg.MessageID, from.ID))
			}
			return
		}
		account.AutoSign = !account.AutoSign
		tx := database.Db.Begin()
		defer tx.RollbackUnlessCommitted()
		DatabasePanicError(tx.Save(account))
		DatabasePanicError(tx.Commit())
		MustSend(bot, tgbotapi.NewCallback(callbackQueryID,
			fmt.Sprintf("Successfully set account %v's auto sign to %v", account.ID, func() string {
				if account.AutoSign {
					return "on"
				} else {
					return "off"
				}
			}())))
		MustSend(bot, generateEditManage(lastMsg.Chat.ID, lastMsg.MessageID, account))
	}
}

func generateRemainingTransferMessage(accountID uint, remaining float64, err error) string {
	switch err {
	case nil:
		return fmt.Sprintf(tplQueryRemainingTransferSucceeded, accountID, remaining)
	case methods.ErrWrongAccountEmailOrPassword:
		return fmt.Sprintf(tplWrongEmailOrPassword, accountID)
	case methods.ErrInvalidCookie:
		return fmt.Sprintf(tplInvalidCookie, accountID)
	default:
		log.Printf("Failed to query remaining transfer for account %v. Error: %v", accountID, err)
		return fmt.Sprintf(tplQueryRemainingTransferFailed, accountID)
	}
}

func OnManageQueryRemainingTransferButtonClick() CallbackQueryHandlerFunc {
	return func(bot *tgbotapi.BotAPI, lastMsg *tgbotapi.Message, from *tgbotapi.User,
		callbackQueryID string, param string) {
		var params ParamID
		DecodeParam(param, &params)
		resp := tgbotapi.NewMessage(lastMsg.Chat.ID, "")
		defer MustSend(bot, &resp)
		defer func(resp *tgbotapi.MessageConfig) {
			QuickAnswerCallbackQueryWithAlert(bot, callbackQueryID, resp.Text)
		}(&resp)
		account, err := getAccountByIDWithSecurityCheck(params.ID, from.ID)
		if err != nil {
			resp.Text = getRespText(err)
			if err == errAccountNotFound {
				MustSend(bot, generateEditList(lastMsg.Chat.ID, lastMsg.MessageID, from.ID))
			}
			return
		}
		QuickSendTextMessage(bot, lastMsg.Chat.ID, fmt.Sprintf(
			"Querying account %v's remaining transfer...", account.ID))
		remaining, err := methods.QueryRemainingTransfer(account)
		resp.Text = generateRemainingTransferMessage(account.ID, remaining, err)
	}
}

func OnManageGetLastSuccessfulSignResultButtonClick() CallbackQueryHandlerFunc {
	return func(bot *tgbotapi.BotAPI, lastMsg *tgbotapi.Message, from *tgbotapi.User,
		callbackQueryID string, param string) {
		var params ParamID
		DecodeParam(param, &params)
		resp := tgbotapi.NewMessage(lastMsg.Chat.ID, "")
		defer MustSend(bot, &resp)
		defer func(resp *tgbotapi.MessageConfig) {
			QuickAnswerCallbackQueryWithAlert(bot, callbackQueryID, resp.Text)
		}(&resp)
		account, err := getAccountByIDWithSecurityCheck(params.ID, from.ID)
		if err != nil {
			resp.Text = getRespText(err)
			if err == errAccountNotFound {
				MustSend(bot, generateEditList(lastMsg.Chat.ID, lastMsg.MessageID, from.ID))
			}
			return
		}
		tx := database.Db
		var signLog models.SignLog
		if v := tx.Where("account_id = ? and status = ?", account.ID, models.SignStatusSuccessful).
			Order("created_at desc").First(&signLog); gorm.IsRecordNotFoundError(v.Error) {
			resp.Text = fmt.Sprintf(tplNoSuccessfulSignResult, account.ID)
			return
		} else if v.Error != nil {
			log.Panic(err)
		}
		resp.Text = fmt.Sprintf(tplLastSuccessfulSignResult,
			account.ID, stringifyTime(signLog.CreatedAt), signLog.GotTransfer)
	}
}

func OnManageUpdateClick() CallbackQueryHandlerFunc {
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
		cache.RecordProcedure(lastMsg.Chat.ID, from.ID, ProcedureCccatManageUpdate, EncodeParam(paramManageUpdating{LastMsgID: lastMsg.MessageID, AccountID: account.ID}))
	}
}

func ProcedureManageUpdate() ProcedureHandlerFunc {
	return func(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, from *tgbotapi.User, param string) {
		var params paramManageUpdating
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
		MustSend(bot, generateEditManage(msg.Chat.ID, params.LastMsgID, account))
	}
}

func OnManageDeleteClick() CallbackQueryHandlerFunc {
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
				MustSend(bot, generateEditList(lastMsg.Chat.ID, lastMsg.MessageID, from.ID))
			}
			return
		}
		methods.DeleteAccount(account)
		resp.Text = fmt.Sprintf("Successfully deleted account %v", account.ID)
		MustSend(bot, generateEditList(lastMsg.Chat.ID, lastMsg.MessageID, from.ID))
	}
}

func OnManageBackToListButtonClick() CallbackQueryHandlerFunc {
	return func(bot *tgbotapi.BotAPI, lastMsg *tgbotapi.Message, from *tgbotapi.User,
		callbackQueryID string, param string) {
		MustSend(bot, generateEditList(lastMsg.Chat.ID, lastMsg.MessageID, from.ID))
	}
}
