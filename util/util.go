package util

import (
	"math/rand"
	"strings"
	"time"
)

// 문자열 슬라이싱용(귀찮음)
func Between(value string, a string, b string) string {
	// Get substring between two strings.
	posFirst := strings.Index(value, a)
	if posFirst == -1 {
		return ""
	}
	posLast := strings.Index(value, b)
	if posLast == -1 {
		return ""
	}
	posFirstAdjusted := posFirst + len(a)
	if posFirstAdjusted >= posLast {
		return ""
	}
	return value[posFirstAdjusted:posLast]
}

func Before(value string, a string) string {
	// Get substring before a string.
	pos := strings.Index(value, a)
	if pos == -1 {
		return ""
	}
	return value[0:pos]
}

func After(value string, a string) string {
	// Get substring after a string.
	pos := strings.LastIndex(value, a)
	if pos == -1 {
		return ""
	}
	adjustedPos := pos + len(a)
	if adjustedPos >= len(value) {
		return ""
	}
	return value[adjustedPos:len(value)]
}

func TrimTrim(s string) string {

	strings.Trim(s, " ")
	strings.TrimLeft(s, " ")
	strings.TrimPrefix(s, " ")
	strings.TrimRight(s, " ")
	strings.TrimSpace(s)

	return s

}

// Rand 용(정말 매우 귀찮기 그지없음)
// map 넣을 시 랜덤 k, v을 반환... 왜 고랭에는 셔플이 없는거지?...
func SelRand(m map[string]string) (k string, v string) {
	i := rand.Intn(len(m))
	for k := range m {
		if i == 0 {
			return k, m[k]
		}
		i--
	}
	panic("never")
}

// 정시 얻기

/*
이걸 활용해서 매일 n시에 기사 크롤링을 해온 후 저장해 뒀다 선별해서 보여줄 수도 있고
이걸 활용해서 매일 n시에 사용자의 작업을 확인한 후 메시지를 보내 줄 수도 있을 것 같음
또는 주변 맛집을 찾아다가 점심시간에 투표 포스팅을 할 수도 있음
*/

func GetHour() *time.Ticker {
	c := make(chan time.Time, 1)
	t := &time.Ticker{C: c}
	go func() {
		for {
			n := time.Now()
			if n.Second() == 0 && n.Minute() == 0 {
				c <- n
			}
			time.Sleep(time.Second)
		}
	}()
	return t
}

// 오늘 날짜 얻기

func GetDay() string {

	n := time.Now()

	return n.Weekday().String()

}

// 오늘이 주말인지 평일인지 얻기

func GetWeekends() bool {

	now := time.Now()
	day := now.Weekday().String()

	switch day {

	case "Monday", "Tuesday", "Wednesday", "Thursday", "Friday":
		return false

	case "Saturday", "Sunday":
		return true

	default:
		return false

	}

}
