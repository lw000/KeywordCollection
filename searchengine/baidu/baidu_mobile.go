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

// TODO: 百度·移动端搜索
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

// Search ...
func (mob *mobileSearch) Search() error {

	return nil
}

func (mob *mobileSearch) test() {
	body, err := ioutil.ReadFile("./html/baidu/百人牛牛_mobile_1.html")
	if err != nil {
		log.Error(err)
		return
	}

	// 解析网页
	html := NewMobileHTML(1, mob.Word())
	html.AddHtmls(string(body))
	err = parseserv.ParseServer().AddTask(html)
	if err != nil {
		log.Error(err)
		return
	}
}

// SearchChrome ...
func (mob *mobileSearch) SearchChrome() error {
	// mob.test()
	// return nil

	var err error
	// 主页-搜索
	err = mob.WebBrowser().Get(mob.Opt().Domain)
	if err != nil {
		log.Error(err)
		return err
	}

	indexKwElem, err := mob.WebBrowser().FindElement(selenium.ByID, "index-kw")
	if err != nil {
		log.Error(err)
		return err
	}

	err = indexKwElem.SendKeys(mob.Word().Keyword)
	if err != nil {
		log.Error(err)
		return err
	}

	time.Sleep(time.Millisecond * mob.Opt().Delay)

	err = indexKwElem.SendKeys(selenium.ReturnKey)
	if err != nil {
		log.Error(err)
		return err
	}

	time.Sleep(time.Millisecond * mob.Opt().Delay)

	body, err := mob.WebBrowser().PageSource()
	if err != nil {
		log.Error(err)
		return err
	}

	var currentPage = 1

	// 解析网页
	html := NewMobileHTML(1, mob.Word())
	html.AddHtmls(string(body))

	// 下一页
	for currentPage < mob.Word().Page {
		// 切换到浏览器底部
		for index := 0; index < 10; index++ {
			err = mob.WebBrowser().KeyDown(selenium.PageDownKey)
			if err != nil {
				log.Error(err)
				break
			}
			time.Sleep(time.Millisecond * mob.Opt().Delay)
		}

		var newNextpageOnlyElem selenium.WebElement
		if currentPage == 1 {
			newNextpageOnlyElem, err = mob.WebBrowser().FindElement(selenium.ByCSSSelector, "#page-controller a.new-nextpage-only")
			if err != nil {
				log.Error(err)
				break
			}
		} else {
			newNextpageOnlyElem, err = mob.WebBrowser().FindElement(selenium.ByCSSSelector, "#page-controller a.new-nextpage")
			if err != nil {
				log.Error(err)
				break
			}
		}

		err = newNextpageOnlyElem.Click()
		if err != nil {
			log.Error(err)
			break
		}

		time.Sleep(time.Millisecond * mob.Opt().Delay)

		body, err = mob.WebBrowser().PageSource()
		if err != nil {
			log.Error(err)
			break
		}

		// log.WithFields(log.Fields{"engine": mob.Opt().Engine, "device": mob.Opt().Type, "currentPage": currentPage}).Info("网页获取")

		html.AddHtmls(string(body))

		currentPage++
	}

	err = parseserv.ParseServer().AddTask(html)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
