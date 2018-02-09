package crawling

import (
	"encoding/xml"
	"log"
	"net/http"
	env "slackbot/envsetting"
	"strings"

	"github.com/parnurzeal/gorequest"
	"github.com/sclevine/agouti"
)

// 이벤트 게시글 json parsing	용
type Write struct {
	Day        string
	AuthorName string
	Text       string
}

type EntryData struct {
	Key   string `xml:"name,attr"`
	Value string `xml:"text"`
}

type ViewEntry struct {
	Key   string      `xml:"unid,attr"`
	Value []EntryData `xml:"entrydata"`
}
type ViewEntries struct {
	XMLName     xml.Name    `xml:viewentries`
	ViewEntries []ViewEntry `xml:"viewentry"`
}

// 이벤트 얻어오기
func GetEvent() map[string]string {

	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()

	token := MakeToken()

	// 사내 이벤트 게시판 xml로 들어가서 파징
	parsed := new(ViewEntries)
	_, body, _ := gorequest.New().Get(
		"",
	).Type("xml").AddCookie(
		&http.Cookie{Name: "", Value: token},
	).End()

	_ = xml.Unmarshal([]byte(body), &parsed)

	// 결과 정리
	var event Write
	var eventlist []Write

	for _, v := range parsed.ViewEntries {
		var entrydata []EntryData
		entrydata = v.Value

		for key, val := range entrydata {

			if event.AuthorName != "" && event.Day != "" && event.Text != "" {
				eventlist = append(eventlist, event)
				event.AuthorName = ""
				event.Day = ""
				event.Text = ""
			}

			switch key {
			case 1:
				event.Day = val.Value
			case 2:
				event.AuthorName = val.Value
			case 3:
				event.Text = val.Value
			}

		}
	}

	// 그 중 최신 3개만 가져옴
	returnlist := make(map[string]string)
	var loop = 0

	for _, v := range eventlist {
		if loop < 3 {
			returnlist[v.Text] = v.Day + " " + v.AuthorName
			loop++
		}
	}

	return returnlist
}

//ltpa 토큰 만들기
func MakeToken() string {

	//agouti 이용. chromedriver, phantomjs가 %PATH%에 있거나
	//mac인경우에는 brew로 설치 필요!

	driver := agouti.ChromeDriver()
	if err := driver.Start(); err != nil {
		log.Fatalln(err)
	}

	defer recover()
	defer driver.Stop()

	page, err := driver.NewPage(agouti.Browser("phantomjs"))
	if err != nil {
		log.Fatalln(err)
	}

	// 접속 (진짜 크롬 창이 뜸)
	if err := page.Navigate(""); err != nil {
		log.Fatalln(err)
	}

	// 로그인
	var interenv env.Interpark
	interenv = env.InterparkLogin(interenv)
	ID := interenv.ID
	PW := interenv.PW
	page.FindByID("").SendKeys(ID)
	page.FindByID("").SendKeys(PW)

	page.FindByClass("").Click()

	// 이벤트 주소 접속
	if err := page.Navigate(""); err != nil {
		log.Fatalln(err)
	}

	// 쿠키 얻기
	cookie, err := page.GetCookies()
	if err != nil {
		log.Fatalln(err)
	}

	// 토큰 추출
	for _, v := range cookie {
		if strings.Contains(v.Name, "") {
			return v.Value
		}
	}

	return ""

}
