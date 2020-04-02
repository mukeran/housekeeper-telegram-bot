package handlers

import (
	"HouseKeeperBot/cache"
	. "HouseKeeperBot/common"
	"HouseKeeperBot/database"
	"HouseKeeperBot/modules/global/methods"
	"HouseKeeperBot/modules/global/models"
	"fmt"
	"github.com/mukeran/telegram-bot-api"
	"strconv"
)

const (
	tplManageAuth                      = `Whitelist mode: %v`
	tplSuccessfullyToggleWhitelistMode = `Successfully set whitelist mode to %v`
)

type paramIntID struct {
	ID int `json:"id"`
}

func generateManageMenuKeyboard() (keyboard tgbotapi.InlineKeyboardMarkup) {
	keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Auth", cache.RecordCallback(CallbackManageAuth, NoParam)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Admin", cache.RecordCallback(CallbackManageAdmin, NoParam)),
		),
	)
	return
}

func generateManageAuthKeyboard() (keyboard tgbotapi.InlineKeyboardMarkup) {
	keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Toggle whitelist mode",
				cache.RecordCallback(CallbackManageToggleWhitelistMode, NoParam)),
			tgbotapi.NewInlineKeyboardButtonData("Whitelist",
				cache.RecordCallback(CallbackManageWhitelist, NoParam)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Blacklist",
				cache.RecordCallback(CallbackManageBlacklist, NoParam)),
			tgbotapi.NewInlineKeyboardButtonData("< Back to manage",
				cache.RecordCallback(CallbackBackToManage, NoParam)),
		),
	)
	return
}

func generateManageAdminKeyboard() (keyboard tgbotapi.InlineKeyboardMarkup) {
	keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("List admins",
				cache.RecordCallback(CallbackManageListAdmins, NoParam)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Add admin",
				cache.RecordCallback(CallbackManageAddAdmin, NoParam)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("< Back to manage",
				cache.RecordCallback(CallbackBackToManage, NoParam)),
		),
	)
	return
}

func generateManageAdminListKeyboard() (keyboard tgbotapi.InlineKeyboardMarkup) {
	tx := database.Db
	var admins []models.User
	DatabasePanicError(tx.Where("is_admin = 1").Find(&admins))
	var buttons [][]tgbotapi.InlineKeyboardButton
	buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("< Back",
			cache.RecordCallback(CallbackManageAdmin, NoParam)),
	))
	for _, admin := range admins {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL(strconv.Itoa(admin.TelegramUserID),
				fmt.Sprintf("tg://user?id=%v", admin.TelegramUserID)),
			tgbotapi.NewInlineKeyboardButtonData("<- Delete <-",
				cache.RecordCallback(CallbackManageDeleteAdmin, EncodeParam(paramIntID{ID: admin.TelegramUserID}))),
		))
	}
	buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("< Back",
			cache.RecordCallback(CallbackManageAdmin, NoParam)),
	))
	keyboard = tgbotapi.NewInlineKeyboardMarkup(buttons...)
	return
}

func generateManageMenu(chatID int64) (resp tgbotapi.MessageConfig) {
	resp = tgbotapi.NewMessage(chatID, "Welcome to HouseKeeperBot manage menu")
	resp.ReplyMarkup = generateManageMenuKeyboard()
	return
}

func generateEditManageMenu(chatID int64, messageID int) (resp tgbotapi.EditMessageTextConfig) {
	resp = tgbotapi.NewEditMessageText(chatID, messageID, "Welcome to HouseKeeperBot manage menu")
	keyboard := generateManageMenuKeyboard()
	resp.ReplyMarkup = &keyboard
	return
}

func generateEditManageAuthMenu(chatID int64, messageID int) (resp tgbotapi.EditMessageTextConfig) {
	whitelistMode := methods.GetConfig(models.ConfigWhitelistMode)
	resp = tgbotapi.NewEditMessageText(chatID, messageID, fmt.Sprintf(tplManageAuth, whitelistMode))
	keyboard := generateManageAuthKeyboard()
	resp.ReplyMarkup = &keyboard
	return
}

func generateEditManageAdminMenu(chatID int64, messageID int) (resp tgbotapi.EditMessageTextConfig) {
	resp = tgbotapi.NewEditMessageText(chatID, messageID, "Please select one option")
	keyboard := generateManageAdminKeyboard()
	resp.ReplyMarkup = &keyboard
	return
}

func generateEditManageAdminListMenu(chatID int64, messageID int) (resp tgbotapi.EditMessageTextConfig) {
	resp = tgbotapi.NewEditMessageText(chatID, messageID, "Admin list:")
	keyboard := generateManageAdminListKeyboard()
	resp.ReplyMarkup = &keyboard
	return
}

func Manage() CommandHandlerFunc {
	return func(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, from *tgbotapi.User) {
		if !methods.IsAdmin(from.ID) {
			QuickSendTextMessage(bot, msg.Chat.ID, "Permission denied")
			return
		}
		MustSend(bot, generateManageMenu(msg.Chat.ID))
	}
}

func OnBackToManageButtonClick() CallbackQueryHandlerFunc {
	return func(bot *tgbotapi.BotAPI, lastMsg *tgbotapi.Message, from *tgbotapi.User,
		callbackQueryID string, param string) {
		if !methods.IsAdmin(from.ID) {
			QuickSendTextMessage(bot, lastMsg.Chat.ID, "Permission denied")
			return
		}
		MustSend(bot, generateEditManageMenu(lastMsg.Chat.ID, lastMsg.MessageID))
	}
}

func OnManageAuthButtonClick() CallbackQueryHandlerFunc {
	return func(bot *tgbotapi.BotAPI, lastMsg *tgbotapi.Message, from *tgbotapi.User,
		callbackQueryID string, param string) {
		if !methods.IsAdmin(from.ID) {
			QuickSendTextMessage(bot, lastMsg.Chat.ID, "Permission denied")
			return
		}
		MustSend(bot, generateEditManageAuthMenu(lastMsg.Chat.ID, lastMsg.MessageID))
	}
}

func OnManageAdminButtonClick() CallbackQueryHandlerFunc {
	return func(bot *tgbotapi.BotAPI, lastMsg *tgbotapi.Message, from *tgbotapi.User,
		callbackQueryID string, param string) {
		if !methods.IsAdmin(from.ID) {
			QuickSendTextMessage(bot, lastMsg.Chat.ID, "Permission denied")
			return
		}
		MustSend(bot, generateEditManageAdminMenu(lastMsg.Chat.ID, lastMsg.MessageID))
	}
}

func OnManageToggleWhitelistModeButtonClick() CallbackQueryHandlerFunc {
	return func(bot *tgbotapi.BotAPI, lastMsg *tgbotapi.Message, from *tgbotapi.User,
		callbackQueryID string, param string) {
		if !methods.IsAdmin(from.ID) {
			QuickSendTextMessage(bot, lastMsg.Chat.ID, "Permission denied")
			return
		}
		whitelistMode := methods.GetConfig(models.ConfigWhitelistMode)
		if whitelistMode == "on" {
			whitelistMode = "off"
		} else {
			whitelistMode = "on"
		}
		methods.SetConfig(models.ConfigWhitelistMode, whitelistMode)
		MustAnswerCallbackQuery(bot, tgbotapi.NewCallback(callbackQueryID,
			fmt.Sprintf(tplSuccessfullyToggleWhitelistMode, whitelistMode)))
		MustSend(bot, generateEditManageAuthMenu(lastMsg.Chat.ID, lastMsg.MessageID))
	}
}

func OnManageWhitelistButtonClick() CallbackQueryHandlerFunc {
	return func(bot *tgbotapi.BotAPI, lastMsg *tgbotapi.Message, from *tgbotapi.User,
		callbackQueryID string, param string) {
		if !methods.IsAdmin(from.ID) {
			QuickSendTextMessage(bot, lastMsg.Chat.ID, "Permission denied")
			return
		}
		QuickSendTextMessage(bot, lastMsg.Chat.ID,
			`Please share a user here to manage whitelist`)
		cache.RecordProcedure(lastMsg.Chat.ID, from.ID, ProcedureManageWhitelist, NoParam)
	}
}

func generateToggleIsWhitelistedKeyboard(telegramUserID int) (keyboard tgbotapi.InlineKeyboardMarkup) {
	keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Toggle status",
				cache.RecordCallback(CallbackManageToggleIsWhitelisted,
					EncodeParam(paramIntID{ID: telegramUserID}))),
		),
	)
	return
}

func generateToggleIsWhitelisted(chatID int64, telegramUserID int) (resp tgbotapi.MessageConfig) {
	resp = tgbotapi.NewMessage(chatID, fmt.Sprintf("User %v is %v the whitelist", telegramUserID,
		func() string {
			if methods.IsWhitelisted(telegramUserID) {
				return "*in*"
			} else {
				return "*not in*"
			}
		}()))
	resp.ParseMode = tgbotapi.ModeMarkdownV2
	resp.ReplyMarkup = generateToggleIsWhitelistedKeyboard(telegramUserID)
	return
}

func generateEditToggleIsWhitelisted(chatID int64, messageID int,
	telegramUserID int) (resp tgbotapi.EditMessageTextConfig) {
	resp = tgbotapi.NewEditMessageText(chatID, messageID, fmt.Sprintf("User %v is %v the whitelist", telegramUserID,
		func() string {
			if methods.IsWhitelisted(telegramUserID) {
				return "*in*"
			} else {
				return "*not in*"
			}
		}()))
	resp.ParseMode = tgbotapi.ModeMarkdownV2
	keyboard := generateToggleIsWhitelistedKeyboard(telegramUserID)
	resp.ReplyMarkup = &keyboard
	return
}

func ProcedureWhitelist() ProcedureHandlerFunc {
	return func(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, from *tgbotapi.User, param string) {
		if msg.Contact == nil {
			QuickSendTextMessage(bot, msg.Chat.ID, "Please share a user to this bot")
			return
		}
		cache.ClearProcedure(msg.Chat.ID, from.ID)
		if !methods.IsAdmin(from.ID) {
			QuickSendTextMessage(bot, msg.Chat.ID, "Permission denied")
			return
		}
		MustSend(bot, generateToggleIsWhitelisted(msg.Chat.ID, msg.Contact.UserID))
	}
}

func OnManageToggleIsWhitelistedButtonClick() CallbackQueryHandlerFunc {
	return func(bot *tgbotapi.BotAPI, lastMsg *tgbotapi.Message, from *tgbotapi.User,
		callbackQueryID string, param string) {
		if !methods.IsAdmin(from.ID) {
			QuickSendTextMessage(bot, lastMsg.Chat.ID, "Permission denied")
			return
		}
		var params paramIntID
		DecodeParam(param, &params)
		tx := database.Db.Begin()
		defer tx.RollbackUnlessCommitted()
		var user models.User
		DatabasePanicError(tx.FirstOrCreate(&user, models.User{TelegramUserID: params.ID}))
		user.IsWhitelisted = !user.IsWhitelisted
		DatabasePanicError(tx.Save(&user))
		DatabasePanicError(tx.Commit())
		QuickAnswerCallbackQuery(bot, callbackQueryID, fmt.Sprintf("Successfully %v whitelist", func() string {
			if user.IsWhitelisted {
				return fmt.Sprintf("added user %v into", params.ID)
			} else {
				return fmt.Sprintf("removed user %v from", params.ID)
			}
		}()))
		MustSend(bot, generateEditToggleIsWhitelisted(lastMsg.Chat.ID, lastMsg.MessageID, params.ID))
	}
}

func OnManageBlacklistButtonClick() CallbackQueryHandlerFunc {
	return func(bot *tgbotapi.BotAPI, lastMsg *tgbotapi.Message, from *tgbotapi.User,
		callbackQueryID string, param string) {
		if !methods.IsAdmin(from.ID) {
			QuickSendTextMessage(bot, lastMsg.Chat.ID, "Permission denied")
			return
		}
		QuickSendTextMessage(bot, lastMsg.Chat.ID,
			`Please share a user here to manage blacklist`)
		cache.RecordProcedure(lastMsg.Chat.ID, from.ID, ProcedureManageBlacklist, NoParam)
	}
}

func generateToggleIsBlacklistedKeyboard(telegramUserID int) (keyboard tgbotapi.InlineKeyboardMarkup) {
	keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Toggle status",
				cache.RecordCallback(CallbackManageToggleIsBlacklisted,
					EncodeParam(paramIntID{ID: telegramUserID}))),
		),
	)
	return
}

func generateToggleIsBlacklisted(chatID int64, telegramUserID int) (resp tgbotapi.MessageConfig) {
	resp = tgbotapi.NewMessage(chatID, fmt.Sprintf("User %v is %v the blacklist", telegramUserID,
		func() string {
			if methods.IsBlacklisted(telegramUserID) {
				return "*in*"
			} else {
				return "*not in*"
			}
		}()))
	resp.ParseMode = tgbotapi.ModeMarkdownV2
	resp.ReplyMarkup = generateToggleIsBlacklistedKeyboard(telegramUserID)
	return
}

func generateEditToggleIsBlacklisted(chatID int64, messageID int,
	telegramUserID int) (resp tgbotapi.EditMessageTextConfig) {
	resp = tgbotapi.NewEditMessageText(chatID, messageID,
		fmt.Sprintf("User %v is %v the blacklist", telegramUserID, func() string {
			if methods.IsBlacklisted(telegramUserID) {
				return "*in*"
			} else {
				return "*not in*"
			}
		}()))
	resp.ParseMode = tgbotapi.ModeMarkdownV2
	keyboard := generateToggleIsBlacklistedKeyboard(telegramUserID)
	resp.ReplyMarkup = &keyboard
	return
}

func ProcedureBlacklist() ProcedureHandlerFunc {
	return func(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, from *tgbotapi.User, param string) {
		if msg.Contact == nil {
			QuickSendTextMessage(bot, msg.Chat.ID, "Please share a user to this bot")
			return
		}
		cache.ClearProcedure(msg.Chat.ID, from.ID)
		if !methods.IsAdmin(from.ID) {
			QuickSendTextMessage(bot, msg.Chat.ID, "Permission denied")
			return
		}
		MustSend(bot, generateToggleIsBlacklisted(msg.Chat.ID, msg.Contact.UserID))
	}
}

func OnManageToggleIsBlacklistedButtonClick() CallbackQueryHandlerFunc {
	return func(bot *tgbotapi.BotAPI, lastMsg *tgbotapi.Message, from *tgbotapi.User,
		callbackQueryID string, param string) {
		if !methods.IsAdmin(from.ID) {
			QuickSendTextMessage(bot, lastMsg.Chat.ID, "Permission denied")
			return
		}
		var params paramIntID
		DecodeParam(param, &params)
		tx := database.Db.Begin()
		defer tx.RollbackUnlessCommitted()
		var user models.User
		DatabasePanicError(tx.FirstOrCreate(&user, models.User{TelegramUserID: params.ID}))
		user.IsBlacklisted = !user.IsBlacklisted
		DatabasePanicError(tx.Save(&user))
		DatabasePanicError(tx.Commit())
		QuickAnswerCallbackQuery(bot, callbackQueryID, fmt.Sprintf("Successfully %v blacklist", func() string {
			if user.IsBlacklisted {
				return fmt.Sprintf("added user %v into", params.ID)
			} else {
				return fmt.Sprintf("removed user %v from", params.ID)
			}
		}()))
		MustSend(bot, generateEditToggleIsBlacklisted(lastMsg.Chat.ID, lastMsg.MessageID, params.ID))
	}
}

func OnManageListAdminsButtonClick() CallbackQueryHandlerFunc {
	return func(bot *tgbotapi.BotAPI, lastMsg *tgbotapi.Message, from *tgbotapi.User,
		callbackQueryID string, param string) {
		if !methods.IsAdmin(from.ID) {
			QuickSendTextMessage(bot, lastMsg.Chat.ID, "Permission denied")
			return
		}
		MustSend(bot, generateEditManageAdminListMenu(lastMsg.Chat.ID, lastMsg.MessageID))
	}
}

func OnManageDeleteAdminButtonClick() CallbackQueryHandlerFunc {
	return func(bot *tgbotapi.BotAPI, lastMsg *tgbotapi.Message, from *tgbotapi.User,
		callbackQueryID string, param string) {
		if !methods.IsAdmin(from.ID) {
			QuickSendTextMessage(bot, lastMsg.Chat.ID, "Permission denied")
			return
		}
		var params paramIntID
		DecodeParam(param, &params)
		if params.ID == from.ID {
			QuickAnswerCallbackQuery(bot, callbackQueryID, "You can't delete yourself from admin list")
			return
		}
		methods.DeleteAdmin(params.ID)
		QuickAnswerCallbackQuery(bot, callbackQueryID,
			fmt.Sprintf("Successfully deleted user %v from admin list", params.ID))
		MustSend(bot, generateEditManageAdminListMenu(lastMsg.Chat.ID, lastMsg.MessageID))
	}
}

func OnManageAddAdminButtonClick() CallbackQueryHandlerFunc {
	return func(bot *tgbotapi.BotAPI, lastMsg *tgbotapi.Message, from *tgbotapi.User,
		callbackQueryID string, param string) {
		if !methods.IsAdmin(from.ID) {
			QuickSendTextMessage(bot, lastMsg.Chat.ID, "Permission denied")
			return
		}
		QuickSendTextMessage(bot, lastMsg.Chat.ID,
			`Please share a user here to add admin`)
		cache.RecordProcedure(lastMsg.Chat.ID, from.ID, ProcedureManageAddAdmin, NoParam)
	}
}

func generateRevertAddAdmin(chatID int64, telegramUserID int) (resp tgbotapi.MessageConfig) {
	resp = tgbotapi.NewMessage(chatID, fmt.Sprintf("Successfully set user %v as admin", telegramUserID))
	resp.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Revert",
				cache.RecordCallback(CallbackManageRevertAddAdmin,
					EncodeParam(paramIntID{ID: telegramUserID}))),
		),
	)
	return
}

func ProcedureAddAdmin() ProcedureHandlerFunc {
	return func(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, from *tgbotapi.User, param string) {
		if msg.Contact == nil {
			QuickSendTextMessage(bot, msg.Chat.ID, "Please share a user to this bot")
			return
		}
		cache.ClearProcedure(msg.Chat.ID, from.ID)
		if !methods.IsAdmin(from.ID) {
			QuickSendTextMessage(bot, msg.Chat.ID, "Permission denied")
			return
		}
		methods.AddAdmin(msg.Contact.UserID)
		MustSend(bot, generateRevertAddAdmin(msg.Chat.ID, msg.Contact.UserID))
	}
}

func OnManageRevertAddAdminButtonClick() CallbackQueryHandlerFunc {
	return func(bot *tgbotapi.BotAPI, lastMsg *tgbotapi.Message, from *tgbotapi.User,
		callbackQueryID string, param string) {
		if !methods.IsAdmin(from.ID) {
			QuickSendTextMessage(bot, lastMsg.Chat.ID, "Permission denied")
			return
		}
		var params paramIntID
		DecodeParam(param, &params)
		methods.DeleteAdmin(params.ID)
		QuickAnswerCallbackQuery(bot, callbackQueryID,
			fmt.Sprintf("Successfully revert add admin for user %v", params.ID))
		resp := tgbotapi.NewEditMessageText(lastMsg.Chat.ID, lastMsg.MessageID,
			"~"+EscapeMarkdownV2(lastMsg.Text)+"~")
		resp.ParseMode = tgbotapi.ModeMarkdownV2
		MustSend(bot, resp)
	}
}
