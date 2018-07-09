package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"log"

	"github.com/PuerkitoBio/goquery"
	"github.com/txgruppi/werr"
	"golang.org/x/net/html/charset"
)

const (
	urlList            = "https://abiturient.kpfu.ru/entrant/abit_entrant_originals_list?p_open=&p_typeofstudy=1&p_faculty=47&p_speciality=1085&p_inst=0&p_category=1"
	nameFindAbiturient = "Пономарев Степан Алексеевич"
	periodUpdate       = 30 * time.Minute
	fileConfig         = "./data/subscribe.txt"
)


// StatusAbiturienta is position abiturient in list
type StatusAbiturienta struct {
	Num             int
	NumWithOriginal int
}

// StatusByName is struct status with name of abiturient
type StatusByName struct {
	Name   string
	Status StatusAbiturienta
}


func tgBotCommandList(bot *telegram.Bot, messageChatID int64) {
	arr, err := getListAbiturient()
	err = werr.Wrap(err)
	if err != nil {
		_log(err, "getListAbiturient()")
		return
	}

	t := template.New("abiturients list")

	t, err = t.Parse(`{{range .}}<pre>{{if eq .Fio "` + nameFindAbiturient + `"}}>>>{{else}}{{if .Original}} * {{else}}   {{end}}{{end}}{{printf "%3d" .Num}} {{printf "%40s" .Fio}} {{index .Points 4|printf "%3d"}}{{if eq .Fio "` + nameFindAbiturient + `"}}<<<{{else}}{{if .Original}} * {{else}}   {{end}}{{end}}</pre>{{printf "\n"}}{{end}}`)
	err = werr.Wrap(err)
	if err != nil {
		_log(err, "t.Parse")
		return
	}

	for inter := 0; inter < len(arr)/20+1; inter++ {
		var b bytes.Buffer
		last := inter*20 + 20
		if last >= len(arr) {
			last = len(arr)
		}
		err = t.Execute(&b, arr[inter*20:last])
		err = werr.Wrap(err)
		if err != nil {
			_log(err, "t.Execute(%+v, %+v)", &b, arr[inter*20:last])
			return
		}

		text := b.String()
		if text == "" {
			text = "Empty :-("
		}
		msg := telegram.NewMessage(messageChatID, text)
		msg.ParseMode = "html"

		_, err := bot.SendMessage(msg)
		err = werr.Wrap(err)
		if err != nil {
			_log(err, "bot.SendMessage (%+v)", msg)
			return
		}
	}
}

func tgBotCommandStatWithGetStatus(bot *telegram.Bot, messageChatID int64, config *Config) {
	name := config.chats[messageChatID].Name

	status, err := getStatusAbiturient(name)
	err = werr.Wrap(err)
	if err != nil {
		_log(err, "getStatusAbiturient()")
		return
	}
	tgBotCommandStat(bot, messageChatID, status)
}

func tgBotCommandStat(bot *telegram.Bot, messageChatID int64, status *StatusAbiturienta) {

	t := template.New("abiturients status")

	t, err := t.Parse("Абитуриент *" + nameFindAbiturient + "*\nПерсональный рейтинг: *{{.Num}}*\nПерсональный рейтинг по оригиналам: *{{.NumWithOriginal}}*")
	err = werr.Wrap(err)
	if err != nil {
		_log(err, "t.Parse")
		return
	}

	var b bytes.Buffer
	err = t.Execute(&b, status)
	err = werr.Wrap(err)
	if err != nil {
		_log(err, "t.Execute(%+v, %+v)", &b, status)
		return
	}

	text := b.String()
	msg := telegram.NewMessage(messageChatID, text)
	msg.ParseMode = "markdown"

	_, err = bot.SendMessage(msg)
	err = werr.Wrap(err)
	if err != nil {
		_log(err, "bot.SendMessage(%+v)", msg)
		return
	}
}

func tgBotCommandSubscribe(bot *telegram.Bot, config *Config, messageChatID int64, abiturientName string) {

	config.Add(messageChatID, abiturientName)
	config.WriteConfig()

	text := "Subscribed for " + abiturientName
	msg := telegram.NewMessage(messageChatID, text)
	msg.ParseMode = "html"

	_, err := bot.SendMessage(msg)
	err = werr.Wrap(err)
	if err != nil {
		_log(err, "bot.SendMessage(%+v)", msg)
		return
	}
}

func tgBotCommandUnSubscribe(bot *telegram.Bot, config *Config, messageChatID int64) {

	delete(config.chats, messageChatID)
	config.WriteConfig()

	text := "UnSubscribed"
	msg := telegram.NewMessage(messageChatID, text)
	msg.ParseMode = "html"
	//			msg.ReplyToMessageID = update.Message.ID

	_, err := bot.SendMessage(msg)
	err = werr.Wrap(err)
	if err != nil {
		_log(err, "bot.SendMessage(%+v)", msg)
		return
	}
}

func tgBotCommandPing(bot *telegram.Bot, messageChatID int64) {

	text := "pong"
	msg := telegram.NewMessage(messageChatID, text)
	msg.ParseMode = "html"
	//			msg.ReplyToMessageID = update.Message.ID

	_, err := bot.SendMessage(msg)
	err = werr.Wrap(err)
	if err != nil {
		_log(err, "bot.SendMessage(%+v)", msg)
		return
	}
}

func tgBotCommandSendChangeStatus(bot *telegram.Bot, config *Config) {
	for key, status := range config {
		status, err := getStatusAbiturient()
		err = werr.Wrap(err)
		if err != nil {
		_log(err, "getStatusAbiturient()")
		return
		}
		if (status.Num == config.Status.Num) && (status.NumWithOriginal == config.Status.NumWithOriginal) {
			// no change exit
			return
		}
	
		config.status.Num = status.Num
		config.status.NumWithOriginal = status.NumWithOriginal
		config.WriteConfig()
	
		for key := range config.chats {
			tgBotCommandStat(bot, key, status)
		}
	}
}

func main() {

	env := os.Getenv("TGBOT_KEY")

	config := new(ConfigType)
	config.ReadConfig()

	app := App.init(config, env)
	app.Run()
}
