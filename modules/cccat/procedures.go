package cccat

import (
	. "HouseKeeperBot/common"
	. "HouseKeeperBot/modules/cccat/handlers"
)

func Procedures() ProcedureHandlerMap {
	return ProcedureHandlerMap{
		ProcedureCccatAdd: ProcedureAdd(),
	}
}
