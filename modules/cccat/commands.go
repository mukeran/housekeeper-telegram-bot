package cccat

import (
	. "HouseKeeperBot/common"
	. "HouseKeeperBot/modules/cccat/handlers"
)

func Commands() CommandHandlerMap {
	return CommandHandlerMap{
		CommandCccatSign:   Sign(),
		CommandCccatAdd:    Add(),
		CommandCccatUpdate: Update(),
		CommandCccatList:   List(),
		CommandCccatDel:    Del(),
	}
}
