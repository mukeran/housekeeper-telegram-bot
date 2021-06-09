package cccat

import (
	. "HouseKeeperBot/common"
	"HouseKeeperBot/database"
	"HouseKeeperBot/modules/cccat/methods"
	"HouseKeeperBot/modules/cccat/models"
	"HouseKeeperBot/schedule"
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/mukeran/telegram-bot-api"
)

const (
	tplScheduledSignSucceeded        = `[Scheduled] Successfully signed account %v (%v). Got %v MB.`
	tplScheduledSignSigned           = `[Scheduled] Account %v (%v) has been signed today.`
	tplScheduledWrongEmailOrPassword = `[Scheduled] CCCAT reports "wrong email or password" for account %v (%v). Please check your settings.`
	tplScheduledInvalidCookie        = `[Scheduled] Account %v (%v) has an invalid cookie. Go to cccat.io to re-login your account and update your user_auth using /cccat_update.`
	tplScheduledSignFailed           = `[Scheduled] Failed to sign account %v (%v).`
)

func generateAccountName(account *models.Account) string {
	if account.HasLoginCredentials {
		return fmt.Sprintf("Email: %v", account.Email)
	} else {
		return fmt.Sprintf("Uid: %v", account.CookieUID)
	}
}

func generateSignReport(account *models.Account, got uint, err error) string {
	accountName := generateAccountName(account)
	switch err {
	case nil:
		return fmt.Sprintf(tplScheduledSignSucceeded, account.ID, accountName, got)
	case methods.ErrSigned:
		return fmt.Sprintf(tplScheduledSignSigned, account.ID, accountName)
	case methods.ErrWrongAccountEmailOrPassword:
		return fmt.Sprintf(tplScheduledWrongEmailOrPassword, account.ID, accountName)
	case methods.ErrInvalidCookie:
		return fmt.Sprintf(tplScheduledInvalidCookie, account.ID, accountName)
	default:
		log.Printf("[Scheduled] Failed to sign account %v. Error: %v", account.ID, err)
		return fmt.Sprintf(tplScheduledSignFailed, account.ID, accountName)
	}
}

func scheduleSign() schedule.TaskFunc {
	return func(bot *tgbotapi.BotAPI, param string) bool {
		var accounts []models.Account
		tx := database.Db
		DatabasePanicError(tx.Find(&accounts, models.Account{AutoSign: true}))
		for _, account := range accounts {
			got, err := methods.Sign(&account)
			QuickSendTextMessage(bot, account.ResultChatID, generateSignReport(&account, got, err))
		}
		return true
	}
}

func RegisterSchedules() {
	now := time.Now()
	schedule.RegisterDurationTask(scheduleSign(), "",
		time.Date(now.Year(), now.Month(), now.Day(), 0, 30, 0, 0, time.UTC),
		time.Hour*24)
}
