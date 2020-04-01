package sogou

import (
	"KeywordCollection/helper"
	"KeywordCollection/searchengine"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"strings"
)

type mobileHTML struct {
	currentPage int
	word        *searchengine.SearchWord // 搜索关键字
	htmls       []string
}

func (mob *mobileHTML) Htmls() []string {
	return mob.htmls
}

func NewMobileHTML(currentPage int, word *searchengine.SearchWord) *mobileHTML {
	return &mobileHTML{
		currentPage: currentPage,
		word:        word,
	}
}

func (mob *mobileHTML) AddHtmls(htmls ...string) {
	mob.htmls = append(mob.htmls, htmls...)
}

func (mob *mobileHTML) Page() int {
	return mob.currentPage
}

func (mob *mobileHTML) Word() *searchengine.SearchWord {
	return mob.word
}

func (mob *mobileHTML) Parse(dom *goquery.Document, fn func(values map[string]interface{})) {
	defer func() {
		log.WithFields(log.Fields{"engine": mob.word.Engine, "device": mob.word.Device, "word": mob.word.Keyword}).Info("数据分析完成")
	}()
	log.WithFields(log.Fields{"currentPage": mob.currentPage, "engine": mob.word.Engine, "device": mob.word.Device, "word": mob.word.Keyword}).Info("数据分析中")

	results := dom.Find(`#mainBodyWapResult`).Find(`#mainBody`).Find(`#resultsWrap`).Find(`.results`)
	// 忽略广告结果
	// results.Find(`.vrResult, .result, .ec_ad_results`).Each(func(i int, s *goquery.Selection) {
	results.Find(`.vrResult, .result, .jsResult`).Each(func(i int, s *goquery.Selection) {
		_, exists := s.Attr("data-v")
		if !exists {
			return
		}

		// 标题
		title := mob.Title(s)
		href := mob.Href(s)
		values := make(map[string]interface{})
		values["title"] = title
		values["href"] = href
		fn(values)
	})
}

// 标题
func (mob *mobileHTML) Title(s *goquery.Selection) string {
	title := s.Find(`h3[class="biz-tit"]`).Find(`a`).Text()

	if title == "" {
		title = s.Find(`h3[class="vr-tit"]`).Find(`a`).Text()
	}

	if title == "" {
		title = s.Find(`h3[class="vr-tit"]`).Find(`span`).Text()
	}

	if title == "" {
		title = s.Find(`h3[class]`).Text()
	}

	title = strings.ReplaceAll(strings.TrimSpace(title), "\n", "")

	return title
}

// 详情
func (mob *mobileHTML) Detail(s *goquery.Selection) string {
	detail := s.Find(`div[class="text-layout"]`).Text()
	if detail == "" {
		detail = s.Find(`div[class="info clamp3"]`).Text()
	}
	detail = strings.ReplaceAll(strings.TrimSpace(detail), "\n", "")
	return detail
}

// 链接标题
func (mob *mobileHTML) HrefTitle(s *goquery.Selection) string {

	return ""
}

// 链接地址
func (mob *mobileHTML) Href(s *goquery.Selection) string {
	var href string
	// hrefSel := s.Find(`.ad_result`).Find(`.citeurl`)
	// hrefSel.RemoveFiltered(`script`)
	// hrefSel.RemoveFiltered(`span`)
	// href = hrefSel.Text()
	if href == "" {
		href = s.Find(`.citeurl`).Text()
	}

	if href != "" {
		href = strings.ReplaceAll(strings.TrimSpace(href), "\n", "")
		href = helper.MatchUrl(href)
	}
	return href
}
