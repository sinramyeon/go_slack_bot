// 웹에서 기사 가져오기

// 참고
//// rss
/*
https://github.com/mmcdole/gofeed
*/

//// 웹 크롤링
/*
https://github.com/PuerkitoBio/goquery
*/

////

package main

import (
	"log"
	"strings"

	"math/rand"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/mmcdole/gofeed"
)

type slackReturn struct {
	Title string
	URL   string
}

// 트위터 읽기
func TwitterScrape(env twitterConfig) map[string]string {

	tweetlist := make(map[string]string)

	// https://apps.twitter.com/app/14097310/keys 를 보시면 됩니다.

	config := oauth1.NewConfig(env.confKey, env.confSecret)
	token := oauth1.NewToken(env.tokenKey, env.tokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)

	// 테스트 하시려면 유저 확인을 해보시면 좋습니다.

	//verifyParams := &twitter.AccountVerifyParams{
	//	SkipStatus:   twitter.Bool(true),
	//	IncludeEmail: twitter.Bool(true),
	//}

	//로그인 유저 얻기
	//user, _, _ := client.Accounts.VerifyCredentials(verifyParams)
	//fmt.Println("연결완료, %s", user)

	// 특정 계정 크롤링 하기

	tweets := GetUserTweets(2, "golangweekly", client)

	for _, v := range tweets {
		tweetlist[v.Text] = tweetlist[v.Entities.Urls[0].URL]
	}

	tweets = GetUserTweets(2, "WEIRDxMEETUP", client)

	for _, v := range tweets {
		tweetlist[v.Text] = tweetlist[v.Entities.Urls[0].URL]
	}

	tweets = GetUserTweets(2, "devsfarm", client)

	for _, v := range tweets {
		tweetlist[v.Text] = tweetlist[v.Entities.Urls[0].URL]
	}

	// 트윗 내용이 키이고 url이 밸류인 걸로 리턴 했습니다(밸류안에 뭘 넣을지 고민중...)
	return tweetlist

}

// 유저별 트윗을 얻는 모듈
// 얻어오고 싶은 양, 유저 아이디, 클라이언트 접속정보를 넣어 보세요...
func GetUserTweets(many int, id string, client *twitter.Client) []twitter.Tweet {

	tweets, _, _ := client.Timelines.UserTimeline(&twitter.UserTimelineParams{
		Count:           many,
		ScreenName:      id,
		IncludeRetweets: twitter.Bool(false),
	})

	return tweets

}

// 갑자기 null값 보냄
// rss 블로그 읽기
func RssScrape() map[string]string {

	log.Println("RssScrape")

	rssURL := []string{
		"https://charsyam.wordpress.com/feed/",
		"http://j.mearie.org/rss",
		"http://feeds.feedburner.com/theyearlyprophet/GGGO?format=xml",
		"http://rss.egloos.com/blog/kwon37xi",
		"http://feeds.feedburner.com/xguru?format=xml",
		"http://thoughts.chkwon.net/feed/",
		"http://feeds.feedburner.com/goodhyun",
		"http://nolboo.github.io/feed.xml",
		"http://html5lab.kr/feed/",
		"http://www.kmshack.kr/rss",
		"http://rss.egloos.com/blog/minjang",
		"http://bomjun.tistory.com/rss",
		"http://kimbyeonghwan.tumblr.com/rss",
		"http://greemate.tistory.com/rss",
		"http://www.se.or.kr/rss",
		"https://subokim.wordpress.com/feed/",
		"http://blog.seulgi.kim/feeds/posts/default",
		"http://moogi.new21.org/tc/rss",
		"http://knight76.tistory.com/rss",
		"http://blog.rss.naver.com/drvoss.xml",
		"https://kimws.wordpress.com/feed/",
		"http://androidkr.blogspot.com/feeds/posts/default",
		"http://feeds.feedburner.com/crazytazo?format=xml",
		"http://forensic-proof.com/feed",
		"http://feeds.feedburner.com/reinblog",
		"http://www.memoriesreloaded.net/feeds/posts/default",
		"http://rss.egloos.com/blog/agile",
		"http://huns.me/feed",
		"http://taegon.kim/feed",
		"http://feeds.feedburner.com/GaeraeBlog?format=xml",
		"https://beyondj2ee.wordpress.com/feed/",
		"http://androidhuman.com/rss",
		"http://www.mickeykim.com/rss",
		"http://www.gisdeveloper.co.kr/rss",
		"http://rss.egloos.com/blog/greentec",

		"http://www.rkttu.com/atom",
		"http://bugsfixed.blogspot.com/feeds/posts/default",
		"http://occamsrazr.net/tt/index.xml",
		"http://ryulib.tistory.com/rss",
		"http://blog.lael.be/feed",

		"http://hoonsbara.tistory.com/rss",
		"http://agebreak.blog.me/rss",
		"http://likejazz.com/rss",
		"https://sangminpark.wordpress.com/feed/",
		"http://rss.egloos.com/blog/parkpd",
		"http://bagjunggyu.blogspot.com/feeds/posts/default",
		"http://blog.naver.com/pjt3591oo",
		"http://feeds.feedburner.com/junyoung?format=xml",
		"http://feeds.feedburner.com/baenefit/slXh",
		"http://whiteship.me/?feed=rss2",
		"http://blog.daum.net/xml/rss/funfunction",
		"http://feeds.feedburner.com/rss_outsider_dev?format=xml",
		"http://blog.suminb.com/feed.xml",

		"http://gamecodingschool.org/feed/",
		"http://rss.egloos.com/blog/seoz",
		"https://arload.wordpress.com/feed/",
		"http://blog.saltfactory.net/feed",
		"http://emptydream.tistory.com/rss",
		"http://www.talk-with-hani.com/rss",
		"http://feeds.feedburner.com/codewiz",
		"http://zetlos.tistory.com/rss",
		"http://hyeonseok.com/rss/",

		"http://toyfab.tistory.com/rss",
		"http://qnibus.com/feed/",
		"http://blog.rss.naver.com/delmadang.xml",
		"https://only2sea.wordpress.com/feed/",
		"http://kwangshin.pe.kr/blog/feed/",
		"http://www.flowdas.com/blog/feeds/rss/",
		"http://www.enshahar.me/feeds/posts/default",
		"http://yonght.tumblr.com/rss",
		"http://blog.hax0r.info/rss",
		"http://feeds.feedburner.com/channy",
		"http://mobicon.tistory.com/rss",
		"http://changsuk.me/?feed=rss2",
		"https://justhackem.wordpress.com/feed/",
		"http://genesis8.tistory.com/rss",
		"http://www.buggymind.com/rss",
		"http://feeds.feedburner.com/sangwook?format=xml",
		"http://www.shalomeir.com/feed/",
		"http://blog.scaloid.org/feeds/posts/default",
		"http://blog.xcoda.net/rss",
		"http://daddycat.blogspot.com/feeds/posts/default",
		"http://feeds.feedburner.com/pyrasis?format=xml",
		"http://www.jimmyrim.com/rss",

		"http://blog.java2game.com/rss",
		"http://blog.lastmind.net/feed",
		"http://devyongsik.tistory.com/rss",
		"http://openlook.org/wp/feed/",
		"http://feeds.feedburner.com/allofsoftware?format=xml",
		"http://www.php5.me/blog/feed/",
		"http://feeds.feedburner.com/gogamza?format=xml",
		"http://www.moreagile.net/feeds/posts/default",
		"http://blrunner.com/rss",
		"http://rss.egloos.com/blog/benelog",
		"http://www.sysnet.pe.kr/rss/getrss.aspx?boardId=635954948",
		"http://health20.kr/rss",
		"http://bcho.tistory.com/rss",
		"http://sungmooncho.com/feed/",
		"http://blog.kivol.net/rss",
		"http://rss.egloos.com/blog/aeternum",
		"http://softwaregeeks.org/feed/",

		"http://blog.doortts.com/rss",
		"http://javacan.tistory.com/rss",
		"http://jacking.tistory.com/rss",

		"http://feeds.feedburner.com/Smartmob",
		"http://kkamagui.tistory.com/rss",
		"http://blog.kazikai.net/?feed=rss2",
		"https://joone.wordpress.com/feed/",
		"http://blog.dahlia.kr/rss",
		"http://blog.fupfin.com/?feed=rss2",
		"http://xrath.com/feed/",
		"http://pragmaticstory.com/feed/",
		"http://rss.egloos.com/blog/recipes",
		"http://iam-hs.com/rss",

		"http://feeds.feedburner.com/gamedevforever?format=xml",
		"http://d2.naver.com/d2.atom",
		"http://www.nextree.co.kr/feed/",
		"http://blog.dramancompany.com/category/develop/feed/",
		"https://engineering.linecorp.com/ko/blog/rss2",
		"http://tech.lezhin.com/rss/",
		"http://blog.secmem.org/rss",
		"https://spoqa.github.io/rss",
		"https://blogs.idincu.com/dev/feed/",
		"http://dev.rsquare.co.kr/feed/",
		"http://feeds.feedburner.com/acornpub",
		"http://blog.embian.com/rss",
		"http://woowabros.github.io/feed",
		"http://eclipse.or.kr/index.php?title=특수기능:최근바뀜&feed=atom",
		"http://blog.weirdx.io/feed/",
		"http://bigmatch.i-um.net/feed/",
		"http://blog.insightbook.co.kr/rss",
		"http://tech.kakao.com/rss/",
		"http://www.codingnews.net/?feed=rss2",
		"http://www.techsuda.com/feed",
		"http://tmondev.blog.me/rss",
		"http://gameplanner.cafe24.com/feed/",
		"http://feeds.feedburner.com/skpreadme?format=xml",
		"http://engineering.vcnc.co.kr/atom.xml",
		"http://feeds.feedburner.com/GoogleDevelopersKorea?format=xml",
		"http://hacks.mozilla.or.kr/feed/",
	}

	rsslist := make(map[string]string)

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 4; i++ {

		//rssURL 중 하나를 골라다
		choosen := rssURL[rand.Intn(len(rssURL)-1)]

		//파징 하기
		fp := gofeed.NewParser()
		feed, err := fp.ParseURL(choosen)

		log.Println("???????", feed, err)

		rsslist[feed.Items[0].Title] = feed.Items[0].Link
	}

	log.Println("rsslist : ", rsslist)

	return rsslist

}

// OKKY 기술 게시글 읽기
func OkkyScrape() map[string]string {

	doc, _ := goquery.NewDocument("https://okky.kr/")

	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()

	okkylist := make(map[string]string)

	// .Find(".class tag .class")
	// .Find("a").Text() a 태그 안에 든 거

	doc.Find(".article-middle-block").Each(func(i int, s *goquery.Selection) {

		title := s.Find("h5").Text()
		url := "https://okky.kr" + s.Find("h5 a").AttrOr("href", "없음")
		//okkylist = append(okkylist, title, url)

		okkylist[title] = url

	})

	return okkylist

}

// 갑자기 null값 보냄
// 깃허브 고 오픈소스 찾기
func GoScrape() map[string]string {

	// 깃허브에 연ㅇ결

	doc, _ := goquery.NewDocument("https://github.com/trending/go?since=daily")

	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()

	githublist := make(map[string]string)

	// EachWithBreak 메서드의 true/false 리턴으로 검색 중 멈추고 싶을 때 멈출 수 있습니다.
	// 나는 너무많으니까 5개만 가져오는데서 멈춤

	var forLoop int = 0

	doc.Find(".repo-list li").EachWithBreak(func(i int, s *goquery.Selection) bool {

		if forLoop > 4 {
			return false
		} else {
			// Trim 시리즈 적용하는 모듈 나중에 만들어야겠음(너무도 귀찮음)
			title := s.Find("h3 a").AttrOr("href", "없음")
			strings.TrimSpace(title)
			desc := s.Find(".py-1 p").Text()
			strings.TrimSpace(desc)
			strings.TrimLeft(desc, " ")

			githublist[title] = desc
			forLoop++
			return true
		}
	})

	return githublist
}

// IT 뉴스 찾기
func NewsScrape() map[string]string {

	doc, _ := goquery.NewDocument("http://www.itworld.co.kr/")

	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()

	newslist := make(map[string]string)

	var forLoop int = 0

	doc.Find(".cio_summary").EachWithBreak(func(i int, s *goquery.Selection) bool {

		if forLoop > 4 {
			return false
		} else {
			title := s.Find("ul li a").Text()
			strings.TrimSpace(title)
			url := s.Find("ul li a").AttrOr("href", "없음")
			strings.TrimSpace(url)
			strings.TrimLeft(url, " ")

			newslist[title] = url
			forLoop++
			return true
		}
	})
	return newslist
}

// 문자열 슬라이싱용(귀찮음)
func between(value string, a string, b string) string {
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

func before(value string, a string) string {
	// Get substring before a string.
	pos := strings.Index(value, a)
	if pos == -1 {
		return ""
	}
	return value[0:pos]
}

func after(value string, a string) string {
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
