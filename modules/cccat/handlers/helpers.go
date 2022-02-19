package handlers

import (
	"errors"
	"fmt"
	"github.com/mukeran/housekeeper-telegram-bot/cache"
	. "github.com/mukeran/housekeeper-telegram-bot/common"
	"github.com/mukeran/housekeeper-telegram-bot/database"
	"github.com/mukeran/housekeeper-telegram-bot/modules/cccat/methods"
	"github.com/mukeran/housekeeper-telegram-bot/modules/cccat/models"
	"log"
	"time"

	tgbotapi "github.com/mukeran/telegram-bot-api"
)

var (
	errAccountNotFound  = errors.New("account not found")
	errPermissionDenied = errors.New("permission denied")
)

func generateAccountListInlineKeyboardButtons(fromID int64, callback string) (buttons [][]tgbotapi.InlineKeyboardButton) {
	var accounts []models.Account
	tx := database.Db
	if v := tx.Where("created_by = ?", fromID).
		Select("id, email, cookie_uid, has_login_credentials").Find(&accounts); v.Error != nil {
		log.Panic(v.Error)
	}
	for _, account := range accounts {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(
			func() string {
				if account.HasLoginCredentials {
					return fmt.Sprintf("%v | %v", account.ID, account.Email)
				}
				return fmt.Sprintf("%v | %v", account.ID, account.CookieUID)
			}(), cache.RecordCallback(callback, EncodeParam(ParamID{ID: account.ID})),
		)))
	}
	return
}

func getAccountByIDWithSecurityCheck(accountID uint, fromID int64) (*models.Account, error) {
	account := methods.GetAccountByID(accountID)
	if account == nil {
		return nil, errAccountNotFound
	}
	if account.CreatedBy != fromID {
		return nil, errPermissionDenied
	}
	return account, nil
}

func getRespText(err error) string {
	switch err {
	case errAccountNotFound:
		return "Bad request!"
	case errPermissionDenied:
		return "Permission denied"
	default:
		return "Unknown error"
	}
}

func stringifyTime(t time.Time) string {
	return t.UTC().Format("2006-01-02 15:04:05") + " UTC"
}
