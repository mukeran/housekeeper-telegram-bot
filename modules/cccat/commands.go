package cccat

import (
	. "github.com/mukeran/housekeeper-telegram-bot/common"
	. "github.com/mukeran/housekeeper-telegram-bot/modules/cccat/handlers"
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
