package so360

import (
	"KeywordCollection/chrome"
	"KeywordCollection/searchengine"
	parseserv "KeywordCollection/server/parseserver"
	log "github.com/sirupsen/logrus"
	"github.com/tebeka/selenium"
	"io/ioutil"
	"time"
)

// TODO: so360·移动端搜索
// mobileSearch ...
type mobileSearch struct {
	searchengine.SearchEngine
}

// NewMobileSearch ...
func NewMobileSearch(chrome *chrome.ChromeDriver, opt *searchengine.SearchOption) *mobileSearch {
	p := &mobileSearch{}
	p.SetOpt(opt)
	p.SetChrome(chrome)
	return p
}

func (mob *mobileSearch) Search() error {
	return nil
}

func (mob *mobileSearch) test() {
	body, err := ioutil.ReadFile("./html/so360/德州扑克_mobile_1.html")
	if err != nil {
		log.Error(err)
		return
	}

	// 解析网页
	html := NewMobileHTML(1, mob.Word())
	html.AddHtmls(string(body))
	if err = parseserv.ParseServer().AddTask(html); err != nil {
		log.Error(err)
		return
	}
}

func (mob *mobileSearch) SearchChrome() error {
	// mob.test()
	// return nil

	// 主页-搜索
	err := mob.WebBrowser().Get(mob.Opt().Domain)
	if err != nil {
		log.Error(err)
		return err
	}

	time.Sleep(time.Millisecond * mob.Opt().Delay)

	qElem, err := mob.WebBrowser().FindElement(selenium.ByID, "q")
	if err != nil {
		log.Error(err)
		return err
	}

	err = qElem.SendKeys(mob.Word().Keyword)
	if err != nil {
		log.Error(err)
		return err
	}

	time.Sleep(time.Millisecond * mob.Opt().Delay)

	searchBtnElem, err := mob.WebBrowser().FindElement(selenium.ByCSSSelector, "form>.search-btn")
	if err != nil {
		log.Error(err)
		return err
	}

	err = searchBtnElem.Click()
	if err != nil {
		log.Error(err)
		return err
	}

	time.Sleep(time.Millisecond * mob.Opt().Delay)

	// 下一页
	var currentPage = 1
	for currentPage < mob.Word().Page {
		currentPage++

		var loadMoreElem selenium.WebElement
		loadMoreElem, err = mob.WebBrowser().FindElement(selenium.ByCSSSelector, "div>.load-more")
		if err != nil {
			log.Error(err)
			break
		}

		if err = loadMoreElem.Click(); err != nil {
			log.Error(err)
			break
		}

		// log.WithFields(log.Fields{"engine": mob.Opt().Engine, "device": mob.Opt().Type, "currentPage": currentPage}).Info("网页获取")

		time.Sleep(time.Millisecond * mob.Opt().Delay)
	}

	body, err := mob.WebBrowser().PageSource()
	if err != nil {
		log.Error(err)
		return err
	}

	html := NewMobileHTML(currentPage, mob.Word())
	html.AddHtmls(string(body))
	if err = parseserv.ParseServer().AddTask(html); err != nil {
		log.Error(err)
		return err
	}

	return nil
}
