package cccat

import (
	. "github.com/mukeran/housekeeper-telegram-bot/common"
	. "github.com/mukeran/housekeeper-telegram-bot/modules/cccat/handlers"
)

func CallbackQueries() CallbackQueryHandlerMap {
	return CallbackQueryHandlerMap{
		CallbackCccatSign:                              OnSignButtonClick(),
		CallbackCccatUpdate:                            OnUpdateButtonClick(),
		CallbackCccatDel:                               OnDelButtonClick(),
		CallbackCccatAddResultToggleAutoSign:           OnAutoSignToggle(),
		CallbackCccatAddResultDelete:                   OnAddResultDeleteButtonClick(),
		CallbackCccatList:                              OnListButtonClick(),
		CallbackCccatManageToggleAutoSign:              OnManageToggleAutoSignButtonClick(),
		CallbackCccatManageQueryRemainingTransfer:      OnManageQueryRemainingTransferButtonClick(),
		CallbackCccatManageGetLastSuccessfulSignResult: OnManageGetLastSuccessfulSignResultButtonClick(),
		CallbackCccatManageUpdate:                      OnManageUpdateClick(),
		CallbackCccatManageDelete:                      OnManageDeleteClick(),
		CallbackCccatManageBackToList:                  OnManageBackToListButtonClick(),
	}
}
