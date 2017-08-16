package main

import (
	"log"
	"net/http"
	"os"
	"slack_test/envsetting"

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
