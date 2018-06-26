package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
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

func main() {

	resp, err := http.Get(urlList)
	if err != nil {
		fmt.Println("HTTP error:", err)
		return
	}

	defer resp.Body.Close()
	// вот здесь и начинается самое интересное
	utf8, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		fmt.Println("Encoding error:", err)
		return
	}

	doc, err := goquery.NewDocumentFromReader(utf8)
	_check(err)

	//fmt.Printf("doc: %+v \n", doc.Text())

	trs := doc.Find("tbody").Eq(1).Find("tr")
	len := trs.Length()
	fmt.Printf("len: %+v \n", len)
	for i := 1; i < len; i++ {
		tds := trs.Eq(i).Find("td")
		pos, err := strconv.ParseInt(tds.Eq(0).Text(), 10, 64)
		_check(err)

		fio := tds.Eq(1).Text()

		points, err := strconv.ParseInt(tds.Eq(6).Text(), 10, 64)
		_check(err)

		fmt.Printf("%4d %40s %3d\n", pos, fio, points)
	}

}
