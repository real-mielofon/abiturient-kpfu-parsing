package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/txgruppi/werr"
	"gitlab.com/toby3d/telegram"
)

func _log(err error, message string, v ...interface{}) {
	log.Println(err)
	log.Println(fmt.Sprintf(message, v...))     // Return the original error message
	if wrapped, ok := err.(*werr.Wrapper); ok { // Try to convert to `*werr.Wrapper`
		lg, _ := wrapped.Log() // Generate the log message
		for _, line := range strings.Split(lg, "\n") {
			log.Println(line) // Print the log message
		}
	}
}

type App struct {
	config *config
	list   ListAbiturents
}

func (app App) init(config *Config) {
	app.config = config
}
func (app App) Run() {
	bot, err := telegram.New(env)
	err = werr.Wrap(err)
	if err != nil {
		_log(err, "telegram.New(%+v)", env)
		return
	}
	log.Printf("Authorized on account %s", bot.Username)

	updatesParameters := &telegram.GetUpdatesParameters{
		Offset:  0,
		Limit:   100,
		Timeout: 120,
	}

	updates := app.bot.NewLongPollingChannel(updatesParameters)

	//app.bot.Debug = true

	ticker := time.NewTicker(periodUpdate)
	for {
		select {
		case update := <-app.updates:
			if update.Message == nil {
				log.Printf("Out message nil")
				continue
			}

			messageText := update.Message.Text
			messageUserName := update.Message.From.Username
			messageChatID := update.Message.Chat.ID

			log.Printf("In  [%s] id:%d %s", messageUserName, messageChatID, messageText)

			switch {
			case (messageText == "/list") || (messageText == "/l"):
				tgBotCommandList(bot, messageChatID)
				log.Printf("Out [%s] id:%d text:%s Ok", messageUserName, messageChatID, messageText)
			case (messageText == "/status") || (messageText == "/s"):
				tgBotCommandStatWithGetStatus(bot, messageChatID, config)
				log.Printf("Out [%s] id:%d text:%s Ok", messageUserName, messageChatID, messageText)
			case strings.HasPrefix(messageText, "/subscribe"):
				abiturientName := strings.TrimPrefix(messageText, "/subscribe"+" ")
				tgBotCommandSubscribe(bot, config, messageChatID, abiturientName)
				log.Printf("Out [%s] id:%d text:%s Ok", messageUserName, messageChatID, messageText)
			case (messageText == "/unsubscribe"):
				tgBotCommandUnSubscribe(bot, config, messageChatID)
				log.Printf("Out [%s] id:%d text:%s Ok", messageUserName, messageChatID, messageText)
			case (messageText == "/ping") || (messageText == "/p"):
				tgBotCommandPing(bot, messageChatID)
				log.Printf("Out [%s] id:%d text:%s Ok", messageUserName, messageChatID, messageText)
			default:
				log.Printf("Out [%s] id:%d text:%s Ok", messageUserName, messageChatID, messageText)
			}
		case <-ticker.C:
			now := time.Now()
			hour := now.Hour()
			if (hour >= 8) && (hour <= 21) {
				tgBotCommandSendChangeStatus(bot, config)
				log.Printf("Out Change Ok")
			}
		}
	}
}

func (app App) GetStatusAbiturient(name string) (*StatusAbiturienta, error) {
	arr := app.list.Get()
	if err != nil {
		return nil, werr.Wrap(err)
	}

	status := StatusAbiturienta{Num: 0, NumWithOriginal: 0}
	for _, ab := range app.list.arr {
		if ab.Fio == name {
			status.NumWithOriginal++
			status.Num = ab.Num
			break
		}
		if ab.Original {
			status.NumWithOriginal++
		}
	}
	return &status, nil
}
