package sogou

import (
	"KeywordCollection/helper"
	"KeywordCollection/searchengine"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"strings"
)

type pcHTML struct {
	currentPage int
	word        *searchengine.SearchWord // 搜索关键字
	htmls       []string
}

func NewPcHTML(currentPage int, word *searchengine.SearchWord) *pcHTML {
	return &pcHTML{
		currentPage: currentPage,
		word:        word,
	}
}

func (pc *pcHTML) Htmls() []string {
	return pc.htmls
}

func (pc *pcHTML) AddHtmls(htmls ...string) {
	pc.htmls = append(pc.htmls, htmls...)
}

func (pc *pcHTML) Page() int {
	return pc.currentPage
}

func (pc *pcHTML) Word() *searchengine.SearchWord {
	return pc.word
}

func (pc *pcHTML) Parse(dom *goquery.Document, fn func(values map[string]interface{})) {
	defer func() {
		log.WithFields(log.Fields{"engine": pc.word.Engine, "device": pc.word.Device, "word": pc.word.Keyword}).Info("数据分析完成")
	}()
	log.WithFields(log.Fields{"currentPage": pc.currentPage, "engine": pc.word.Engine, "device": pc.word.Device, "word": pc.word.Keyword}).Info("数据分析中")

	// dom.Find(`.sponsored`).Find(`.biz_sponsor`).Find(`div[class="biz_rb "]`).Each(func(i int, s *goquery.Selection) {
	// 	title := s.Find(`h3[class="biz_title"]`).Find(`a`).Text()
	// 	if title == "" {
	// 		// title = s.Find(`h3[class="vrTitle"]`).Find(`a`).Text()
	// 	}
	//
	// 	// href, exists := s.Find(`h3[class="biz_title"]>a`).Attr("href")
	// 	// if exists {
	// 	// 	if href != "" {
	// 	// 		href = getRealURL(constant.DevicePc, href, func(body string) string {
	// 	// 			dom1, er := goquery.NewDocumentFromReader(strings.NewReader(body))
	// 	// 			if er != nil {
	// 	// 				log.Error(er)
	// 	// 				return ""
	// 	// 			}
	// 	// 			url, exists1 := dom1.Find(`[http-equiv]`).Attr("content")
	// 	// 			if exists1 {
	// 	// 				ss := strings.Split(url, "=")
	// 	// 				if len(ss) > 1 {
	// 	// 					return ss[1]
	// 	// 				}
	// 	// 			}
	// 	// 			return ""
	// 	// 		})
	// 	// 	}
	// 	// }
	// 	href := s.Find(`a[class="cite"]`).Text()
	//
	// 	title = strings.ReplaceAll(strings.TrimSpace(title), "\n", "")
	// 	href = strings.ReplaceAll(strings.TrimSpace(href), "\n", "")
	//
	// 	values := make(map[string]interface{})
	// 	values["title"] = title
	// 	values["href"] = href
	// 	fn(values)
	// })

	dom.Find(`.results`).Find(`div[class="vrwrap"],div[class="rb"]`).Each(func(i int, s *goquery.Selection) {
		_, exists := s.Attr("id")
		if exists {
			return
		}

		title := pc.Title(s)

		// href, exists := s.Find(`h3[class="pt"],[class="vrTitle"]`).Find(`a`).Attr("href")
		// if exists {
		// 	if href != "" {
		// 		href = "https://www.sogou.com" + href
		// 		// href = getRealUrlAddress("pc", href, func(body string) string {
		// 		// 	dom1, er := goquery.NewDocumentFromReader(strings.NewReader(body))
		// 		// 	if er != nil {
		// 		// 		log.Println(er)
		// 		// 		return ""
		// 		// 	}
		// 		// 	url, exists1 := dom1.Find(`[http-equiv]`).Attr("content")
		// 		// 	if exists1 {
		// 		// 		ss := strings.Split(url, "=")
		// 		// 		if len(ss) > 1 {
		// 		// 			return ss[1]
		// 		// 		}
		// 		// 	}
		// 		// 	return ""
		// 		// })
		// 	}
		// }

		href := pc.Href(s)

		values := make(map[string]interface{})
		values["title"] = title
		values["href"] = href
		fn(values)
	})
}

// 标题
func (pc *pcHTML) Title(s *goquery.Selection) string {
	title := s.Find(`h3[class="pt"],[class="vrTitle"]`).Find(`a`).Text()
	if title == "" {
		title = s.Find(`div[class="vrTitle"]`).Text()
	}
	if title != "" {
		title = strings.ReplaceAll(strings.TrimSpace(title), "\n", "")
	}

	return title
}

// 详情
func (pc *pcHTML) Detail(s *goquery.Selection) string {

	return ""
}

// 链接标题
func (pc *pcHTML) HrefTitle(s *goquery.Selection) string {

	return ""
}

// 链接地址
func (pc *pcHTML) Href(s *goquery.Selection) string {
	href := s.Find(`div[class="fb"]>cite`).Text()
	if href != "" {
		href = helper.MatchUrl(href)
	}
	return href
}
