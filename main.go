package main

import (
	"fmt"
	"github.com/imdario/mergo"
	"github.com/mukeran/housekeeper-telegram-bot/cache"
	. "github.com/mukeran/housekeeper-telegram-bot/common"
	"github.com/mukeran/housekeeper-telegram-bot/database"
	"github.com/mukeran/housekeeper-telegram-bot/modules/cccat"
	"github.com/mukeran/housekeeper-telegram-bot/modules/global"
	"github.com/mukeran/housekeeper-telegram-bot/schedule"
	tgbotapi "github.com/mukeran/telegram-bot-api"
	"io"
	"log"
	"os"
	"time"
)

var (
	commandHandlerMap       CommandHandlerMap
	procedureHandlerMap     ProcedureHandlerMap
	callbackQueryHandlerMap CallbackQueryHandlerMap
	logFile                 *os.File
)

func main() {
	setupLogFile()
	defer func() {
		if logFile != nil {
			logFile.Close()
		}
	}()

	shouldInitialize := database.Connect()
	if shouldInitialize {
		initializeDatabase()
	}
	cache.Connect()

	log.Println("Connecting to Telegram server...")
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}
	log.Println("Successfully connected to Telegram server!")

	if os.Getenv("DEBUG") != "" {
		bot.Debug = true
		log.Println("Debug mode is on")
	}
	log.Printf("Authorized on account %s\n", bot.Self.UserName)

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
	updates := bot.GetUpdatesChan(u)

	log.Println("Start to fetch updates...")
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

func setupLogFile() {
	var err error
	if err = os.Mkdir("logs", 0755); err != nil && !os.IsExist(err) {
		panic(err)
	}
	if logFile, err = os.OpenFile(fmt.Sprintf("logs/log_%v.log", time.Now().Format(time.RFC3339)),
		os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644); err == nil {
		mw := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(mw)
	}
}
