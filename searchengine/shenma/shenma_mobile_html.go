package shenma

import (
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

	// dom.Find(`div[data-reco],div[ad_aid]`).Each(func(i int, s *goquery.Selection) {
	// 忽略广告结果
	dom.Find(`div[data-reco]`).Each(func(i int, s *goquery.Selection) {
		var (
			title string
			href  string
		)
		_, exists := s.Attr(`ad_aid`)
		if exists {
			// 标题
			title = s.Find(`a[class="c-title cpc-two-line cpc-title"]`).Find(`span[click_area="title"]`).Text()
			// 链接地址
			href, _ = s.Find(`a[class="c-title cpc-two-line cpc-title"]`).Attr("href")
		} else {
			// 标题
			headerTitle := s.Find(`.c-header-inner`).Find(`.c-header-title`)
			title = headerTitle.Find(`span`).Text()
			if title == "" {
				script := headerTitle.Find(`script[type="text/mask"]`)
				title = script.Text()
				title = deutf8(title)
				title = strings.ReplaceAll(title, "<em>", "")
				title = strings.ReplaceAll(title, "</em>", "")
			}
			// 链接地址
			href, _ = s.Find(`a[class="c-header-inner c-flex-1"]`).Attr("href")
		}

		title = strings.ReplaceAll(strings.TrimSpace(title), "\n", "")
		href = strings.ReplaceAll(strings.TrimSpace(href), "\n", "")

		values := make(map[string]interface{})
		values["title"] = title
		values["href"] = href
		fn(values)
	})
}

// 标题
func (mob *mobileHTML) Title(s *goquery.Selection) string {
	return ""
}

// 详情
func (mob *mobileHTML) Detail(s *goquery.Selection) string {
	return ""
}

// 链接标题
func (mob *mobileHTML) HrefTitle(s *goquery.Selection) string {
	return ""
}

// 链接地址
func (mob *mobileHTML) Href(s *goquery.Selection) string {
	return ""
}
