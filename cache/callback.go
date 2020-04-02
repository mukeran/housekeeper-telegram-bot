package cache

import (
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"log"
)

const (
	fieldCallbackQueryHandler = "callbackQueryHandler"
)

var (
	ErrCallbackNotFound = errors.New("callback not found")
)

func getCallbackID(token string) string {
	return fmt.Sprintf("callback-%v", token)
}

func RecordCallback(callbackQueryHandler string, param string) (token string) {
	conn := Redis.Get()
	defer conn.Close()
	token = uuid.New().String()
	_, err := conn.Do("HMSET", getCallbackID(token),
		fieldCallbackQueryHandler, callbackQueryHandler,
		fieldParam, param)
	if err != nil {
		log.Panic(err)
	}
	return
}

func ResumeCallback(token string) (callbackQueryHandler string, param string, err error) {
	conn := Redis.Get()
	defer conn.Close()
	values, err := redis.Values(conn.Do("HMGET", getCallbackID(token),
		fieldCallbackQueryHandler, fieldParam))
	if err != nil {
		log.Panic(err)
	}
	if values[0] == nil {
		return "", "", ErrCallbackNotFound
	}
	if _, err = redis.Scan(values, &callbackQueryHandler, &param); err != nil {
		log.Panic(err)
	}
	return
}
