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

// TODO: sogou·移动端搜索
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
	body, err := ioutil.ReadFile("./html/sogou/百人牛牛_mobile_1.html")
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

func (mob *mobileSearch) SearchChrome() error {
	// mob.test()
	// return nil

	// 主页-搜索
	var err error
	err = mob.WebBrowser().Get(mob.Opt().Domain)
	if err != nil {
		log.Error(err)
		return err
	}

	var keywordElem selenium.WebElement
	keywordElem, err = mob.WebBrowser().FindElement(selenium.ByID, "keyword")
	if err != nil {
		log.Error(err)
		return err
	}

	err = keywordElem.SendKeys(mob.Word().Keyword)
	if err != nil {
		log.Error(err)
		return err
	}

	err = keywordElem.SendKeys(selenium.ReturnKey)
	if err != nil {
		log.Error(err)
		return err
	}

	time.Sleep(time.Millisecond * mob.Opt().Delay)

	// 点击下一页
	var currentPage = 1

	for currentPage < mob.Word().Page {
		currentPage++

		var ajaxNextPageElem selenium.WebElement
		ajaxNextPageElem, err = mob.WebBrowser().FindElement(selenium.ByID, "ajax_next_page")
		if err != nil {
			log.Error(err)
			break
		}

		err = ajaxNextPageElem.Click()
		if err != nil {
			log.Error(err)
			break
		}

		// log.WithFields(log.Fields{"engine": mob.Opt().Engine, "device": mob.Opt().Type, "currentPage": currentPage}).Info("网页获取")

		time.Sleep(time.Millisecond * mob.Opt().Delay)
	}

	// 获取网页内容
	body, err := mob.WebBrowser().PageSource()
	if err != nil {
		log.Error(err)
		return err
	}

	html := NewMobileHTML(currentPage, mob.Word())
	html.AddHtmls(string(body))
	err = parseserv.ParseServer().AddTask(html)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
