package app

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"strings"
	"time"

	"github.com/real-mielofon/abiturient-kpfu-parsing/config"
	"github.com/real-mielofon/abiturient-kpfu-parsing/status"

	"github.com/txgruppi/werr"
	"gitlab.com/toby3d/telegram"
)

const (
	urlList = "https://abiturient.kpfu.ru/entrant/abit_entrant_originals_list?p_open=&p_typeofstudy=1&p_faculty=47&p_speciality=1085&p_inst=0&p_category=1"
	//nameFindAbiturient = "Пономарев Степан Алексеевич"
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
	cfg          *config.Config
	env          string
	periodUpdate time.Duration

	bot  *telegram.Bot
	list ListAbiturents
}

func New(config *config.Config, env string, periodUpdate time.Duration) (app App) {
	app.cfg = config
	app.env = env
	app.periodUpdate = periodUpdate
	return app
}
func (app App) Run() {
	bot, err := telegram.New(app.env)
	err = werr.Wrap(err)
	if err != nil {
		_log(err, "telegram.New(%+v)", app.env)
		return
	}
	app.bot = bot
	log.Printf("Authorized on account %s", app.bot.Username)

	updatesParameters := &telegram.GetUpdatesParameters{
		Offset:  0,
		Limit:   100,
		Timeout: 120,
	}

	updates := app.bot.NewLongPollingChannel(updatesParameters)

	//app.bot.Debug = true

	ticker := time.NewTicker(app.periodUpdate)
	for {
		select {
		case update := <-updates:
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
				app.tgBotCommandList(messageChatID)
				log.Printf("Out [%s] id:%d text:%s Ok", messageUserName, messageChatID, messageText)
			case (messageText == "/status") || (messageText == "/s"):
				app.tgBotCommandStatWithGetStatus(messageChatID)
				log.Printf("Out [%s] id:%d text:%s Ok", messageUserName, messageChatID, messageText)
			case strings.HasPrefix(messageText, "/subscribe"):
				abiturientName := strings.TrimPrefix(messageText, "/subscribe"+" ")
				app.tgBotCommandSubscribe(messageChatID, abiturientName)
				log.Printf("Out [%s] id:%d text:%s Ok", messageUserName, messageChatID, messageText)
			case (messageText == "/unsubscribe"):
				app.tgBotCommandUnSubscribe(messageChatID)
				log.Printf("Out [%s] id:%d text:%s Ok", messageUserName, messageChatID, messageText)
			case (messageText == "/ping") || (messageText == "/p"):
				app.tgBotCommandPing(messageChatID)
				log.Printf("Out [%s] id:%d text:%s Ok", messageUserName, messageChatID, messageText)
			default:
				log.Printf("Out [%s] id:%d text:%s Ok", messageUserName, messageChatID, messageText)
			}
		case <-ticker.C:
			now := time.Now()
			hour := now.Hour()
			if (hour >= 8) && (hour <= 21) {
				app.tgBotCommandSendChangeStatus()
				log.Printf("Out Change Ok")
			}
		}
	}
}

func (app App) GetStatusAbiturient(name string) (*status.StatusAbiturienta, error) {
	list, err := app.list.Get()
	err = werr.Wrap(err)
	if err != nil {
		return nil, err
	}

	status := status.StatusAbiturienta{Num: 0, NumWithOriginal: 0}
	for _, ab := range list {
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

func (app App) tgBotCommandList(messageChatID int64) {

	name := app.cfg.Chats[messageChatID].Name

	arr, err := app.list.Get()
	err = werr.Wrap(err)
	if err != nil {
		_log(err, "getListAbiturient()")
		return
	}

	t := template.New("abiturients list")

	t, err = t.Parse(`{{range .}}<pre>{{if eq .Fio "` + name + `"}}>>>{{else}}{{if .Original}} * {{else}}   {{end}}{{end}}{{printf "%3d" .Num}} {{printf "%40s" .Fio}} {{index .Points 4|printf "%3d"}}{{if eq .Fio "` + name + `"}}<<<{{else}}{{if .Original}} * {{else}}   {{end}}{{end}}</pre>{{printf "\n"}}{{end}}`)
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

		app.SendHTML(messageChatID, text)
	}
}

func (app App) tgBotCommandStatWithGetStatus(messageChatID int64) {
	statusByName, found := app.cfg.Chats[messageChatID]
	if !found {
		app.SendText(messageChatID, "No subscribe")
		return
	}
	name := statusByName.Name

	status, err := app.GetStatusAbiturient(name)
	err = werr.Wrap(err)
	if err != nil {
		_log(err, "getStatusAbiturient()")
		return
	}
	app.tgBotCommandStat(messageChatID, name, status)
}

func (app App) tgBotCommandStat(messageChatID int64, name string, status *status.StatusAbiturienta) {

	t := template.New("abiturients status")

	t, err := t.Parse("Абитуриент *" + name + "*\nПерсональный рейтинг: *{{.Num}}*\nПерсональный рейтинг по оригиналам: *{{.NumWithOriginal}}*")
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

	app.SendMarkDown(messageChatID, text)

}

func (app App) SendMarkDown(chatID int64, text string) {
	app.Send(chatID, text, "markdown")
}
func (app App) SendText(chatID int64, text string) {
	app.Send(chatID, text, "text")
}
func (app App) SendHTML(chatID int64, text string) {
	app.Send(chatID, text, "html")
}
func (app App) Send(chatID int64, text string, parseMode string) {
	msg := telegram.NewMessage(chatID, text)
	msg.ParseMode = parseMode

	_, err := app.bot.SendMessage(msg)
	err = werr.Wrap(err)
	if err != nil {
		_log(err, "bot.SendMessage(%+v)", msg)
		return
	}
}

func (app App) tgBotCommandSubscribe(messageChatID int64, abiturientName string) {

	app.cfg.Add(messageChatID, abiturientName)
	app.cfg.WriteConfig()

	text := "Subscribed for " + abiturientName
	app.SendHTML(messageChatID, text)
}

func (app App) tgBotCommandUnSubscribe(messageChatID int64) {

	delete(app.cfg.Chats, messageChatID)
	app.cfg.WriteConfig()

	text := "UnSubscribed"
	app.SendHTML(messageChatID, text)
}

func (app App) tgBotCommandPing(messageChatID int64) {

	text := "pong"
	app.SendHTML(messageChatID, text)
}

func (app App) tgBotCommandSendChangeStatus() {
	configChanged := false
	for key, statusWithName := range app.cfg.Chats {
		name := statusWithName.Name
		statusSaved := statusWithName.Status

		sts, err := app.GetStatusAbiturient(name)
		err = werr.Wrap(err)
		if err != nil {
			_log(err, "getStatusAbiturient()")
			break
		}
		if sts.IsEqual(statusSaved) {
			// no change exit
			break
		}

		configChanged = true
		app.cfg.Chats[key] = status.StatusByName{Name: name, Status: *sts}
		app.tgBotCommandStat(key, name, sts)
	}
	if configChanged {
		app.cfg.WriteConfig()
	}

}
