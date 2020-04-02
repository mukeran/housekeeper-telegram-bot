package schedule

import (
	tgbotapi "github.com/mukeran/telegram-bot-api"
	"log"
	"time"
)

const (
	checkingDuration        = time.Minute * 10
	retryDurationWhenFailed = time.Minute * 10
)

func Start(bot *tgbotapi.BotAPI) {
	ticker := time.NewTicker(checkingDuration)
	go func() {
		now := time.Now()
		for range ticker.C {
			log.Printf("Start to execute schedule tasks")
			for i, task := range taskList {
				if !now.Before(task.NextTime) {
					b := task.Func(bot, task.Param)
					if !b {
						taskList[i].NextTime = task.NextTime.Add(retryDurationWhenFailed)
						log.Printf("Failed to execute scheduled task %v. Will retry at %v", i,
							taskList[i].NextTime.Format(time.RFC3339))
					} else {
						for !now.Before(taskList[i].SupposedNextTime) {
							taskList[i].SupposedNextTime = task.Rule.GetNextTime(task.SupposedNextTime)
						}
						taskList[i].NextTime = taskList[i].SupposedNextTime
						log.Printf("Successfully executed scheduled task %v. Next executing time is %v", i,
							taskList[i].NextTime.Format(time.RFC3339))
					}
				}
			}
		}
	}()
}
