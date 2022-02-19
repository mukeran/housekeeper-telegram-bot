package global

import (
	. "github.com/mukeran/housekeeper-telegram-bot/common"
	. "github.com/mukeran/housekeeper-telegram-bot/modules/global/handlers"
)

func CallbackQueries() CallbackQueryHandlerMap {
	return CallbackQueryHandlerMap{
		CallbackBackToManage:              OnBackToManageButtonClick(),
		CallbackManageAuth:                OnManageAuthButtonClick(),
		CallbackManageAdmin:               OnManageAdminButtonClick(),
		CallbackManageToggleWhitelistMode: OnManageToggleWhitelistModeButtonClick(),
		CallbackManageWhitelist:           OnManageWhitelistButtonClick(),
		CallbackManageToggleIsWhitelisted: OnManageToggleIsWhitelistedButtonClick(),
		CallbackManageBlacklist:           OnManageBlacklistButtonClick(),
		CallbackManageToggleIsBlacklisted: OnManageToggleIsBlacklistedButtonClick(),
		CallbackManageListAdmins:          OnManageListAdminsButtonClick(),
		CallbackManageDeleteAdmin:         OnManageDeleteAdminButtonClick(),
		CallbackManageAddAdmin:            OnManageAddAdminButtonClick(),
		CallbackManageRevertAddAdmin:      OnManageRevertAddAdminButtonClick(),
	}
}
