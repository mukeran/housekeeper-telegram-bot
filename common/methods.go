package common

import (
	"encoding/base64"
	"encoding/json"
	"github.com/imdario/mergo"
	"github.com/jinzhu/gorm"
	tgbotapi "github.com/mukeran/telegram-bot-api"
	"log"
	"strings"
)

// interface{} to string
func Interface2String(arg interface{}) (string, bool) {
	if arg == nil {
		return "", false
	}
	switch arg.(type) {
	case string:
		return arg.(string), true
	default:
		return "", false
	}
}

// interface{} to []byte
func Interface2Bytes(arg interface{}) ([]byte, bool) {
	if arg == nil {
		return nil, false
	}
	switch arg.(type) {
	case []byte:
		return arg.([]byte), true
	default:
		return nil, false
	}
}

// interface{} to H (map[string]interface{})
func Interface2H(arg interface{}) (H, bool) {
	if arg == nil {
		return nil, false
	}
	switch arg.(type) {
	case map[string]interface{}:
		return arg.(map[string]interface{}), true
	default:
		return nil, false
	}
}

// interface{} to bool
func Interface2Bool(arg interface{}) (bool, bool) {
	if arg == nil {
		return false, false
	}
	switch arg.(type) {
	case bool:
		return arg.(bool), true
	default:
		return false, false
	}
}

// interface{} to uint
func Interface2Uint(arg interface{}) (uint, bool) {
	if arg == nil {
		return 0, false
	}
	switch arg.(type) {
	case uint:
		return arg.(uint), true
	case float64:
		return uint(arg.(float64)), true
	default:
		return 0, false
	}
}

// Reduce the reduction of panic
func MustSend(bot *tgbotapi.BotAPI, chattable tgbotapi.Chattable) {
	_, err := bot.Send(chattable)
	if err != nil {
		log.Panic(err)
	}
}

func QuickSendTextMessage(bot *tgbotapi.BotAPI, chatID int64, text string) {
	MustSend(bot, tgbotapi.NewMessage(chatID, text))
}

func MustAnswerCallbackQuery(bot *tgbotapi.BotAPI, config tgbotapi.CallbackConfig) {
	_, err := bot.AnswerCallbackQuery(config)
	if err != nil {
		log.Panic(err)
	}
}

func QuickAnswerCallbackQuery(bot *tgbotapi.BotAPI, callbackQueryID string, text string) {
	MustAnswerCallbackQuery(bot, tgbotapi.NewCallback(callbackQueryID, text))
}

func QuickAnswerCallbackQueryWithAlert(bot *tgbotapi.BotAPI, callbackQueryID string, text string) {
	MustAnswerCallbackQuery(bot, tgbotapi.NewCallbackWithAlert(callbackQueryID, text))
}

func EncodeParam(param interface{}) string {
	b, err := json.Marshal(param)
	if err != nil {
		log.Panic(err)
	}
	return base64.StdEncoding.EncodeToString(b)
}

func DecodeParam(encoded string, dst interface{}) {
	var (
		b   []byte
		err error
	)
	b, err = base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		log.Panic(err)
	}
	err = json.Unmarshal(b, dst)
	if err != nil {
		log.Panic(err)
	}
}

func DatabasePanicError(results ...*gorm.DB) {
	for _, result := range results {
		if result.Error != nil {
			log.Panic(result.Error)
		}
	}
}

func EscapeMarkdownV2(str string) string {
	str = strings.ReplaceAll(str, "_", `\_`)
	str = strings.ReplaceAll(str, "*", `\*`)
	str = strings.ReplaceAll(str, "[", `\[`)
	str = strings.ReplaceAll(str, "]", `\]`)
	str = strings.ReplaceAll(str, "(", `\(`)
	str = strings.ReplaceAll(str, ")", `\)`)
	str = strings.ReplaceAll(str, "~", `\~`)
	str = strings.ReplaceAll(str, "`", "\\`")
	str = strings.ReplaceAll(str, "~", `\~`)
	str = strings.ReplaceAll(str, ">", `\>`)
	str = strings.ReplaceAll(str, "#", `\#`)
	str = strings.ReplaceAll(str, "+", `\+`)
	str = strings.ReplaceAll(str, "-", `\-`)
	str = strings.ReplaceAll(str, "=", `\=`)
	str = strings.ReplaceAll(str, "|", `\|`)
	str = strings.ReplaceAll(str, "{", `\}`)
	str = strings.ReplaceAll(str, "}", `\}`)
	str = strings.ReplaceAll(str, ".", `\.`)
	str = strings.ReplaceAll(str, "!", `\!`)
	return str
}

func MergeHandlerMaps(dest interface{}, sources ...interface{}) {
	for _, source := range sources {
		if source == nil {
			continue
		}
		err := mergo.Merge(dest, source)
		if err != nil {
			log.Panic(err)
		}
	}
}
