package baidu

import (
	"KeywordCollection/helper"
	"KeywordCollection/searchengine"
	"errors"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"net/url"
	"strings"
)

// mobileHTML ...
type mobileHTML struct {
	currentPage int
	word        *searchengine.SearchWord // 搜索关键字
	htmls       []string
}

func NewMobileHTML(currentPage int, word *searchengine.SearchWord) *mobileHTML {
	return &mobileHTML{
		currentPage: currentPage,
		word:        word,
	}
}

func (mob *mobileHTML) Htmls() []string {
	return mob.htmls
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

	results := dom.Find(`div[class="search-page"]`).Find(`div[class="se-page-bd "]`).Find(`div[class="results"]`)

	// // 广告结果
	// ecAdResults := results.Find(`div[class="ec_wise_ad"]`).Find(`div[class="ec_ad_results"]`)
	// ecAdResults.Each(func(i int, s *goquery.Selection) {
	// 	s.Find(`div[data-lp][data-rank]`).Each(func(i int, s *goquery.Selection) {
	// 		// 标题
	// 		title := s.Find(`a[class="c-blocka ec_title "]`).Find(`h3[class="c-title c-color-link c-gap-top-small c-gap-bottom-small c-line-clamp3"]`).Text()
	// 		// 链接
	// 		href := s.Find(`div[class="c-showurl c-line-clamp1"]`).
	// 			Find(`span[class="c-showurl"]`).Text()
	//
	// 		title = strings.ReplaceAll(strings.TrimSpace(title), "\n", "")
	// 		href = strings.ReplaceAll(strings.TrimSpace(href), "\n", "")
	// 		if href != "" {
	// 			href = helper.MatchUrl(href)
	// 		}
	//
	// 		value := make(map[string]interface{})
	// 		value["title"] = title
	// 		value["href"] = href
	//
	// 		fn(value)
	// 	})
	// })

	// 普通结果
	results.Find(`div[class="c-result result"]`).Each(func(i int, s *goquery.Selection) {
		tpl, exists := s.Attr("tpl")
		if exists {
			switch tpl {
			case "www_normal":
				fallthrough
			case "h5_mobile":
				fallthrough
			case "vid_pocket":
				fallthrough
				// 图片
			case "image_horizonal_sam_tag":
				fallthrough
				// 先关游戏平台
			case "sigma_celebrity_rela":
				// 标题
				title := s.Find(`header[class="c-gap-bottom-small"]`).Find(`span[class="c-title-text"]`).Text()

				// 链接
				href := s.Find(`div[class="c-showurl c-line-clamp1"]`).Find(`span[class="c-showurl c-footer-showurl"]`).Text()

				title = strings.ReplaceAll(strings.TrimSpace(title), "\n", "")
				href = strings.ReplaceAll(strings.TrimSpace(href), "\n", "")
				if href != "" {
					href = helper.MatchUrl(href)
				}

				value := make(map[string]interface{})
				value["title"] = title
				value["href"] = href
				fn(value)
				// 百科
			case "sg_kg_entity":
				// 标题
				title := s.Find(`section`).Find(`div[class="c-title"]`).Text()
				// 链接
				var href string
				href, exists = s.Find(`section`).Find(`a`).Attr("data-url")
				if exists {

				}

				title = strings.ReplaceAll(strings.TrimSpace(title), "\n", "")
				href = strings.ReplaceAll(strings.TrimSpace(href), "\n", "")
				if href != "" {
					href = helper.MatchUrl(href)
				}

				value := make(map[string]interface{})
				value["title"] = title
				value["href"] = href
				fn(value)
			}
		}
	})
}

func (mob *mobileHTML) firstPageURL(s *goquery.Selection) (*url.URL, error) {
	href, exists := s.Find(`div[class="se-page-controller"]`).
		Find(`div[class="new-pagenav c-flexbox"]`).
		Find(`a[class="new-nextpage-only"]`).Attr("href")
	if !exists {
		return nil, errors.New("href 不存在")
	}

	// nextURL, er := url.ParseRequestURI(href)
	nextURL, er := url.Parse(href)
	if er != nil {
		log.Error(er)
		return nil, er
	}

	return nextURL, nil
}

func (mob *mobileHTML) nextPageURL(s *goquery.Selection) (*url.URL, error) {
	href, exists := s.Find(`div[class="se-page-controller"]`).
		Find(`div[class="new-pagenav c-flexbox"]`).
		Find(`div[class="new-pagenav-right"]`).
		Find(`a[class="new-nextpage"]`).Attr("href")
	if !exists {
		return nil, errors.New("href 不存在")
	}

	nextURL, er := url.ParseRequestURI(href)
	if er != nil {
		log.Error(er)
		return nil, er
	}

	return nextURL, nil
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
