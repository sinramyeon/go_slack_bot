package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"slack_test/crawling"
	"slack_test/envsetting"
	"slack_test/util"
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

	userID := ev.Msg.User

	receivedMsg := ev.Msg.Text

	// 다른 채널에 쳤을때
	if ev.Channel != s.channelID {
		log.Printf("다른 채널 : %s %s", ev.Channel, s.channelID)
		return nil
	}

	log.Println(userID, " : ", receivedMsg)

	if strings.Contains(receivedMsg, `¯\_(ツ)_/¯`) {
		if strings.Contains(ev.Msg.Username, "go돌이") {
			log.Println("봇이 한 대화라 무시 했어요.")
			return nil
		}

		s.client.PostMessage(ev.Channel, `¯\_(ツ)_/¯`, slack.PostMessageParameters{})
	}

	// 봇에게 한 멘션이 아닐 때
	if !(strings.HasPrefix(receivedMsg, fmt.Sprintf("<@%s> ", s.botID))) {

		// 봇이 한 말이면 무시하자!
		if strings.Contains(ev.Msg.Username, "go돌이") {
			log.Println("봇이 한 대화라 무시 했어요.")
			return nil
		}

		// 1. 기사 찾기
		if strings.Contains(receivedMsg, "기사") || strings.Contains(receivedMsg, "뉴스") || strings.Contains(receivedMsg, "소식") {

			log.Println("기사 크롤링 시.")
			m := crawling.NewsScrape()

			if !(len(m) == 0) {

				PostMessage(m, s, ev, "cc1512")

			} else {
				s.client.PostMessage(ev.Channel, "알 수 없는 에러가 발생했습니다. 다시 시도해 주세요.", slack.PostMessageParameters{})
			}
		}

		// 2. 오키 게시글 찾기
		if strings.Contains(receivedMsg, "오키") || strings.Contains(receivedMsg, "옼희") {

			log.Println("오키 크롤링 시.")
			m := crawling.OkkyScrape()

			if len(m) == 0 {
				s.client.PostMessage(ev.Channel, "알 수 없는 에러가 발생했습니다. 다시 시도해 주세요.", slack.PostMessageParameters{})
			} else {

				PostMessage(m, s, ev, "104293")
			}

		}

		// 3. 블로그 찾기
		if strings.Contains(receivedMsg, "블로그") {

			log.Println("블로그 크롤링 시.")

			m := crawling.RssScrape()

			if len(m) == 0 {
				s.client.PostMessage(ev.Channel, "알 수 없는 에러가 발생했습니다. 다시 시도해 주세요.", slack.PostMessageParameters{})
			} else {

				PostMessage(m, s, ev, "2a4f2e")
			}
		}

		// 4. 트위터 찾기
		if strings.Contains(receivedMsg, "트윗") || strings.Contains(receivedMsg, "트위터") {

			log.Println("트위터 크롤링 시.")

			m := crawling.TwitterScrape(tweetenv)

			if len(m) == 0 {
				s.client.PostMessage(ev.Channel, "알 수 없는 에러가 발생했습니다. 다시 시도해 주세요.", slack.PostMessageParameters{})
			} else {

				PostMessage(m, s, ev, "42c7d6")
			}
		}

		// 5. git 사용자이름 입력 시, 오늘의 깃허브 커밋여부 반환
		if strings.HasPrefix(receivedMsg, "git") {
			GitCommitMessage(receivedMsg, s, ev)
		}

		// 6. 근무자 입력 시, 현재 슬랙에 로그인 해 있는 상태인 사용자 반환
		if strings.Contains(receivedMsg, "근무자") {

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
		}

		// 7. 날씨 알려주기
		if strings.Contains(receivedMsg, "날씨") {

			log.Println("날씨 확인 시")
			w := util.GetHourWeather()
			if len(w) == 0 {
				s.client.PostMessage(ev.Channel, "알 수 없는 에러가 발생했습니다. 다시 시도해 주세요.", slack.PostMessageParameters{})
			} else {
				PostMessage(w, s, ev, "2c99ce")
			}

		}

		// 8. 도움 입력 시, 도움말을 전송
		if strings.Contains(receivedMsg, "도움") {

			GetHelp(s, ev)

		}

		// 9. 매일 아침 무료 책을 전송
		if strings.Contains(receivedMsg, "책") || strings.Contains(receivedMsg, "무료") || strings.Contains(receivedMsg, "공짜") {

			Freebook(s, ev.Channel)

		}

		return nil

	}

	// 봇에게 멘션 했을 시
	if strings.HasPrefix(receivedMsg, fmt.Sprintf("<@%s> ", s.botID)) {

		log.Println("봇에게 멘션했을 시.")

		// 봇이 한 말이면 무시하자!
		if strings.Contains(ev.Msg.Username, "go돌이") {
			log.Println("봇이 한 대화라 무시 했어요.")
			return nil
		}

		// select 메뉴
		if strings.Contains(receivedMsg, "도움") {

			/*

				attachment := slack.Attachment{

					Text:       "무엇을 도와드릴까요? :newspaper: ",
					Color:      "#f9a41b",
					CallbackID: "news",
					Actions: []slack.AttachmentAction{

						{

							Name: actionSelect,
							Type: "select",

							Options: []slack.AttachmentActionOption{

								{
									Text:  "IT 기사 읽기",
									Value: "ITNews",
								},
								{
									Text:  "OKKY 읽기",
									Value: "OKKY",
								},
								{
									Text:  "TWITTER 읽기",
									Value: "TWITTER",
								},
								{
									Text:  "기술 블로그 읽기",
									Value: "BLOG",
								},
								{
									Text:  "도움말",
									Value: "HELP",
								},
							},
						},
					},
				}

				params := slack.PostMessageParameters{

					Attachments: []slack.Attachment{
						attachment,
					},
				}

				if _, _, err := s.client.PostMessage(ev.Channel, "", params); err != nil {
					return fmt.Errorf("failed to post message: %s", err)
				}
			*/

			GetHelp(s, ev)

		} else {
			s.client.PostMessage(ev.Channel, "무엇을 도와드릴까요? 도움, 도움말 이라고 입력해보세요~", slack.PostMessageParameters{})
		}

		/*

			 else if strings.Contains(receivedMsg, "점심") {

				log.Println("오늘 점심은 뭘로 드실래요?")

				attachment := slack.Attachment{

					Text:       "오늘의 점심",
					Color:      "#f9a41b",
					CallbackID: "button",
					Actions: []slack.AttachmentAction{

						{
							Name:  "lunch",
							Text:  "한솥",
							Type:  "button",
							Value: "hansot",
						},
						{
							Name:  "lunch",
							Text:  "샐러디",
							Type:  "button",
							Value: "salady",
						},
						{
							Name:  "lunch",
							Text:  "따로 먹을래요",
							Type:  "button",
							Value: "myown",
							Style: "danger",
							Confirm: &slack.ConfirmationField{

								Title:       "오늘은 따로 드시겠어요?",
								Text:        "도시락 멤버에서 빼 드립니다.",
								OkText:      "그래",
								DismissText: "아니",
							},
						},
					},
				}

				params := slack.PostMessageParameters{

					Attachments: []slack.Attachment{
						attachment,
					},
				}

				if _, _, err := s.client.PostMessage(ev.Channel, "", params); err != nil {
					return fmt.Errorf("failed to post message: %s", err)
				}

			}
		*/

	}

	log.Println("return nil")
	return nil
}
