package sogou

import (
	"KeywordCollection/chrome"
	"KeywordCollection/searchengine"
	parseserv "KeywordCollection/server/parseserver"
	log "github.com/sirupsen/logrus"
	"github.com/tebeka/selenium"
	"io/ioutil"
	"time"
)

// TODO: sogou·PC端搜索
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

func (pc *pcSearch) Search() error {
	return nil
}

func (pc *pcSearch) test() {
	body, err := ioutil.ReadFile("./html/sogou/百人牛牛_pc_1.html")
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

	queryElem, err := pc.WebBrowser().FindElement(selenium.ByID, "query")
	if err != nil {
		log.Error(err)
		return err
	}
	err = queryElem.SendKeys(pc.Word().Keyword)
	if err != nil {
		log.Error(err)
		return err
	}

	stbElem, err := pc.WebBrowser().FindElement(selenium.ByID, "stb")
	if err != nil {
		log.Error(err)
		return err
	}
	err = stbElem.Click()
	if err != nil {
		log.Error(err)
		return err
	}

	time.Sleep(time.Millisecond * pc.Opt().Delay)

	body, err := pc.WebBrowser().PageSource()
	if err != nil {
		log.Error(err)
		return err
	}

	// 解析网页
	// html := &pcSearchHTML{search: pc, currentPage: currentPage}
	// err = html.Do(string(body))
	// if err != nil {
	// 	log.Error(err)
	// 	return err
	// }

	html := NewPcHTML(currentPage, pc.Word())
	html.AddHtmls(string(body))

	// 点击下一页
	for currentPage < pc.Word().Page {
		currentPage++

		var sogouNextElem selenium.WebElement
		sogouNextElem, err = pc.WebBrowser().FindElement(selenium.ByCSSSelector, "#pagebar_container .np")
		// sogouNextElem, err = webBrowser.FindElement(selenium.ByID, "sogou_next")
		if err != nil {
			log.Error(err)
			return err
		}

		err = sogouNextElem.Click()
		if err != nil {
			log.Error(err)
			return err
		}

		time.Sleep(time.Millisecond * pc.Opt().Delay)

		body, err = pc.WebBrowser().PageSource()
		if err != nil {
			log.Error(err)
			return err
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
