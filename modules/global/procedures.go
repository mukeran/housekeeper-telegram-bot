package global

import (
	. "HouseKeeperBot/common"
	. "HouseKeeperBot/modules/global/handlers"
)

func Procedures() ProcedureHandlerMap {
	return ProcedureHandlerMap{
		ProcedureManageWhitelist: ProcedureWhitelist(),
		ProcedureManageBlacklist: ProcedureBlacklist(),
		ProcedureManageAddAdmin:  ProcedureAddAdmin(),
	}
}
