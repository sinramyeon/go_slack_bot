package util

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// 제휴식당 저장용 배열
var Eatoptions = []InterEat{
	InterEat{
		Name:  "포시즌플러스",
		Price: "6000",
		Add:   "",
	},
	InterEat{
		Name:  "본도시락",
		Price: "5500",
		Add:   "",
	},
	InterEat{
		Name:  "김밥천국",
		Price: "5500",
		Add:   "",
	},
	InterEat{
		Name:  "로봇김밥",
		Price: "5500",
		Add:   "서울 강남구 삼성동 156-5",
	},
	InterEat{
		Name:  "오감만족",
		Price: "5500",
		Add:   "서울시 강남구 삼성동 157-10 GS25 편의점 옆",
	},
	InterEat{
		Name:  "우리집밥",
		Price: "6000",
		Add:   "서울특별시 강남구 삼성동 158-4 두산위브센티움1층",
	},
	InterEat{
		Name:  "병천순대",
		Price: "6000",
		Add:   "서울시 강남구 삼성동 158-6번지 지하1층",
	}, InterEat{
		Name:  "태항",
		Price: "6000",
		Add:   "서울특별시 강남구 삼성동 153-23",
	}, InterEat{
		Name:  "우주최강도시락",
		Price: "6000",
		Add:   "서울시 강남구 삼성동 156-5 1층",
	}, InterEat{
		Name:  "명동칼국수",
		Price: "6000",
		Add:   "서울시 강남구 삼성동 129-17",
	}, InterEat{
		Name:  "은성식당",
		Price: "6000",
		Add:   "서울특별시 강남구 삼성동 144-14 채널리저브빌딩 B1",
	}, InterEat{
		Name:  "강남순두부",
		Price: "6000",
		Add:   "서울시 강남구 삼성동 144-14 채널리저브빌딩 B1",
	}, InterEat{
		Name:  "샐러디",
		Price: "6000",
		Add:   "",
	}, InterEat{
		Name:  "예손순두부",
		Price: "7000",
		Add:   "서울특별시 강남구 테헤란로 87길 동성빌딩 지하 1층",
	}, InterEat{
		Name:  "육대장",
		Price: "7000",
		Add:   "",
	},
}

// json 파싱용
type MapRequest struct {
	Query string `json:"query"`
	Key   string `json:"appkey"`
}

type InterEat struct {
	Name  string
	Price string
	Add   string
}

type Documents struct {
	Meta      map[string]interface{} `json:"meta"`
	Documents []struct {
		PlaceName   string `json:"place_name"`
		PlaceURL    string `json:"place_url"`
		Category    string `json:"category_name"`
		Address     string `json:"address_name"`
		RoadAddress string `json:"road_address_name"`
		Phone       string `json:"phone"`
	} `json:"documents"`
}

type EatingPlace struct {
	PlaceURL    string `json:"place_url"`
	Category    string `json:"category_name"`
	Address     string `json:"address_name"`
	RoadAddress string `json:"road_address_name"`
	Phone       string `json:"phone"`
}

// 지도API 호출
func CallMapAPI() []byte {

	jsonreq, err := json.Marshal(MapRequest{
		Key: "1b5396fee020a19230beade9dbee9381",
	})

	if err != nil {
		log.Fatalln(err)
		recover()
	}

	return jsonreq

}

// 식당 검색
func ReqEatMap(where string) map[string]EatingPlace {

	eatspot := make(map[string]EatingPlace)
	var docu Documents

	//API호출
	jsonreq := CallMapAPI()

	//연결준비 (query가 한글이므로 꼭 QueryEscape()필요)
	mapurl := "https://dapi.kakao.com/v2/local/search/keyword.json?query="
	myurl := mapurl + url.QueryEscape(where)

	// 1. 한번 url로 값을 전송해 본다.
	req, err := http.NewRequest("GET", myurl, bytes.NewBuffer(jsonreq))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "KakaoAK 1b5396fee020a19230beade9dbee9381")

	client := &http.Client{}

	// 2. 값을 받아온다
	resp, err := client.Do(req)
	if err != nil {
		log.Panicln(err)
	}
	// 항상 Close함수는 defer로 까먹지 말고 만듭니다(자주 까먹음)
	defer resp.Body.Close()

	// 3. 값을 읽고 리턴해 줍니다.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panicln(err)
	}

	if err := json.Unmarshal(body, &docu); err != nil {
		log.Fatal(err)
	}

	for _, v := range docu.Documents {

		// 삼성, 서초 근처 식당만 반환
		if MapNear(v.Address) {
			eatspot[v.PlaceName] = EatingPlace{
				Address:     v.Address,
				Category:    v.Category,
				Phone:       v.Phone,
				PlaceURL:    v.PlaceURL,
				RoadAddress: v.RoadAddress,
			}

		} else {
			continue
		}

	}

	return eatspot

}

// 제휴식당 주소에 등록되어있는 주소 반환
func EatMap(where string) map[string]EatingPlace {

	eatspot := make(map[string]EatingPlace)

	for _, v := range Eatoptions {
		if strings.Contains(where, v.Name) {
			eatspot[v.Name] = EatingPlace{
				Address: v.Add,
			}
		} else {
			continue
		}
	}

	return eatspot

}

// 회사 근처에 식당이 있는지 확인
func MapNear(where string) bool {

	if strings.Contains(where, "강남") || strings.Contains(where, "서초") {
		return true
	} else {
		return false
	}
}

// 제휴식당 목록 포스팅
func PostEatPlace() map[string]string {

	eatplace := make(map[string]string)

	for _, v := range Eatoptions {
		eatplace[v.Name] = v.Add
	}

	return eatplace

}
