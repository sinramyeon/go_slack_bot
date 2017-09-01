package main

import (
	"fmt"
	"log"
	"slackbot/crawling"
	"slackbot/envsetting"
	"slackbot/util"
	"strings"

	"github.com/nlopes/slack"
)

// 시간별로 채널에 메세지 보내기
func (s *SlackListener) PostByTime(env envsetting.EnvConfig) {

	// 정확히 n시 0분 0초가 딱 정시 되는 순간 작동!
	for n := range util.GetHour().C {

		day := util.GetDay()
		hour, _, _ := n.Clock()

		// 주말 빼고 전송
		switch hour {

		case 10:

			Freebook(s, env.ChannelID)

			if strings.Contains(day, "Thursday") {
				s.client.PostMessage(env.ChannelID, "", slack.PostMessageParameters{

					Attachments: []slack.Attachment{

						slack.Attachment{
							Title: "주간보고서 제출일",
							Text:  "오늘 오전에는 주간 보고서를 제출해 주세요",
							Color: "#f45642",
						},
					},
				})
			}

		case 12:

			if !util.GetWeekends() {
				PostTimeMessage(s, env, "a470e0", "점심알림", "점심 식사 하시러 갈 시간입니다!", "오늘도 맛있는 점심 되세요.")
			}

		case 17:

			if !util.GetWeekends() {
				w := util.TommorrowWeather()

				for k, v := range w {

					attachment := slack.Attachment{

						Color: "#99d5cf",
						Title: k,
						Text:  v,
					}

					params := slack.PostMessageParameters{

						Attachments: []slack.Attachment{
							attachment,
						},
					}

					s.client.PostMessage(env.ChannelID, "", params)

				}
			}

		case 18:

			if !util.GetWeekends() {

				PostTimeMessage(s, env, "ff0033", "퇴근알림", "퇴근 할 시간입니다!", "오늘도 수고하셨어요.")

			}

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

// 시간별 메시지 전송
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
	attachment := slack.Attachment{

		Color: "#296346",
		Title: "봇 사용 커맨드",
		Text: `go돌이 사용설명서

		go돌이에게 멘션을 보내 보세요!

			[날씨] 삼성동 오늘, 내일, 모레 날씨를 알려드려요
			[먼지] 미세먼지 상태나, [자외선]
			[뉴스, 기사, 소식] IT뉴스라인을 보여드려요
			[오키, 옼희] 오키 주간 기술 트렌드를 안내해 드려요
			[블로그] 기술블로그들의 rss 피드를 얻어옵니다
			[트위터, 트윗] 엄선된 트위터를 크롤링해 옵니다
			[근무자] 입력 시 현재 슬랙에 로그인 해 있는 사용자를 안내해 드립니다
			[무료, 공짜, 책] 오늘의 packt 사 무료 ebook을 놓치지 마세요
			[행사, 이벤트] 인터파크 임직원분들만을 위한 행사소식을 알려드립니다
						
			매일 아침 10시에는 packt사 무료 ebook을, 12시에는 점심을, 저녁에는 퇴근시간과 내일 날씨를 알려 드립니다.
			`,
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
