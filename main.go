package main

import (
	"HouseKeeperBot/cache"
	. "HouseKeeperBot/common"
	"HouseKeeperBot/database"
	"HouseKeeperBot/modules/cccat"
	"HouseKeeperBot/modules/global"
	"HouseKeeperBot/schedule"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/imdario/mergo"
	tgbotapi "github.com/mukeran/telegram-bot-api"
)

var (
	commandHandlerMap       CommandHandlerMap
	procedureHandlerMap     ProcedureHandlerMap
	callbackQueryHandlerMap CallbackQueryHandlerMap
)

func main() {
	shouldInitialize := database.Connect()
	if shouldInitialize {
		initializeDatabase()
	}
	cache.Connect()

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}
	//bot.Debug = true
	if file, err := os.OpenFile(fmt.Sprintf("logs/log_%v.log", time.Now().Format(time.RFC3339)),
		os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666); err == nil {
		defer file.Close()
		log.SetOutput(file)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)

	_ = mergo.Merge(&commandHandlerMap, global.Commands())
	_ = mergo.Merge(&commandHandlerMap, cccat.Commands())

	_ = mergo.Merge(&procedureHandlerMap, global.Procedures())
	_ = mergo.Merge(&procedureHandlerMap, cccat.Procedures())

	_ = mergo.Merge(&callbackQueryHandlerMap, global.CallbackQueries())
	_ = mergo.Merge(&callbackQueryHandlerMap, cccat.CallbackQueries())

	cccat.RegisterSchedules()
	schedule.Start(bot)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		msg := update.Message
		if msg != nil {
			go handleMessage(bot, *msg)
		}
		cq := update.CallbackQuery
		if cq != nil {
			go handleCallbackQuery(bot, *cq)
		}
	}
}
