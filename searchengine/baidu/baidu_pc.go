package baidu

import (
	"KeywordCollection/chrome"
	"KeywordCollection/searchengine"
	parseserv "KeywordCollection/server/parseserver"
	"io/ioutil"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/tebeka/selenium"
)

// TODO: 百度·PC搜索
// pcSearch ...
type pcSearch struct {
	searchengine.SearchEngine
}

// NewPcSearch ...
func NewPcSearch(chrome *chrome.ChromeDriver, opt *searchengine.SearchOption) *pcSearch {
	p := &pcSearch{}
	p.SetOpt(opt)
	p.SetChrome(chrome)
	return p
}

func (pc *pcSearch) test() {
	body, err := ioutil.ReadFile("./html/baidu/百人牛牛_pc_1.html")
	if err != nil {
		log.Error(err)
		return
	}

	html := NewPcHTML(1, pc.Word())
	html.AddHtmls(string(body))
	err = parseserv.ParseServer().AddTask(html)
	if err != nil {
		log.Error(err)
		return
	}
}

// Search ...
func (pc *pcSearch) Search() error {

	return nil
}

// SearchChrome ...
func (pc *pcSearch) SearchChrome() error {
	// pc.test()
	// return nil

	var err error
	var currentPage = 1

	// 主页-搜索
	err = pc.WebBrowser().Get(pc.Opt().Domain)
	if err != nil {
		log.Error(err)
		return err
	}

	kwElem, err := pc.WebBrowser().FindElement(selenium.ByID, "kw")
	if err != nil {
		log.Error(err)
		return err
	}

	err = kwElem.SendKeys(pc.Word().Keyword)
	if err != nil {
		log.Error(err)
		return err
	}

	suElem, err := pc.WebBrowser().FindElement(selenium.ByID, "su")
	if err != nil {
		log.Error(err)
		return err
	}
	err = suElem.Click()

	time.Sleep(time.Millisecond * pc.Opt().Delay)

	body, err := pc.WebBrowser().PageSource()
	if err != nil {
		log.Error(err)
		return err
	}

	// 解析网页
	html := NewPcHTML(1, pc.Word())
	html.AddHtmls(string(body))

	// 点击下一页
	for currentPage < pc.Word().Page {
		currentPage++

		var stbElem selenium.WebElement
		stbElem, err = pc.WebBrowser().FindElement(selenium.ByCSSSelector, "#page>a:last-child")
		if err != nil {
			log.Error(err)
			break
		}

		err = stbElem.Click()
		if err != nil {
			log.Error(err)
			break
		}

		time.Sleep(time.Millisecond * pc.Opt().Delay)

		body, err = pc.WebBrowser().PageSource()
		if err != nil {
			log.Error(err)
			break
		}

		// log.WithFields(log.Fields{"engine": pc.Opt().Engine, "device": pc.Opt().Type, "currentPage": currentPage}).Info("网页获取")

		html.AddHtmls(string(body))
	}

	err = parseserv.ParseServer().AddTask(html)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
