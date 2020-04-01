package so360

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

	// var e = 1
	// dom.Find(`div[class="r-results"]`).Find(`[data-page]`).Each(func(i int, s *goquery.Selection) {
	dom.Find(`div[class="r-results"]`).Find(`.g-card`).Each(func(i int, s *goquery.Selection) {
		// data-type="0|normal"                             广告
		// data-type="2|normal"                             广告

		// data-cat="mso-svideo"							视频
		// data-cat="mso-baike"								360百科
		// data-cat="mso-news"
		// data-cat="mso-image"                             图片
		// data-cat="mso-recommend-normal-rel-1_bottom"		猜你关注
		// data-cat="mso-recommend-normal-rel-1"			相关游戏
		// data-cat="mso-wenda-stepnew-step"                360问答
		// data-cat="own_guide_recommend"

		// data-mohe-type="360pic"							360图片

		// dataType, exists := s.Attr(`data-type`)
		// if exists {
		// 	log.Info(e, "data-type", dataType)
		// }
		//

		// dataCat, exists := s.Attr(`data-cat`)
		// if exists {
		// 	log.WithFields(log.Fields{"Engine": mob.word.Engine, "Device": mob.word.Device, "data-cat": dataCat, "e": e}).Info(mob.word.Engine)
		// }
		//
		// e++

		// 标题
		title := mob.Title(s)
		// 链接地址
		href := mob.Href(s)

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
	title := s.Find(`div[class="e_idea_list"]>a`).Text()
	if title == "" {
		title = s.Find(`h3[class="res-title"]`).Text()
	}
	if title == "" {
		title = s.Find(`span[class="title"]`).Text()
	}

	if title != "" {
		return strings.ReplaceAll(strings.TrimSpace(title), "\n", "")
	}
	return ""
}

// 详情
func (mob *mobileHTML) Detail(s *goquery.Selection) string {
	detail := s.Find(`div[class="e-base-box"]`).Find(`div[class="e-fw-desc"]`).Text()
	if detail == "" {
		detail = s.Find(`div[class="e_idea_listfelx"]`).Find(`a[class="e-graphic-right"]`).Text()
	}

	if detail == "" {
		detail = s.Find(`.summary`).Text()
	}

	if detail != "" {
		return strings.ReplaceAll(strings.TrimSpace(detail), "\n", "")
	}
	return ""
}

// 链接标题
func (mob *mobileHTML) HrefTitle(s *goquery.Selection) string {
	hreftitle := s.Find(`div[class="e_ad_brand"]>a`).Text()
	if hreftitle == "" {
		hreftitle = s.Find(`div[class="res-supplement"]>cite`).Text()
	}
	if hreftitle == "" {
		hreftitle = s.Find(`div[class="mohe-nav g-mt"]`).Find(`span[class="mohe-site"]`).Text()
	}

	if hreftitle != "" {
		return strings.ReplaceAll(strings.TrimSpace(hreftitle), "\n", "")
	}
	return ""
}

// 链接地址
func (mob *mobileHTML) Href(s *goquery.Selection) string {
	href, exists := s.Attr("data-pcurl")
	if !exists {
		href = s.Find(`div[class="e_ad_brand"]`).Find(`a[class="e_fw_brand_link"]`).Text()
		if href == "" {
			href = s.Find(`.res-supplement`).Find(`cite`).Find(`.res-site-url`).Text()
		}
	}

	if href != "" {
		href = helper.MatchUrl(href)
	}
	return href
}
