package schedule

import (
	"log"
	"time"
)

var (
	taskList []Task
)

func RegisterDurationTask(f TaskFunc, param string, firstTime time.Time, duration time.Duration) {
	if duration == 0 {
		log.Panic("duration can't be zero")
	}
	task := Task{
		Func:             f,
		Param:            param,
		SupposedNextTime: firstTime,
		NextTime:         firstTime,
		Rule:             Rule{Duration: duration},
	}
	taskList = append(taskList, task)
}
