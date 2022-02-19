package methods

import (
	. "github.com/mukeran/housekeeper-telegram-bot/common"
	"github.com/mukeran/housekeeper-telegram-bot/database"
	"github.com/mukeran/housekeeper-telegram-bot/modules/global/models"
)

func GetUserByTelegramUserID(telegramUserID int64) *models.User {
	tx := database.Db.Begin()
	defer tx.RollbackUnlessCommitted()
	var user models.User
	DatabasePanicError(tx.FirstOrCreate(&user, models.User{TelegramUserID: telegramUserID}))
	DatabasePanicError(tx.Commit())
	return &user
}

func IsAdmin(telegramUserID int64) bool {
	tx := database.Db
	var count uint
	DatabasePanicError(tx.Table(models.TableUser).
		Where("telegram_user_id = ? and is_admin = 1", telegramUserID).Count(&count))
	return count != 0
}

func SetAdmin(telegramUserID int64, isAdmin bool) {
	user := GetUserByTelegramUserID(telegramUserID)
	tx := database.Db.Begin()
	defer tx.RollbackUnlessCommitted()
	user.IsAdmin = isAdmin
	DatabasePanicError(tx.Save(user))
	DatabasePanicError(tx.Commit())
}

func AddAdmin(telegramUserID int64) {
	SetAdmin(telegramUserID, true)
}

func DeleteAdmin(telegramUserID int64) {
	SetAdmin(telegramUserID, false)
}

func IsWhitelisted(telegramUserID int64) bool {
	var count uint
	tx := database.Db
	DatabasePanicError(tx.Table(models.TableUser).
		Where("telegram_user_id = ? and (is_whitelisted = 1 or is_admin = 1)", telegramUserID).Count(&count))
	return count != 0
}

func IsBlacklisted(telegramUserID int64) bool {
	var count uint
	tx := database.Db
	DatabasePanicError(tx.Table(models.TableUser).
		Where("telegram_user_id = ? and is_blacklisted = 1", telegramUserID).Count(&count))
	return count != 0
}
