package methods

import (
	. "HouseKeeperBot/common"
	"HouseKeeperBot/database"
	"HouseKeeperBot/modules/global/models"
)

func GetUserByTelegramUserID(telegramUserID int) *models.User {
	tx := database.Db.Begin()
	defer tx.RollbackUnlessCommitted()
	var user models.User
	DatabasePanicError(tx.FirstOrCreate(&user, models.User{TelegramUserID: telegramUserID}))
	DatabasePanicError(tx.Commit())
	return &user
}

func IsAdmin(telegramUserID int) bool {
	tx := database.Db
	var count uint
	DatabasePanicError(tx.Table(models.TableUser).
		Where("telegram_user_id = ? and is_admin = 1", telegramUserID).Count(&count))
	return count != 0
}

func SetAdmin(telegramUserID int, isAdmin bool) {
	user := GetUserByTelegramUserID(telegramUserID)
	tx := database.Db.Begin()
	defer tx.RollbackUnlessCommitted()
	user.IsAdmin = isAdmin
	DatabasePanicError(tx.Save(user))
	DatabasePanicError(tx.Commit())
}

func AddAdmin(telegramUserID int) {
	SetAdmin(telegramUserID, true)
}

func DeleteAdmin(telegramUserID int) {
	SetAdmin(telegramUserID, false)
}

func IsWhitelisted(telegramUserID int) bool {
	var count uint
	tx := database.Db
	DatabasePanicError(tx.Table(models.TableUser).
		Where("telegram_user_id = ? and is_whitelisted = 1", telegramUserID).Count(&count))
	return count != 0
}

func IsBlacklisted(telegramUserID int) bool {
	var count uint
	tx := database.Db
	DatabasePanicError(tx.Table(models.TableUser).
		Where("telegram_user_id = ? and is_blacklisted = 1", telegramUserID).Count(&count))
	return count != 0
}
