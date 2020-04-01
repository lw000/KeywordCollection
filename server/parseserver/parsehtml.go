package parseserv

import (
	"KeywordCollection/constant"
	"KeywordCollection/searchengine"
	"KeywordCollection/searchengine/nethtp"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"strings"
)

type HTMLData interface {
	// Add ...
	AddHtmls(htmls ...string)
	// Htmls ...
	Htmls() []string
	// CurrentPage ...
	Page() int
	// Word 关键字信息
	Word() *searchengine.SearchWord
	// Parse ...
	Parse(dom *goquery.Document, fn func(values map[string]interface{}))
}

type HTMLParse interface {
	// Start ...
	Start() error
	// Do ...
	Do() error
	// Word 关键字信息
	Word() *searchengine.SearchWord
	// Store 关键字存储信息
	Store() searchengine.WebSearchStore
	// Close ...
	Close()
}

type htmlParse struct {
	data     HTMLData
	store    searchengine.WebSearchStore
	saveHtml bool
}

func NewHTMLParse(data HTMLData, saveHtml bool) *htmlParse {
	return &htmlParse{
		data:     data,
		saveHtml: saveHtml,
	}
}

func (hp *htmlParse) Start() error {
	hp.store = searchengine.NewSearchWordsStore(hp.data.Word()).Start()
	return nil
}

func (hp *htmlParse) Close() {
	hp.store.Close()
}

func (hp *htmlParse) Store() searchengine.WebSearchStore {
	return hp.store
}

func (hp *htmlParse) Word() *searchengine.SearchWord {
	return hp.data.Word()
}

func (hp *htmlParse) Do() error {
	htmls := hp.data.Htmls()
	for _, body := range htmls {
		dom, err := goquery.NewDocumentFromReader(strings.NewReader(body))
		if err != nil {
			log.WithFields(log.Fields{"word": hp.data.Word()}).Error(err)
			return err
		}

		page := hp.data.Page()
		word := hp.data.Word()

		// 保存网页内容
		if hp.saveHtml {
			nethtp.SaveHTML(constant.GetStoreDir(word.Engine), word.Keyword, word.Device, page, body)
		}

		hp.data.Parse(dom, func(values map[string]interface{}) {
			err = hp.store.Store(values)
			if err != nil {
				log.WithFields(log.Fields{"word": hp.data.Word()}).Error(err)
				return
			}
		})

		// hp.store.Foreach(func(index int, data *searchengine.StoreItem) bool {
		// 	log.WithFields(log.Fields{"Word": hp.data.Word()}).Info(data.String())
		// 	return true
		// })
	}

	return nil
}
