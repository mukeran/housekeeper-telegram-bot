package main

import (
	"github.com/mukeran/housekeeper-telegram-bot/cache"
	. "github.com/mukeran/housekeeper-telegram-bot/common"
	tgbotapi "github.com/mukeran/telegram-bot-api"
)

func handleMessage(bot *tgbotapi.BotAPI, msg tgbotapi.Message) {
	defer recovery(bot, &msg)
	if !auth(bot, msg.Chat.ID, msg.From.ID) {
		return
	}
	if !msg.IsCommand() {
		callback, param := cache.ResumeProcedure(msg.Chat.ID, msg.From.ID)
		callbackFunc := procedureHandlerMap[callback]
		if callbackFunc == nil {
			QuickSendTextMessage(bot, msg.Chat.ID, "Please send a command")
		} else {
			callbackFunc(bot, &msg, msg.From, param)
		}
	} else {
		f, existing := commandHandlerMap[msg.Command()]
		if !existing {
			QuickSendTextMessage(bot, msg.Chat.ID, "No such command")
		} else {
			f(bot, &msg, msg.From)
		}
	}
}

func handleCallbackQuery(bot *tgbotapi.BotAPI, query tgbotapi.CallbackQuery) {
	defer recovery(bot, query.Message)
	if !auth(bot, query.Message.Chat.ID, query.From.ID) {
		return
	}
	callback, param, err := cache.ResumeCallback(query.Data)
	if err != nil {
		QuickSendTextMessage(bot, query.Message.Chat.ID, "Bad request!")
		return
	}
	callbackFunc := callbackQueryHandlerMap[callback]
	if callbackFunc == nil {
		resp := tgbotapi.NewMessage(query.Message.Chat.ID, "Bad request!")
		MustSend(bot, resp)
	} else {
		callbackFunc(bot, query.Message, query.From, query.ID, param)
	}
}
