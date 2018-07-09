package main

//Abiturient is struct data of abiturient
type Abiturient struct {
	Num      int
	Fio      string
	Points   [5]int
	Original bool
}

type ListAbiturents struct{
    arr []Abiturient
    }

func (l ListAbiturents) get() []Abiturient {
    if Length(l.arr) < 0 {
        l.fetchListAbiturient()
    }
    return l.arr
}
   

func (l ListAbiturents) fetchListAbiturient() (error) {
	var cl = &http.Client{}
	cl.Timeout = 60 * time.Second

	resp, err := cl.Get(urlList)
	if err != nil {
		return nil, werr.Wrap(err)
	}

	defer resp.Body.Close()
	// вот здесь и начинается самое интересное
	utf8, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		return nil, werr.Wrap(err)
	}

	doc, err := goquery.NewDocumentFromReader(utf8)
	if err != nil {
		return nil, werr.Wrap(err)
	}

	//fmt.Printf("doc: %+v \n", doc.Text())

	trs := doc.Find("tbody").Eq(1).Find("tr")
	len := trs.Length()
	l.arr := make([]Abiturient, len-1)

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
			return nil, werr.Wrap(err)
		}
		ab.Num = int(num)

		ab.Fio = tds.Eq(1).Text()

		s = tds.Eq(2).Text()
		num, err = strconv.ParseInt(s, 10, 32)
		if err != nil {
			ab.Points[0] = 0
		} else {
			ab.Points[0] = int(num)
		}

		num, err = strconv.ParseInt(tds.Eq(3).Text(), 10, 32)
		if err != nil {
			ab.Points[1] = 0
		} else {
			ab.Points[1] = int(num)
		}

		num, err = strconv.ParseInt(tds.Eq(4).Text(), 10, 32)
		if err != nil {
			ab.Points[2] = 0
		} else {
			ab.Points[2] = int(num)
		}

		num, err = strconv.ParseInt(tds.Eq(5).Text(), 10, 32)
		if err != nil {
			ab.Points[3] = 0
		} else {
			ab.Points[3] = int(num)
		}

		num, err = strconv.ParseInt(tds.Eq(6).Text(), 10, 32)
		if err != nil {
			ab.Points[4] = ab.Points[0] + ab.Points[1] + ab.Points[2] + ab.Points[3]
		} else {
			ab.Points[4] = int(num)
		}

		l.arr[i-1] = ab
	}
	return nil
}
