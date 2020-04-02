package schedule

import (
	tgbotapi "github.com/mukeran/telegram-bot-api"
	"time"
)

type TaskFunc func(bot *tgbotapi.BotAPI, param string) bool

type Task struct {
	Func             TaskFunc
	Param            string
	SupposedNextTime time.Time
	NextTime         time.Time
	Rule             Rule
}
