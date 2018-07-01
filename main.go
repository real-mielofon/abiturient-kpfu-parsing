package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"time"

	"log"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-errors/errors"
	"gitlab.com/toby3d/telegram"
	"golang.org/x/net/html/charset"
)

const (
	urlList            = "https://abiturient.kpfu.ru/entrant/abit_entrant_originals_list?p_open=&p_typeofstudy=1&p_faculty=47&p_speciality=1085&p_inst=0&p_category=1"
	nameFindAbiturient = "Пономарев Степан Алексеевич"
	periodUpdate       = 30 * time.Minute
	fileConfig         = "./data/subscribe.txt"
)

func _check(err error) {
	if err != nil {
		panic(err)
	}
}

//Abiturient is struct data of abiturient
type Abiturient struct {
	Num      int
	Fio      string
	Points   [5]int
	Original bool
}

// StatusAbiturienta is position abiturient in list
type StatusAbiturienta struct {
	Num             int
	NumWithOriginal int
}

func getStatusAbiturient() (*StatusAbiturienta, error) {
	arr, err := getListAbiturient()
	if err != nil {
		return nil, err
	}
	status := StatusAbiturienta{Num: 0, NumWithOriginal: 0}
	for _, ab := range arr {
		if ab.Fio == nameFindAbiturient {
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

func getListAbiturient() ([]Abiturient, error) {
	var cl = &http.Client{}
	cl.Timeout = 60 * time.Second

	resp, err := cl.Get(urlList)
	if err != nil {
		log.Println("HTTP error:", err, err.(*errors.Error).ErrorStack())
		return nil, fmt.Errorf("HTTP error: %v", err)
	}

	defer resp.Body.Close()
	// вот здесь и начинается самое интересное
	utf8, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		log.Println("Encoding error:", err, err.(*errors.Error).ErrorStack())
		return nil, fmt.Errorf("HTTP error: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(utf8)
	if err != nil {
		log.Println("goquery.NewDocumentFromReader:", err, err.(*errors.Error).ErrorStack())
		return nil, fmt.Errorf("goquery.NewDocumentFromReader: %v", err)
	}

	//fmt.Printf("doc: %+v \n", doc.Text())

	trs := doc.Find("tbody").Eq(1).Find("tr")
	len := trs.Length()
	arr := make([]Abiturient, len-1)

	//	fmt.Printf("len: %+v \n", len)
	for i := 1; i < len; i++ {

		tr := trs.Eq(i)
		tds := tr.Find("td")
		var ab Abiturient

		if style, styleExist := tr.Attr("style"); styleExist {
			ab.Original = (style == "font-weight:bold;")
		} else {
			ab.Original = false
		}

		s := tds.Eq(0).Text()
		num, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			log.Printf("strconv.ParseInt(%s), %v, %s\n", s, err, err.(*errors.Error).ErrorStack())
			return nil, fmt.Errorf("strconv.ParseInt(tds.Eq(0).Text(): %v", err)
		}
		ab.Num = int(num)

		ab.Fio = tds.Eq(1).Text()

		s = tds.Eq(2).Text()
		num, err = strconv.ParseInt(s, 10, 32)
		if err != nil {
			//			log.Printf("strconv.ParseInt(%s), %v\n", s, err)
			ab.Points[0] = 0
		} else {
			ab.Points[0] = int(num)
		}

		num, err = strconv.ParseInt(tds.Eq(3).Text(), 10, 32)
		if err != nil {
			//			log.Printf("strconv.ParseInt(%s), %v\n", s, err)
			ab.Points[1] = 0
		} else {
			ab.Points[1] = int(num)
		}

		num, err = strconv.ParseInt(tds.Eq(4).Text(), 10, 32)
		if err != nil {
			//			log.Printf("strconv.ParseInt(%s), %v\n", s, err)
			ab.Points[2] = 0
		} else {
			ab.Points[2] = int(num)
		}

		num, err = strconv.ParseInt(tds.Eq(5).Text(), 10, 32)
		if err != nil {
			//			log.Printf("strconv.ParseInt(%s), %v\n", s, err)
			ab.Points[3] = 0
		} else {
			ab.Points[3] = int(num)
		}

		num, err = strconv.ParseInt(tds.Eq(6).Text(), 10, 32)
		if err != nil {
			//			log.Printf("strconv.ParseInt(%s), %v\n", s, err)
			ab.Points[4] = ab.Points[0] + ab.Points[1] + ab.Points[2] + ab.Points[3]
		} else {
			ab.Points[4] = int(num)
		}

		arr[i-1] = ab

		//		log.Printf("%4d %40s %3d %s\n", ab.Num, ab.Fio, ab.Points[4], strconv.FormatBool(ab.Original))
	}
	return arr, nil
}

func tgBotCommandList(bot *telegram.Bot, messageChatID int64) {
	arr, err := getListAbiturient()
	if err != nil {
		log.Printf("Error!!!: %v %s", err, err.(*errors.Error).ErrorStack())
		return
	}

	t := template.New("abiturients list")

	t, err = t.Parse(`{{range .}}<pre>{{if eq .Fio "` + nameFindAbiturient + `"}}>>>{{else}}{{if .Original}} * {{else}}   {{end}}{{end}}{{printf "%3d" .Num}} {{printf "%40s" .Fio}} {{index .Points 4|printf "%3d"}}{{if eq .Fio "` + nameFindAbiturient + `"}}<<<{{else}}{{if .Original}} * {{else}}   {{end}}{{end}}</pre>{{printf "\n"}}{{end}}`)
	if err != nil {
		log.Printf("Error getStatus!!!: %v, %s", err, err.(*errors.Error).ErrorStack())
		return
	}

	for inter := 0; inter < len(arr)/20+1; inter++ {
		var b bytes.Buffer
		last := inter*20 + 20
		if last >= len(arr) {
			last = len(arr)
		}
		err = t.Execute(&b, arr[inter*20:last])
		if err != nil {
			log.Printf("Error!!!: %v, %s", err, err.(*errors.Error).ErrorStack())
			return
		}

		text := b.String()
		if text == "" {
			text = "Empty :-("
		}
		msg := telegram.NewMessage(messageChatID, text)
		msg.ParseMode = "html"
		//				msg.ReplyToMessageID = update.Message.ID

		_, err := bot.SendMessage(msg)
		if err != nil {
			log.Printf("Error!!!: %v, %s", err, err.(*errors.Error).ErrorStack())
			return
		}
	}
}

func tgBotCommandStatWithGetStatus(bot *telegram.Bot, messageChatID int64) {
	status, err := getStatusAbiturient()
	if err != nil {
		log.Printf("Error getStatus!!!: %v, %s", err, err.(*errors.Error).ErrorStack())
		return
	}
	tgBotCommandStat(bot, messageChatID, status)
}

func tgBotCommandStat(bot *telegram.Bot, messageChatID int64, status *StatusAbiturienta) {

	t := template.New("abiturients status")

	t, err := t.Parse("Абитуриент *" + nameFindAbiturient + "*\nПерсональный рейтинг: *{{.Num}}*\nПерсональный рейтинг по оригиналам: *{{.NumWithOriginal}}*")
	if err != nil {
		log.Printf("Error: %v, %s", err, err.(*errors.Error).ErrorStack())
		return
	}

	var b bytes.Buffer
	err = t.Execute(&b, status)
	if err != nil {
		log.Printf("Error: %v, %s", err, err.(*errors.Error).ErrorStack())
		return
	}

	text := b.String()
	msg := telegram.NewMessage(messageChatID, text)
	msg.ParseMode = "markdown"
	//			msg.ReplyToMessageID = update.Message.ID

	_, err = bot.SendMessage(msg)
	if err != nil {
		log.Printf("Error: %v, %s", err, err.(*errors.Error).ErrorStack())
		return
	}
}

func tgBotCommandSubscribe(bot *telegram.Bot, config *ConfigType, messageChatID int64) {

	config.Add(messageChatID)
	config.WriteConfig()

	text := "Subscribed"
	msg := telegram.NewMessage(messageChatID, text)
	msg.ParseMode = "html"
	//			msg.ReplyToMessageID = update.Message.ID

	_, err := bot.SendMessage(msg)
	if err != nil {
		log.Printf("Error: %v, %s", err, err.(*errors.Error).ErrorStack())
		return
	}
}

func tgBotCommandUnSubscribe(bot *telegram.Bot, config *ConfigType, messageChatID int64) {

	delete(config.chats, messageChatID)
	config.WriteConfig()

	text := "UnSubscribed"
	msg := telegram.NewMessage(messageChatID, text)
	msg.ParseMode = "html"
	//			msg.ReplyToMessageID = update.Message.ID

	_, err := bot.SendMessage(msg)
	if err != nil {
		log.Printf("Error: %v, %s", err, err.(*errors.Error).ErrorStack())
		return
	}
}

func tgBotCommandPing(bot *telegram.Bot, messageChatID int64) {

	text := "pong"
	msg := telegram.NewMessage(messageChatID, text)
	msg.ParseMode = "html"
	//			msg.ReplyToMessageID = update.Message.ID

	_, err := bot.SendMessage(msg)
	if err != nil {
		log.Printf("Error: %v, %s", err, err.(*errors.Error).ErrorStack())
		return
	}
}

func tgBotCommandSendChangeStatus(bot *telegram.Bot, config *ConfigType) {
	status, err := getStatusAbiturient()
	if err != nil {
		log.Printf("Error: %v, %s", err, err.(*errors.Error).ErrorStack())
		return
	}
	if (status.Num == config.status.Num) && (status.NumWithOriginal == config.status.NumWithOriginal) {
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

func main() {

	env := os.Getenv("TGBOT_KEY")

	config := new(ConfigType)
	config.ReadConfig()

	bot, err := telegram.New(env)
	if err != nil {
		log.Printf("Error: %v, %s", err, err.(*errors.Error).ErrorStack())
		return
	}

	//bot.Debug = true

	log.Printf("Authorized on account %s", bot.Username)

	updatesParameters := &telegram.GetUpdatesParameters{
		Offset:  0,
		Limit:   100,
		Timeout: 120,
	}

	updates := bot.NewLongPollingChannel(updatesParameters)

	ticker := time.NewTicker(periodUpdate)
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
				tgBotCommandList(bot, messageChatID)
				log.Printf("Out [%s] id:%d text:%s Ok", messageUserName, messageChatID, messageText)
			case (messageText == "/status") || (messageText == "/s"):
				tgBotCommandStatWithGetStatus(bot, messageChatID)
				log.Printf("Out [%s] id:%d text:%s Ok", messageUserName, messageChatID, messageText)
			case (messageText == "/subscribe") || (messageText == "/sub"):
				tgBotCommandSubscribe(bot, config, messageChatID)
				log.Printf("Out [%s] id:%d text:%s Ok", messageUserName, messageChatID, messageText)
			case (messageText == "/unsubscribe") || (messageText == "/uns"):
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
