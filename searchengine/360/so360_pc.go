package so360

import (
	"KeywordCollection/chrome"
	"KeywordCollection/searchengine"
	parseserv "KeywordCollection/server/parseserver"
	"io/ioutil"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/tebeka/selenium"
)

// TODO: so360·PC端搜索

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
	body, err := ioutil.ReadFile("./html/so360/德州扑克_pc_1.html")
	if err != nil {
		log.Error(err)
		return
	}

	html := NewPcHTML(1, pc.Word())
	html.AddHtmls(string(body))
	if err = parseserv.ParseServer().AddTask(html); err != nil {
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

	inputElem, err := pc.WebBrowser().FindElement(selenium.ByID, "input")
	if err != nil {
		log.Error(err)
		return err
	}

	if err = inputElem.SendKeys(pc.Word().Keyword); err != nil {
		log.Error(err)
		return err
	}

	searchButtonElem, err := pc.WebBrowser().FindElement(selenium.ByID, "search-button")
	if err != nil {
		log.Error(err)
		return err
	}

	if err = searchButtonElem.Click(); err != nil {
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
	html := NewPcHTML(currentPage, pc.Word())
	html.AddHtmls(string(body))

	// 点击下一页
	for currentPage < pc.Word().Page {
		// 切换到浏览器底部
		for index := 0; index < 2; index++ {
			err = pc.WebBrowser().KeyDown(selenium.PageDownKey)
			if err != nil {
				log.Error(err)
				break
			}
			time.Sleep(time.Millisecond * pc.Opt().Delay)
		}

		currentPage++

		var stbElem selenium.WebElement
		stbElem, err = pc.WebBrowser().FindElement(selenium.ByCSSSelector, "#snext")
		if err != nil {
			log.Error(err)
			break
		}

		if err = stbElem.Click(); err != nil {
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

		// 解析HTML
		html.AddHtmls(string(body))
	}

	if err = parseserv.ParseServer().AddTask(html); err != nil {
		log.Error(err)
		return err
	}

	return nil
}
