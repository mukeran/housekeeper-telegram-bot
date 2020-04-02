package main

import (
	tgbotapi "github.com/mukeran/telegram-bot-api"
	"log"
)

func recovery(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	if err := recover(); err != nil {
		log.Printf("Panic recovered from %v", err)
		resp := tgbotapi.NewMessage(msg.Chat.ID, "Bot error! Please contact administrator.")
		if _, err := bot.Send(resp); err != nil {
			log.Printf("Error occurs when sending error notice. %v", err)
		}
	}
}
