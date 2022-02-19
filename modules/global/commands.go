package global

import (
	. "github.com/mukeran/housekeeper-telegram-bot/common"
	. "github.com/mukeran/housekeeper-telegram-bot/modules/global/handlers"
)

func Commands() CommandHandlerMap {
	return CommandHandlerMap{
		CommandStart:  Start(),
		CommandManage: Manage(),
		CommandHelp:   Help(),
	}
}
