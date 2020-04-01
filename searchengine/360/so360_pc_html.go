package so360

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

	// 广告结果
	// dom.Find(`#m-spread-left`).Find(`#e_idea_pp>li`).Each(func(i int, s *goquery.Selection) {
	// 	title := pc.Title(s)
	// 	href := pc.Href(s)
	//
	// 	values := make(map[string]interface{})
	// 	values["title"] = title
	// 	values["href"] = href
	// 	fn(values)
	// })

	// 普通结果
	dom.Find(`ul[class="result"]`).Find(`li[class="res-list"]`).Each(func(i int, s *goquery.Selection) {
		title := pc.Title(s)
		href := pc.Href(s)

		values := make(map[string]interface{})
		values["title"] = title
		values["href"] = href
		fn(values)
	})

	// 广告结果
	// dom.Find(`#m-spread-bottom`).Find(`#e_idea_pp_vip_bottom>li`).Each(func(i int, s *goquery.Selection) {
	// 	title := pc.Title(s)
	// 	href := pc.Href(s)
	//
	// 	values := make(map[string]interface{})
	// 	values["title"] = title
	// 	values["href"] = href
	// 	fn(values)
	// })
}

// 标题
func (pc *pcHTML) Title(s *goquery.Selection) string {
	title := s.Find(`a[class="e_haosou_fw_bg_title"]`).Text()
	if title == "" {
		title = s.Find(`h3[class="res-title"],h3[class="res-title "],h3[class="title g-ellipsis"],h3[class="title"]`).Find("a").Text()
	}
	if title == "" {
		title = s.Find(`a`).Text()
	}

	if title != "" {
		return strings.ReplaceAll(strings.TrimSpace(title), "\n", "")
	}
	return ""
}

// 详情
func (pc *pcHTML) Detail(s *goquery.Selection) string {
	detail := s.Find(`div[class="e_haosou_fw_sm"]`).Text()
	if detail == "" {
		detail = s.Find(`div[class="res-comm-con"],p[class="res-desc"],div[class="mh-title"],div[class="res-rich so-rich-news clearfix"]`).Text()
	}
	if detail == "" {
		detail = s.Find(`div[class="e_haosou_fw_sm"]`).Text()
	}

	if detail != "" {
		return strings.ReplaceAll(strings.TrimSpace(detail), "\n", "")
	}
	return ""
}

// 链接标题
func (pc *pcHTML) HrefTitle(s *goquery.Selection) string {
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
func (pc *pcHTML) Href(s *goquery.Selection) string {
	var (
		href string
		// exists bool
	)
	// href, exists = s.Find(`a[class="e_haosou_fw_bg_title"]`).Attr("e-landurl")
	// if exists {
	// 	if href != "" {
	//
	// 	}
	// }
	//
	// if href == "" {
	// 	href, exists = s.Find(`h3[class="res-title"]>a`).Attr("href")
	// 	if exists {
	// 		if href != "" {
	// 		}
	// 	}
	// }
	//
	// if href == "" {
	// 	href, exists = s.Find(`a`).Attr("e-landurl")
	// 	if exists {
	// 		if href != "" {
	//
	// 		}
	// 	}
	// }

	href = s.Find(".res-linkinfo").Find("cite").Text()
	if href == "" {
		href = s.Find(".mh-url").Find("cite").Text()
	}
	if href != "" {
		href = helper.MatchUrl(href)
	}
	return href
}
