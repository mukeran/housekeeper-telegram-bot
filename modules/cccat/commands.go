package cccat

import (
	. "HouseKeeperBot/common"
	. "HouseKeeperBot/modules/cccat/handlers"
)

func Commands() CommandHandlerMap {
	return CommandHandlerMap{
		CommandCccatSign: Sign(),
		CommandCccatAdd:  Add(),
		CommandCccatList: List(),
		CommandCccatDel:  Del(),
	}
}
