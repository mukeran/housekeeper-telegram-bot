package cccat

import (
	. "HouseKeeperBot/common"
	. "HouseKeeperBot/modules/cccat/handlers"
)

func CallbackQueries() CallbackQueryHandlerMap {
	return CallbackQueryHandlerMap{
		CallbackCccatSign:                              OnSignButtonClick(),
		CallbackCccatDel:                               OnDelButtonClick(),
		CallbackCccatAddResultToggleAutoSign:           OnAutoSignToggle(),
		CallbackCccatAddResultDelete:                   OnAddResultDeleteButtonClick(),
		CallbackCccatList:                              OnListButtonClick(),
		CallbackCccatManageToggleAutoSign:              OnManageToggleAutoSignButtonClick(),
		CallbackCccatManageQueryRemainingTransfer:      OnManageQueryRemainingTransferButtonClick(),
		CallbackCccatManageGetLastSuccessfulSignResult: OnManageGetLastSuccessfulSignResultButtonClick(),
		CallbackCccatManageDelete:                      OnManageDeleteClick(),
		CallbackCccatManageBackToList:                  OnManageBackToListButtonClick(),
	}
}
