package main

import (
	"fmt"
	"log"
	"slack_test/crawling"
	"slack_test/envsetting"
	"slack_test/util"
	"strings"

	"github.com/nlopes/slack"
)

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

		// 라. 블로그 입력 시(RSS)

		if strings.Contains(receivedMsg, "블로그") {

			log.Println("블로그 크롤링 시.")

			m := crawling.RssScrape()

			if len(m) == 0 {
				s.client.PostMessage(ev.Channel, "알 수 없는 에러가 발생했습니다. 다시 시도해 주세요.", slack.PostMessageParameters{})
			} else {

				PostMessage(m, s, ev, "2a4f2e")
			}
		}

		// 3. 트위터 찾기

		if strings.Contains(receivedMsg, "트윗") || strings.Contains(receivedMsg, "트위터") {

			log.Println("트위터 크롤링 시.")

			m := crawling.TwitterScrape(tweetenv)

			if len(m) == 0 {
				s.client.PostMessage(ev.Channel, "알 수 없는 에러가 발생했습니다. 다시 시도해 주세요.", slack.PostMessageParameters{})
			} else {

				PostMessage(m, s, ev, "42c7d6")
			}
		}

		// 4. git 사용자이름 입력 시, 오늘의 깃허브 커밋여부 반환

		if strings.HasPrefix(receivedMsg, "git") {
			GitCommitMessage(receivedMsg, s, ev)
		}

		// 5. 근무자 입력 시, 현재 슬랙에 로그인 해 있는 상태인 사용자 반환

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

		// 6. 도움 입력 시, 도움말을 전송

		if strings.Contains(receivedMsg, "도움") {

			GetHelp(s, ev)

		}

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

// 시간별로 채널에 메세지 보내기
func (s *SlackListener) PostByTime(env envsetting.EnvConfig) {

	// 정확히 n시 0분 0초가 딱 정시 되는 순간 작동!
	for n := range util.GetHour().C {

		day := util.GetDay()
		hour, _, _ := n.Clock()

		switch hour {

		case 9:

			Freebook(s, env.ChannelID)

		case 10:

			if strings.Contains(day, "Monday") {
				s.client.PostMessageTo(env.ChannelID, "수습평가표 쓰세요", "U6DKDJMPV", slack.PostMessageParameters{})
			}

			if strings.Contains(day, "Thursday") {
				s.client.PostMessageTo(env.ChannelID, "주간보고서 쓰세요", "U6DKDJMPV", slack.PostMessageParameters{})
			}

		case 12:

			if !util.GetWeekends() {
				PostTimeMessage(s, env, "a470e0", "점심알림", "점심 식사 하시러 갈 시간입니다!", "오늘도 맛있는 점심 되세요.")
			}

			// 시간별 커밋 알림봇 구현

			/*
				case 14:

					b, _ := util.GetGitCommit("hero0926")
					if !b {
						attachment := slack.Attachment{

							Color:      "#635129",
							AuthorName: "Commit-bot",
							Title:      "아직 한 커밋이 없어요!",
						}
						params := slack.PostMessageParameters{
							Attachments: []slack.Attachment{
								attachment,
							},
						}

						//제가 새로 만든 유저에게 멘션을 보내는 메서드(풀 리퀘스트는 받아질 것인가?)
						//사용법 (보낼 채널, 보낼 텍스트, 보낼 유저(아이디), 파라미터)
						//그냥 쓰시려면 s.client.PostMessage(env.ChannelID, "<@유저아이디> ", params)
						//꼭 <> 를 넣어줘야 가더라고요...
						//s.client.PostMessageTo(env.ChannelID, "", "U6DKDJMPV", params)
						/*
							func (api *Client) PostMessageTo(channel, text string, id string, params PostMessageParameters) (string, string, error) {
								respChannel, respTimestamp, _, err := api.SendMessageContext(
									context.Background(),
									channel,
									MsgOptionText("<@"+id+"> "+text, params.EscapeText),
									MsgOptionAttachments(params.Attachments...),
									MsgOptionPostMessageParameters(params),
								)
								return respChannel, respTimestamp, err
							}


						//또는 디엠을 보내고 싶을때는 채널명에 유저ID를 쓰시면 됩니다.
						s.client.PostMessage("U6DKDJMPV", "", params)
					}
				case 15:
					b, _ := util.GetGitCommit("hero0926")
					if !b {
						attachment := slack.Attachment{

							Color:      "#633f29",
							AuthorName: "Commit-bot",
							Title:      "아직도! 한 커밋이 없어요!",
						}
						params := slack.PostMessageParameters{
							Attachments: []slack.Attachment{
								attachment,
							},
						}
						s.client.PostMessage("U6DKDJMPV", "", params)
					}
				case 16:
					b, _ := util.GetGitCommit("hero0926")
					if !b {
						attachment := slack.Attachment{

							Color:      "#632b29",
							AuthorName: "Commit-bot",
							Title:      "아직!!!!!!! 한개도 커밋이 없어요!",
						}
						params := slack.PostMessageParameters{
							Attachments: []slack.Attachment{
								attachment,
							},
						}
						s.client.PostMessage("U6DKDJMPV", "", params)
					}
				case 17:
					b, _ := util.GetGitCommit("hero0926")
					if !b {
						attachment := slack.Attachment{

							Color:      "#680e0e",
							AuthorName: "Commit-bot",
							Title: `Commit-bot is watching your commit...
							PLZ commit soon...(아직도 안했다는 소리이다.)`,
						}
						params := slack.PostMessageParameters{
							Attachments: []slack.Attachment{
								attachment,
							},
						}
						s.client.PostMessage("U6DKDJMPV", "", params)
					}

				case 18:

					b, c := util.GetGitCommit("hero0926")

					if !b {

						attachment := slack.Attachment{

							Color:      "#ff0033",
							AuthorName: "긴급 알림",
							Title:      "퇴근 할 시간인데도 커밋을 하지 않았습니다!",
							Text:       "뭔가 하고 가시던지 집에 가서 해보세요!",
						}
						params := slack.PostMessageParameters{
							Attachments: []slack.Attachment{
								attachment,
							},
						}

						s.client.PostMessage("U6DKDJMPV", "", params)

					} else {

						attachment := slack.Attachment{

							Color:      "#ff0033",
							AuthorName: "수고의 알림",
							Title:      "퇴근 할 시간입니다!",
							Text: `오늘도 수고하셨어요. ` +
								"오늘은" + fmt.Sprint(c) + "개의 커밋을 하였습니다.",
						}
						params := slack.PostMessageParameters{
							Attachments: []slack.Attachment{
								attachment,
							},
						}

						s.client.PostMessage("U6DKDJMPV", "", params)
					}
			*/
		case 18:

			if !util.GetWeekends() {
				PostTimeMessage(s, env, "ff0033", "퇴근알림", "퇴근 할 시간입니다!", "오늘도 수고하셨어요.")
			}

			// 야근봇 구현
			// 퇴근 후 일정시간 자동 백업 등을 수행할 수 있을 것 같음...

			/*
				case 19, 20, 21:

					Users, _ := s.client.GetUsers()
					var logineduser []string

					for _, v := range Users {
						if v.Presence == "active" && v.IsBot == false {
							logineduser = append(logineduser, v.Name)
						}
					}

					PostTimeMessage(s, env, "63294e", "아직 불철주야 일하고 계신 분", strings.Join(logineduser, "\n"), "님께서"+fmt.Sprint(hour)+"시까지 수고해주시고 계십니다.")
			*/
		}
	}
}

// 봇 답장용 메서드

func PostMessage(m map[string]string, s *SlackListener, ev *slack.MessageEvent, color string) {

	for k, v := range m {

		attachment := slack.Attachment{

			Color: "#" + color,
			Title: k,
			Text:  v,
		}

		params := slack.PostMessageParameters{

			Attachments: []slack.Attachment{
				attachment,
			},
		}

		s.client.PostMessage(ev.Channel, "", params)

	}

}

func PostTimeMessage(s *SlackListener, env envsetting.EnvConfig, color string, authorname string, title string, text string) {

	attachment := slack.Attachment{

		Color:      "#" + color,
		AuthorName: authorname,
		Title:      title,
		Text:       text,
	}
	params := slack.PostMessageParameters{
		Attachments: []slack.Attachment{
			attachment,
		},
	}
	s.client.PostMessage(env.ChannelID, "", params)
}

// 깃 커밋 확인

func GitCommitMessage(receivedMsg string, s *SlackListener, ev *slack.MessageEvent) {

	log.Println("깃 커밋 확인 시.")
	id := receivedMsg[strings.Index(receivedMsg, " ")+1:]
	util.TrimTrim(id)

	// 사용자가 커밋을 하지 않았을 경우

	b, c := util.GetGitCommit(id)

	if !b {

		if c == 1 {

			s.client.PostMessage(ev.Channel, "그런 유저가 없어요...", slack.PostMessageParameters{})

		} else {

			attachment := slack.Attachment{

				Color:     "#e20000",
				Title:     id + "님께서는 아직 커밋하신 적이 없습니다!",
				TitleLink: "https://github.com/" + id,
				Text:      "내용을 확인 해 주세요",
			}

			params := slack.PostMessageParameters{

				Attachments: []slack.Attachment{
					attachment,
				},
			}

			s.client.PostMessage(ev.Channel, "", params)

		}
	} else {

		attachment := slack.Attachment{

			Color:     "#e20000",
			Title:     id + "님께서는 오늘 " + fmt.Sprint(c) + "개의 커밋을 했습니다!",
			TitleLink: "https://github.com/" + id,
			Text:      "앞으로도 수고해 주세요",
		}

		params := slack.PostMessageParameters{

			Attachments: []slack.Attachment{
				attachment,
			},
		}

		s.client.PostMessage(ev.Channel, "", params)
	}

}

// 도움말

func GetHelp(s *SlackListener, ev *slack.MessageEvent) {
	log.Println("도움말!")

	attachment := slack.Attachment{

		Color: "#296346",
		Title: "봇 사용 커맨드",
		Text: `안녕하세요? IT봇입니다.
				IT봇 사용을 위해서 참고해주세요~
				1. 기사, 뉴스, 소식 키워드 입력 시 오늘의 IT 뉴스라인을 보실 수 있습니다.
				2. 오키, 옼희 입력 시 오키 주간 기술 트렌드를 보실 수 있습니다.
				3. 블로그 입력 시 엄선된 기술블로그들의 rss 피드를 얻어옵니다.
				4. 트위터, 트윗 입력 시 엄선된 트위터를 크롤링해 옵니다.
				5. git 사용자id(Ex - git hero0926) 입력 시 오늘의 커밋상황을 안내해 드립니다.
				6. 근무자 입력 시 현재 슬랙에 로그인 해 있는 사용자를 안내해 드립니다.
				7. 무료, 공짜, 책 입력 시 오늘의 packt 사 무료 ebook을 알려드립니다.`,
	}
	params := slack.PostMessageParameters{
		Attachments: []slack.Attachment{
			attachment,
		},
	}
	s.client.PostMessage(ev.Channel, "", params)
}

// 무료 책

func Freebook(s *SlackListener, channel string) {

	attachment := slack.Attachment{

		Color:     "#92f1f4",
		Title:     crawling.PacktFreeBook(),
		TitleLink: "https://www.packtpub.com/packt/offers/free-learning",
		Text:      "안녕하세요! 오늘의 무료 책을 득템하세요!",
	}

	params := slack.PostMessageParameters{

		Attachments: []slack.Attachment{
			attachment,
		},
	}

	s.client.PostMessage(channel, "", params)

}
