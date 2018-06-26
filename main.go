package main

import (
	"time"
	"fmt"
	"net/http"
	"strconv"

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

type Abiturient struct {
	Num    int
	Fio    string
	Points [5]int
}

func getListAbiturient() ([]Abiturient, error) {
	var cl = &http.Client{}
	cl.Timeout = 60*time.Second

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
		tds := trs.Eq(i).Find("td")
		var ab Abiturient

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
			ab.Points[4] = ab.Points[0] + ab.Points[1] + ab.Points[2]+ ab.Points[3]
		} else {
			ab.Points[4] = int(num)
		}

		arr[i-1] = ab

		log.Printf("%4d %40s %3d\n", ab.Num, ab.Fio, ab.Points[3])
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
			text := ""

			arr, err := getListAbiturient()
			if err != nil {
				log.Printf("Error!!!: %v", err)
				text = fmt.Sprintf("Error!!!: %v", err)
			}
			for _, ab := range arr {
				text = text + fmt.Sprintf("%4d %40s %3d\n", ab.Num, ab.Fio, ab.Points[4])
			}

			msg := telegram.NewMessage(update.Message.Chat.ID, text)
			msg.ReplyToMessageID = update.Message.ID

			bot.SendMessage(msg)
		}
	}
}
