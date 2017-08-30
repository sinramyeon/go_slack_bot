package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// sk플래닛 개발지원센터 주소
// https://developers.skplanetx.com
/* 우선 가입부터 해서 appkey를 발급받으세요!
 */

//간편날씨
const Simpleurl = "http://apis.skplanetx.com/weather/summary?version=1&lat=%s&lon=%s&stnid=%s"

// 시간별 현재날씨
const Nowurl = "http://apis.skplanetx.com/weather/current/hourly?lon=127.02583&village=&county=korea&lat=%s&lon=%s&city=%s&version=1"

// 자외선지수
const UVurl = " http://apis.skplanetx.com/weather/windex/uvindex?lon=%s&lat=%s&version=1"

// 미세먼지
const Pollutionurl = " http://apis.skplanetx.com/weather/dust?lon=%s&lat=%s&version=1"

// ------------------------------------------------------------------------------------------------------------ //

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

// ------------------------------------------------------------------------------------------------------------ //

func CallAPI() []byte {
	// 1. api를 호출

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

func ReqWeather(url string, jsonreq []byte) []byte {

	req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonreq))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("appkey", "0b37b2f4-0d8b-3cac-b53c-20218fe07af8")
	client := &http.Client{}

	// 2. api에서 값을 받아옴

	resp, err := client.Do(req)
	if err != nil {
		log.Panicln(err)
	}
	defer resp.Body.Close()

	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		log.Panicln(err2)
	}
	return body
}

func GetBasicWeather() map[string]string {

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

		var Today string = `최저 : ` + v.Today.Temperature.Tmax +
			"최고 : " + v.Today.Temperature.Tmin +
			`도이며, 기상 상황은 ` + v.Today.Sky.Name + `과 같겠습니다.`

		var Tomorrow string = `최저 : ` + v.Tomorrow.Temperature.Tmax +
			"최고 : " + v.Tomorrow.Temperature.Tmin +
			`도이며, 기상 상황은 ` + v.Tomorrow.Sky.Name + `과 같겠습니다.`

		var DayAfterTomorrow string = `최저 : ` + v.DayAfterTomorrow.Temperature.Tmax +
			"최고 : " + v.DayAfterTomorrow.Temperature.Tmin +
			`도이며, 기상 상황은 ` + v.DayAfterTomorrow.Sky.Name + `과 같겠습니다.`

		weatherlist["오늘 날씨를 알려드립니다."] = Today
		weatherlist["내일 날씨를 알려드립니다."] = Tomorrow
		weatherlist["미리 모레 날씨를 전해드립니다."] = DayAfterTomorrow
	}

	return weatherlist
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

		var Tomorrow string = `최저 : ` + v.Tomorrow.Temperature.Tmax +
			"최고 : " + v.Tomorrow.Temperature.Tmin +
			`도이며, 기상 상황은 ` + v.Tomorrow.Sky.Name + `과 같겠습니다.`

		weatherlist["내일 날씨를 알려드립니다."] = Tomorrow

	}

	return weatherlist
}

func GetHourWeather() map[string]string {

	weatherlist := make(map[string]string)

	jsonreq := CallAPI()

	// 2. 서초구 시간별 날씨

	seocho := fmt.Sprintf(Nowurl, "37.508214", "127.056541", "seoul")
	body := ReqWeather(seocho, jsonreq)

	var hourweather HourWeather
	if err := json.Unmarshal(body, &hourweather); err != nil {
		log.Fatal(err)
	}

	for _, v := range hourweather.Weather.Hourly {
		var where string = v.Grid.City + " " + v.Grid.County + " " + v.Grid.Village
		var weather string = v.Sky.Name + "\n최고온도는 " + v.Temperature.Tmax +
			"도이며, 최저온도는 " + v.Temperature.Tmin + " 도 입니다."

		weatherlist[where] = weather
	}

	return weatherlist
}
