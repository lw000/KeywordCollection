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

// pcHTML ...
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

	dom.Find(`#content_left`).Find(`[data-click]`).Each(func(i int, s *goquery.Selection) {
		var (
			title string
			href  string
			// cmatchid string
			exists bool
		)
		// // 忽略广告结果
		// cmatchid, exists = s.Attr("cmatchid")
		// if cmatchid != "" {
		// 	title = s.Find(`a`).Text()
		//
		// 	href, exists = s.Find(`a`).Attr("data-landurl")
		// 	if exists {
		// 	}
		//
		// 	title = strings.ReplaceAll(strings.TrimSpace(title), "\n", "")
		// 	href = strings.ReplaceAll(strings.TrimSpace(href), "\n", "")
		// 	if href != "" {
		// 		href = helper.MatchUrl(href)
		// 	}
		// 	m := make(map[string]interface{})
		// 	m["title"] = title
		// 	m["href"] = href
		// 	fn(m)
		// }

		tpl, exists := s.Attr("tpl")
		if !exists {
			return
		}

		switch tpl {
		// 歌手
		case "singer_v2":
		//
		case "exactqa":
		//
		case "fraudphone":
			// 标题
			title = s.Find(`.op_fraudphone_net`).Text()
			// 链接
			href = s.Find(`.result-right`).Find(`.c-showurl`).Text()
		case "se_com_default":
			// 标题
			title = s.Find(`.t>a`).Text()
			// 链接
			cShowURL := s.Find(`div[class="f13"]`).Find(`a[class="c-showurl"]`)
			href = cShowURL.Text()
			if href == "" || helper.MatchUrl(href) == "" {
				href, exists = cShowURL.Attr("href")
				if exists {
					if href != "" {
						href = getRealURL(pc.word.Device, href)
					}
				}
			}
			// 	// TODO:百科内容
			// case "bk_polysemy":
			// 	// 标题
			// 	title = s.Find(`h3[class="t c-gap-bottom-small"]>a`).Text()
			// 	// 链接
			// 	href = s.Find(`.c-showurl`).Text()
			// 	if href == "" || helper.MatchUrl(href) == "" {
			// 		href, exists = s.Find(`h3[class="t c-gap-bottom-small"]>a`).Attr("href")
			// 		if exists {
			// 			if href != "" {
			// 				href = getRealURL(pc.word.Device, href)
			// 			}
			// 		}
			// 	}
			// 	// TODO:百度经验
			// case "jingyan_summary":
			// 	// 标题
			// 	title = s.Find(`.t>a`).Text()
			// 	// 链接
			// 	href, exists = s.Find(`h3[class="t c-gap-bottom-small"]>a`).Attr("href")
			// 	if exists {
			// 		if href != "" {
			// 			href = getRealURL(pc.word.Device, href)
			// 		}
			// 	}
			// }

			title = strings.ReplaceAll(strings.TrimSpace(title), "\n", "")
			href = strings.ReplaceAll(strings.TrimSpace(href), "\n", "")
			m := make(map[string]interface{})
			m["title"] = title
			m["href"] = href
			fn(m)
		}
	})
}

func (pc *pcHTML) nextPageURL(s *goquery.Selection) (*url.URL, error) {
	href, exists := s.Find(`div[id="page"]`).Find(`a[class="n"]`).Attr("href")
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
