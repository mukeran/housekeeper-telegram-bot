package cccat

import (
	. "github.com/mukeran/housekeeper-telegram-bot/common"
	. "github.com/mukeran/housekeeper-telegram-bot/modules/cccat/handlers"
)

func Procedures() ProcedureHandlerMap {
	return ProcedureHandlerMap{
		ProcedureCccatAdd:          ProcedureAdd(),
		ProcedureCccatUpdate:       ProcedureUpdate(),
		ProcedureCccatManageUpdate: ProcedureManageUpdate(),
	}
}
