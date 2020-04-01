package shenma

import (
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

	dom.Find(`body`).Find(`div[class="article ali_row"]`).Each(func(i int, s *goquery.Selection) {
		title := s.Find(`h2`).Find(`a`).Text()
		// detail := s.Find(`p`).Text()
		// hreftitle := s.Find(`div[class="other"]`).Text()
		href, _ := s.Find(`h2`).Find(`a`).Attr("href")

		title = strings.ReplaceAll(strings.TrimSpace(title), "\n", "")
		// detail = strings.ReplaceAll(strings.TrimSpace(detail), "\n", "")
		// hreftitle = strings.ReplaceAll(strings.TrimSpace(hreftitle), "\n", "")
		href = strings.ReplaceAll(strings.TrimSpace(href), "\n", "")

		value := make(map[string]interface{})
		value["title"] = title
		value["href"] = href

		values := make(map[string]interface{})
		values["title"] = title
		values["href"] = href
		fn(values)
	})
}

// 标题
func (pc *pcHTML) Title(s *goquery.Selection) string {

	return ""
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
	return ""
}
