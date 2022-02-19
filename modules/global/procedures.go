package global

import (
	. "github.com/mukeran/housekeeper-telegram-bot/common"
	. "github.com/mukeran/housekeeper-telegram-bot/modules/global/handlers"
)

func Procedures() ProcedureHandlerMap {
	return ProcedureHandlerMap{
		ProcedureManageWhitelist: ProcedureWhitelist(),
		ProcedureManageBlacklist: ProcedureBlacklist(),
		ProcedureManageAddAdmin:  ProcedureAddAdmin(),
	}
}
