package global

import (
	. "HouseKeeperBot/common"
	. "HouseKeeperBot/modules/global/handlers"
)

func Commands() CommandHandlerMap {
	return CommandHandlerMap{
		CommandStart:  Start(),
		CommandManage: Manage(),
		CommandHelp:   Help(),
	}
}
