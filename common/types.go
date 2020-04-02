package common

import tgbotapi "github.com/mukeran/telegram-bot-api"

/// Handler function definitions
type (
	CommandHandlerFunc       func(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, from *tgbotapi.User)
	ProcedureHandlerFunc     func(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, from *tgbotapi.User, param string)
	CallbackQueryHandlerFunc func(bot *tgbotapi.BotAPI, lastMsg *tgbotapi.Message, from *tgbotapi.User,
		callbackQueryID string, param string)
)

/// Handler map definitions
type (
	CommandHandlerMap       map[string]CommandHandlerFunc
	ProcedureHandlerMap     map[string]ProcedureHandlerFunc
	CallbackQueryHandlerMap map[string]CallbackQueryHandlerFunc
)

/// Key-value map
type H map[string]interface{}

const (
	NoParam = ""
)

type ParamID struct {
	ID uint `json:"id"`
}
