package cache

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
)

const (
	fieldMessageHandler = "messageHandler"
	fieldParam          = "param"
)

func getProcedureID(chatID int64, userID int) string {
	return fmt.Sprintf("%v%v", chatID, userID)
}

func RecordProcedure(chatID int64, userID int, procedureHandler string, param string) {
	conn := Redis.Get()
	defer conn.Close()
	_, err := conn.Do("HMSET", getProcedureID(chatID, userID),
		fieldMessageHandler, procedureHandler,
		fieldParam, param)
	if err != nil {
		log.Panic(err)
	}
}

func ResumeProcedure(chatID int64, userID int) (procedureHandler string, param string) {
	conn := Redis.Get()
	defer conn.Close()
	values, err := redis.Values(conn.Do("HMGET", getProcedureID(chatID, userID),
		fieldMessageHandler, fieldParam))
	if err != nil {
		log.Panic(err)
	}
	if _, err := redis.Scan(values, &procedureHandler, &param); err != nil {
		log.Panic(err)
	}
	return
}

func ClearProcedure(chatID int64, userID int) {
	conn := Redis.Get()
	defer conn.Close()
	_, err := conn.Do("DEL", getProcedureID(chatID, userID))
	if err != nil {
		log.Panic(err)
	}
}
