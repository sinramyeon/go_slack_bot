package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"slackbot/crawling"
	"slackbot/envsetting"
	"slackbot/util"
	"strings"

	"github.com/nlopes/slack"
)

const (
	// action is used for slack attament action.
	actionSelect = "select"
	actionStart  = "start"
	actionCancel = "cancel"
	buttonSelect = "button"
)

type SlackListener struct {
	client    *slack.Client
	botID     string
	channelID string
}

func main() {

	os.Exit(_main(os.Args[1:]))

}

func _main(args []string) int {

	defer func() {
		if err := recover(); err != nil {
			log.Fatal(err)
			return
		}
	}()

	// 1. 설정
	var env envsetting.EnvConfig
	env = envsetting.Envconfig(env)
	api := slack.New(env.BotToken)

	var tweetenv envsetting.TwitterConfig
	tweetenv = envsetting.Twitterconfing(tweetenv)

	slackListener := &SlackListener{
		client:    api,
		botID:     env.BotID,
		channelID: env.ChannelID,
	}

	// DEBUG설정 - 개발시에만 켜주세요
	//api.SetDebug(true)
	//로그인 테스트하기
	groups, err := api.GetGroups(false)
	if err != nil {
		log.Printf("%s 로그인 중 에러가 발생하였습니다. : %s\n", groups, err)
		return 0
	}

	// 2. 메시지 받는 설정

	go slackListener.ListenAndResponse(tweetenv)
	go slackListener.PostByTime(env)

	log.Printf("[INFO] Server listening on :%s", env.Port)
	if err := http.ListenAndServe(":"+env.Port, nil); err != nil {
		log.Printf("[ERROR] %s", err)
		return 1
	}
	return 0

}

func (s *SlackListener) ListenAndResponse(tweetenv envsetting.TwitterConfig) {

	rtm := s.client.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {

		// 타입이 특정 인터페이스를 구현하는지 검사
		// interface{}.(타입).(구현하는지 궁금한 인터페이스)
		switch ev := msg.Data.(type) {
		/////////////interface.(type) 형식을 눈에 익혀두자~
		//Data 인터페이스의 type에 따라 switch 문 돌리는 중...
		//slack 의 messageEvent일때 처리
		case *slack.MessageEvent:
			if err := s.handleMessageEvent(ev, tweetenv); err != nil {
				log.Printf("[ERROR] 처리중 에러가 발생하였습니다.: %s", err)
			}
		}
	}
}

//////////////////////////////////////////////////////////////////

// 메시지 받고 보내기
func (s *SlackListener) handleMessageEvent(ev *slack.MessageEvent, tweetenv envsetting.TwitterConfig) error {

	receivedMsg := ev.Msg.Text

	// 다른 채널에 쳤을때
	if ev.Channel != s.channelID {
		log.Printf("다른 채널 : %s %s", ev.Channel, s.channelID)
		return nil
	}

	if strings.Contains(receivedMsg, `¯\_(ツ)_/¯`) {
		if strings.Contains(ev.Msg.Username, "go돌이") {
			return nil
		}

		s.client.PostMessage(ev.Channel, `¯\_(ツ)_/¯`, slack.PostMessageParameters{})
	}

	// 봇에게 한 멘션이 아닐 때
	if !(strings.HasPrefix(receivedMsg, fmt.Sprintf("<@%s> ", s.botID))) {

		// 봇이 한 말이면 무시하자!
		if strings.Contains(ev.Msg.Username, "go돌이") {
			return nil
		}

		// 도움말을 전하자
		if strings.Contains(receivedMsg, "도움") || strings.Contains(receivedMsg, "도와") || strings.Contains(receivedMsg, "help") {

			GetHelp(s, ev)

		}

		return nil

	}

	// 봇에게 멘션 했을 시
	if strings.HasPrefix(receivedMsg, fmt.Sprintf("<@%s> ", s.botID)) {

		// 봇이 한 말이면 무시하자!
		if strings.Contains(ev.Msg.Username, "go돌이") {
			return nil
		}

		// 도움말!
		if strings.Contains(receivedMsg, "도움") || strings.Contains(receivedMsg, "도와") || strings.Contains(receivedMsg, "help") {
			GetHelp(s, ev)
		} else if strings.Contains(receivedMsg, "기사") || strings.Contains(receivedMsg, "뉴스") || strings.Contains(receivedMsg, "소식") {

			log.Println("기사 크롤링 시.")
			m := crawling.NewsScrape()

			if !(len(m) == 0) {

				PostMessage(m, s, ev, "7a7a6a")

			} else {
				s.client.PostMessage(ev.Channel, "알 수 없는 에러가 발생했습니다. 다시 시도해 주세요.", slack.PostMessageParameters{})
			}
		} else if strings.Contains(receivedMsg, "오키") || strings.Contains(receivedMsg, "옼희") {

			log.Println("오키 크롤링 시.")
			m := crawling.OkkyScrape()

			if len(m) == 0 {
				s.client.PostMessage(ev.Channel, "알 수 없는 에러가 발생했습니다. 다시 시도해 주세요.", slack.PostMessageParameters{})
			} else {

				PostMessage(m, s, ev, "104293")
			}

		} else if strings.Contains(receivedMsg, "블로그") {

			log.Println("블로그 크롤링 시.")

			m := crawling.RssScrape()

			if len(m) == 0 {
				s.client.PostMessage(ev.Channel, "알 수 없는 에러가 발생했습니다. 다시 시도해 주세요.", slack.PostMessageParameters{})
			} else {

				PostMessage(m, s, ev, "2a4f2e")
			}
		} else if strings.Contains(receivedMsg, "트윗") || strings.Contains(receivedMsg, "트위터") {

			log.Println("트위터 크롤링 시.")

			m := crawling.TwitterScrape(tweetenv)

			if len(m) == 0 {
				s.client.PostMessage(ev.Channel, "알 수 없는 에러가 발생했습니다. 다시 시도해 주세요.", slack.PostMessageParameters{})
			} else {

				PostMessage(m, s, ev, "42c7d6")
			}
		} else if strings.Contains(receivedMsg, "근무자") {

			log.Println("현재 로그인 해 있는 사용자 확인 시")
			Users, _ := s.client.GetUsers()
			var logineduser []string

			for _, v := range Users {
				if v.Presence == "active" && v.IsBot == false {
					logineduser = append(logineduser, v.Name)
				}
			}

			attachment := slack.Attachment{

				Color: "#292963",
				Title: "현재 로그인 해 있는 사용자",
				Text:  strings.Join(logineduser, "\n"),
			}
			params := slack.PostMessageParameters{
				Attachments: []slack.Attachment{
					attachment,
				},
			}
			s.client.PostMessage(ev.Channel, "", params)
		} else if strings.Contains(receivedMsg, "날씨") {

			log.Println("날씨 확인 시")

			if strings.Contains(receivedMsg, "내일") || strings.Contains(receivedMsg, "모레") {
				w := util.GetBasicWeather()
				if len(w) == 0 {
					s.client.PostMessage(ev.Channel, "알 수 없는 에러가 발생했습니다. 다시 시도해 주세요.", slack.PostMessageParameters{})
				} else {
					PostMessage(w, s, ev, "2c99ce")
				}
			} else {
				w := util.GetHourWeather()
				if len(w) == 0 {
					s.client.PostMessage(ev.Channel, "알 수 없는 에러가 발생했습니다. 다시 시도해 주세요.", slack.PostMessageParameters{})
				} else {
					PostMessage(w, s, ev, "2c99ce")
				}
			}
		} else if strings.Contains(receivedMsg, "먼지") {
			log.Println("미세먼지 확인 시")
			w := util.GetTodayPollution()
			if len(w) == 0 {
				s.client.PostMessage(ev.Channel, "알 수 없는 에러가 발생했습니다. 다시 시도해 주세요.", slack.PostMessageParameters{})
			} else {
				PostMessage(w, s, ev, "2c2d25")
			}
		} else if strings.Contains(receivedMsg, "자외선") || strings.Contains(receivedMsg, "선크림") || strings.Contains(receivedMsg, "태양") {
			log.Println("자외선 확인 시")
			where, url, comment := util.GetTodayUV()
			attachment := slack.Attachment{

				Color:    "#f45942",
				Title:    where,
				Text:     comment,
				ImageURL: url,
			}

			params := slack.PostMessageParameters{

				Attachments: []slack.Attachment{
					attachment,
				},
			}

			s.client.PostMessage(ev.Channel, "나가기 전 자외선 지수 확인하세요.", params)

		} else if strings.Contains(receivedMsg, "책") || strings.Contains(receivedMsg, "무료") || strings.Contains(receivedMsg, "공짜") {
			Freebook(s, ev.Channel)
		} else if strings.Contains(receivedMsg, "이벤트") || strings.Contains(receivedMsg, "행사") {

			log.Println("행사 확인 시")
			w := crawling.GetEvent()
			if len(w) == 0 {
				s.client.PostMessage(ev.Channel, "알 수 없는 에러가 발생했습니다. 다시 시도해 주세요.", slack.PostMessageParameters{})
			} else {
				s.client.PostMessage(ev.Channel, "인터파크 임직원만을 위한 행사를 알려드려요!", slack.PostMessageParameters{})
				PostMessage(w, s, ev, "cc1512")
			}
		} else {
			log.Println("헛소리 했을 시 + ", util.After(receivedMsg, " "))
			s.client.PostMessage(ev.Channel, util.After(receivedMsg, " ")+" 는 무슨 말이예요?!", slack.PostMessageParameters{})
		}

	}

	return nil
}
