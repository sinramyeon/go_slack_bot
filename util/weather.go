package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

//Disclaimer

//THIS WEATHERAPI IS FOR KOREANS.
//I'LL ADD DarkSkyAPI soon(or later).
//NO ENGLISH COMMENT PASSING BY

// sk플래닛 개발지원센터 주소
// https://developers.skplanetx.com
/* 우선 가입부터 해서 appkey를 발급받으세요!
 */

//간편날씨 주소
const Simpleurl = "http://apis.skplanetx.com/weather/summary?version=1&lat=%s&lon=%s&stnid=%s"

// 시간별 현재날씨 주소
const Nowurl = "http://apis.skplanetx.com/weather/current/hourly?lon=127.02583&village=&county=korea&lat=%s&lon=%s&city=%s&version=1"

// 자외선지수 주소
const UVurl = "http://apis.skplanetx.com/weather/windex/uvindex?version=1&lat=%s&lon=%s"

// 미세먼지 주소
const Pollutionurl = "http://apis.skplanetx.com/weather/dust?version=1&lat=%s&lon=%s"

// Get으로 보낼 리퀘스트에 쓰일 값들
type WeatherRequest struct {
	UserId         string    `json:"x-skpop-userId"`
	AcceptLanguage string    `json:"Accept-Language"`
	Date           time.Time `json:"Date"`
	Accept         string    `json:"Accept"`
	AccessToken    string    `json:access_token`
}

// response에서 받아올 값들

// 1. 간편날씨
type JsonMap map[string]interface{}

type JsonTest struct {
	Weather Summary `json:"weather"`
}

type Summary struct {
	Summary []struct {
		Today            Weather `json:"today"`
		Tomorrow         Weather `json:"tomorrow"`
		DayAfterTomorrow Weather `json:"dayAfterTomorrow`
	} `json:"summary"`
}

type Basic struct {
	Today            []Weather `json:"today"`
	Tomorrow         []Weather `json:"tomorrow"`
	DayAfterTomorrow []Weather `json:"dayAfterTomorrow`
}

type Weather struct {
	Temperature Temperature `json:"temperature"`
	Sky         Sky         `json:"sky"`
}

type Temperature struct {
	Tmax string `json:"tmax"`
	Tmin string `json:"tmin"`
}

type Sky struct {
	Name string `json:"name"`
}

// 2. 시간별 날씨
type HourWeather struct {
	Weather Hourly `json:"weather"`
}

type Hourly struct {
	Hourly []struct {
		Grid        Grid        `json:"grid"`
		Temperature Temperature `json:"temperature"`
		Sky         Sky         `json:"sky"`
	} `json:"hourly"`
}

type Grid struct {
	City    string `json:"city"`
	County  string `json:"county"`
	Village string `json:"village"`
}

// 자외선 지수
type UVWeather struct {
	Weather Windex `json:"weather"`
}

type Windex struct {
	Windex UVindex `json:"wIndex"`
}

type UVindex struct {
	UVindex []struct {
		Grid  Grid `json:"grid"`
		Day00 Day  `json:"day00"`
		Day01 Day  `json:"day01"`
		Day02 Day  `jsonL"day02"`
	}
}

type Day struct {
	ImageUrl string `json:"imageUrl"`
	Index    string `json:"index"`
	Comment  string `json:"comment"`
}

// 미세먼지

type Pollution struct {
	Weather Dust `json:"weather"`
}

type Dust struct {
	Dust []struct {
		Time    string  `json:"timeObservation"`
		Station Station `json:"station"`
		Pm10    Pm10    `json:"pm10"`
	}
}

type Station struct {
	Name string `json:"name"`
}

type Pm10 struct {
	Grade string `json:"grade"`
	Value string `json:"value"`
}

// API 호출함수
func CallAPI() []byte {

	jsonreq, err := json.Marshal(WeatherRequest{

		UserId:         "sh8kim@interpark.com",
		AcceptLanguage: "ko_KR",
		Date:           time.Now(),
		Accept:         "application/json",
		AccessToken:    "",
	})

	if err != nil {
		log.Fatalln(err)
		recover()
	}

	return jsonreq

}

// 날씨 얻기
func ReqWeather(url string, jsonreq []byte) []byte {

	// 1. 한번 url로 값을 전송해 본다.
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonreq))
	req.Header.Set("Content-Type", "application/json")
	// 2. skplanet에서 받은 appKey를 등록한다
	// 주의!!!!!!!!!!!!!!! 제발 appKey를 Github 등에 올리지 마세요!!!!!!!!!!
	// 제발!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	req.Header.Set("appkey", "0b37b2f4-0d8b-3cac-b53c-20218fe07af8")
	client := &http.Client{}

	// 3. 값을 받아온다
	resp, err := client.Do(req)
	if err != nil {
		log.Panicln(err)
	}
	// 항상 Close함수는 defer로 까먹지 말고 만듭니다(자주 까먹음)
	defer resp.Body.Close()

	// 4. 값을 읽고 리턴해 줍니다.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panicln(err)
	}

	return body
}

// 간편날씨 얻기
func GetBasicWeather() map[string]string {

	weatherlist := make(map[string]string)

	// api 연결용 json맵을 만들어 둡니다.
	jsonreq := CallAPI()

	// 1. 서초구 간편날씨
	// 위도, 경도, 관측소 주소 얻는 곳(주의 : 느림) 왜이렇게 느리나요? ㅡㅡ
	// http://minwon.kma.go.kr/main/obvStn.do
	seocho := fmt.Sprintf(Simpleurl, "", "", "401")

	// 서초구 주소로 간편날씨 api를 호출했습니다.
	body := ReqWeather(seocho, jsonreq)

	// 받아온 json값을 처리합니다.
	var weather JsonTest
	if err := json.Unmarshal(body, &weather); err != nil {
		log.Fatal(err)
	}

	// 오늘, 내일, 모레 날씨값에 접속할 수 있습니다.
	for _, v := range weather.Weather.Summary {
		var Today string = `최고 : ` + v.Today.Temperature.Tmax +
			"도, 최저 : " + v.Today.Temperature.Tmin +
			`도이며, 기상 상황은 ` + v.Today.Sky.Name + `과 같겠습니다.`

		var Tomorrow string = `최고 : ` + v.Tomorrow.Temperature.Tmax +
			"도, 최저 : " + v.Tomorrow.Temperature.Tmin +
			`도이며, 기상 상황은 ` + v.Tomorrow.Sky.Name + `과 같겠습니다.`

		var DayAfterTomorrow string = `최고 : ` + v.DayAfterTomorrow.Temperature.Tmax +
			"도, 최저 : " + v.DayAfterTomorrow.Temperature.Tmin +
			`도이며, 기상 상황은 ` + v.DayAfterTomorrow.Sky.Name + `과 같겠습니다.`

		weatherlist["오늘 날씨를 알려드립니다."] = Today
		weatherlist["내일 날씨를 알려드립니다."] = Tomorrow
		weatherlist["미리 모레 날씨를 전해드립니다."] = DayAfterTomorrow
	}

	return weatherlist
}

// 지금 이시간 날씨예보 받기
func GetHourWeather() map[string]string {

	weatherlist := make(map[string]string)

	// api 연결용 json맵을 만들어 둡니다.
	jsonreq := CallAPI()

	// 서초구 시간별 날씨
	// api마다 요구하는 값들이 조금씩 다릅니다... 어디는 관측소 위치, 어디는 위도 경도...
	// 그런데 날씨 api는 보통 위도, 경도를 요구하니 위도 경도로 하시면 될 것 같아요.
	seocho := fmt.Sprintf(Nowurl, "37.508214", "127.056541", "seoul")

	// 서초구 주소로 간편날씨 api를 호출했습니다.
	body := ReqWeather(seocho, jsonreq)

	// 받아온 json값을 처리합니다.
	var hourweather HourWeather
	if err := json.Unmarshal(body, &hourweather); err != nil {
		log.Fatal(err)
	}

	// 지금 날씨값에 접속할 수 있습니다.
	for _, v := range hourweather.Weather.Hourly {
		var where string = v.Grid.City + " " + v.Grid.County + " " + v.Grid.Village
		var weather string = v.Sky.Name + "\n최고온도는 " + v.Temperature.Tmax +
			"도이며, 최저온도는 " + v.Temperature.Tmin + " 도 입니다."

		weatherlist[where] = weather
	}

	return weatherlist
}

// 오늘 자외선 지수 받기
func GetTodayUV() (a string, b string, c string) {

	jsonStr := CallAPI()

	seocho := fmt.Sprintf(UVurl, "37.508214", "127.056541")
	body := ReqWeather(seocho, jsonStr)

	// 받아온 json값을 처리합니다.
	var uv UVWeather
	if err := json.Unmarshal(body, &uv); err != nil {
		log.Fatal(err)
	}

	var where string
	var url string
	var comment string
	for _, v := range uv.Weather.Windex.UVindex {
		where = v.Grid.City + v.Grid.County
		url = v.Day00.ImageUrl
		comment = v.Day00.Comment
	}

	return where, url, comment

}

// 오늘 미세먼지
func GetTodayPollution() map[string]string {

	jsonStr := CallAPI()

	seocho := fmt.Sprintf(Pollutionurl, "37.508214", "127.056541")
	body := ReqWeather(seocho, jsonStr)

	// 받아온 json값을 처리합니다.
	var pollution Pollution
	if err := json.Unmarshal(body, &pollution); err != nil {
		log.Fatal(err)
	}

	Pollution := make(map[string]string)
	for _, v := range pollution.Weather.Dust {

		var today = v.Time

		pos := strings.Index(today, " ")
		today = today[0:pos]

		var str = today + "일 " + v.Station.Name
		Pollution[str] = v.Pm10.Grade + " : " + v.Pm10.Value
	}

	return Pollution
}

func TommorrowWeather() map[string]string {

	weatherlist := make(map[string]string)

	jsonreq := CallAPI()

	// 1. 서초구 간편날씨
	// 위도, 경도, 관측소 주소 얻는 곳(주의 : 느림)
	// http://minwon.kma.go.kr/main/obvStn.do
	seocho := fmt.Sprintf(Simpleurl, "", "", "401")
	body := ReqWeather(seocho, jsonreq)

	var weather JsonTest
	if err := json.Unmarshal(body, &weather); err != nil {
		log.Fatal(err)
	}

	for _, v := range weather.Weather.Summary {

		var Tomorrow string = `최고 : ` + v.Tomorrow.Temperature.Tmax +
			", 최저 : " + v.Tomorrow.Temperature.Tmin +
			`도이며, 기상 상황은 ` + v.Tomorrow.Sky.Name + `과 같겠습니다.`

		weatherlist["내일 날씨를 알려드립니다."] = Tomorrow

	}

	return weatherlist
}
