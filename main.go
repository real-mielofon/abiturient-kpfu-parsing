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
	"gitlab.com/toby3d/telegram"
	"golang.org/x/net/html/charset"
)

const (
	urlList = "https://abiturient.kpfu.ru/entrant/abit_entrant_originals_list?p_open=&p_typeofstudy=1&p_faculty=47&p_speciality=1085&p_inst=0&p_category=1"
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

func getListAbiturient() ([]Abiturient, error) {
	var cl = &http.Client{}
	cl.Timeout = 60 * time.Second

	resp, err := cl.Get(urlList)
	if err != nil {
		log.Println("HTTP error:", err)
		return nil, fmt.Errorf("HTTP error: %v", err)
	}

	defer resp.Body.Close()
	// вот здесь и начинается самое интересное
	utf8, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		log.Println("Encoding error:", err)
		return nil, fmt.Errorf("HTTP error: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(utf8)
	if err != nil {
		log.Println("goquery.NewDocumentFromReader:", err)
		return nil, fmt.Errorf("goquery.NewDocumentFromReader: %v", err)
	}

	//fmt.Printf("doc: %+v \n", doc.Text())

	trs := doc.Find("tbody").Eq(1).Find("tr")
	len := trs.Length()
	arr := make([]Abiturient, len-1)

	fmt.Printf("len: %+v \n", len)
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
			log.Printf("strconv.ParseInt(%s), %v\n", s, err)
			return nil, fmt.Errorf("strconv.ParseInt(tds.Eq(0).Text(): %v", err)
		}
		ab.Num = int(num)

		ab.Fio = tds.Eq(1).Text()

		s = tds.Eq(2).Text()
		num, err = strconv.ParseInt(s, 10, 32)
		if err != nil {
			log.Printf("strconv.ParseInt(%s), %v\n", s, err)
			ab.Points[0] = 0
		} else {
			ab.Points[0] = int(num)
		}

		num, err = strconv.ParseInt(tds.Eq(3).Text(), 10, 32)
		if err != nil {
			log.Printf("strconv.ParseInt(%s), %v\n", s, err)
			ab.Points[1] = 0
		} else {
			ab.Points[1] = int(num)
		}

		num, err = strconv.ParseInt(tds.Eq(4).Text(), 10, 32)
		if err != nil {
			log.Printf("strconv.ParseInt(%s), %v\n", s, err)
			ab.Points[2] = 0
		} else {
			ab.Points[2] = int(num)
		}

		num, err = strconv.ParseInt(tds.Eq(5).Text(), 10, 32)
		if err != nil {
			log.Printf("strconv.ParseInt(%s), %v\n", s, err)
			ab.Points[3] = 0
		} else {
			ab.Points[3] = int(num)
		}

		num, err = strconv.ParseInt(tds.Eq(6).Text(), 10, 32)
		if err != nil {
			log.Printf("strconv.ParseInt(%s), %v\n", s, err)
			ab.Points[4] = ab.Points[0] + ab.Points[1] + ab.Points[2] + ab.Points[3]
		} else {
			ab.Points[4] = int(num)
		}

		arr[i-1] = ab

		log.Printf("%4d %40s %3d %s\n", ab.Num, ab.Fio, ab.Points[4], strconv.FormatBool(ab.Original))
	}
	return arr, nil
}

func main() {

	env := os.Getenv("TGBOT_KEY")

	bot, err := telegram.New(env)
	if err != nil {
		log.Panic(err)
	}

	//bot.Debug = true

	log.Printf("Authorized on account %s", bot.Username)

	var updatesParameters telegram.GetUpdatesParameters
	updatesParameters.Timeout = 60

	updates := bot.NewLongPollingChannel(&updatesParameters)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.Username, update.Message.Text)
		switch update.Message.Text {
		case "/list":
			arr, err := getListAbiturient()
			if err != nil {
				log.Printf("Error!!!: %v", err)
			}

			t := template.New("fieldname example")
			t, err = t.Parse(`{{range .}}<pre>{{if eq .Fio "Пономарев Степан Алексеевич"}}>>>{{else}}{{if .Original}} * {{else}}   {{end}}{{end}}{{printf "%3d" .Num}} {{printf "%40s" .Fio}} {{index .Points 4|printf "%3d"}}{{if eq .Fio "Пономарев Степан Алексеевич"}}<<<{{else}}{{if .Original}} * {{else}}   {{end}}{{end}}</pre>{{printf "\n"}}{{end}}`)
			if err != nil {
				log.Panic(err)
			}

			for inter := 0; inter < len(arr)/20+1; inter++ {
				var b bytes.Buffer
				last := inter*20 + 20
				if last >= len(arr) {
					last = len(arr)
				}
				err = t.Execute(&b, arr[inter*20:last])
				if err != nil {
					log.Panic(err)
				}

				text := b.String()
				msg := telegram.NewMessage(update.Message.Chat.ID, text)
				msg.ParseMode = "html"
				msg.ReplyToMessageID = update.Message.ID

				_, err := bot.SendMessage(msg)
				if err != nil {
					log.Panic(err)
				}
			}
		}
	}
}
