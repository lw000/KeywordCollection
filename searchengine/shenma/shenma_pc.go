package shenma

import (
	"KeywordCollection/chrome"
	"KeywordCollection/constant"
	"KeywordCollection/searchengine"
	"KeywordCollection/searchengine/nethtp"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"github.com/tebeka/selenium"
)

// TODO: 极速版搜索
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
	// body, err := ioutil.ReadFile("./html/shenma/百人牛牛_mobile_1.html")
	// if err != nil {
	// 	log.Error(err)
	// 	return
	// }

	// dom, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	// if err != nil {
	// 	log.Error(err)
	// 	return
	// }
	// pc.parseCenterHTML(1, dom)
}

func (pc *pcSearch) SearchChrome() error {
	// pc.test()
	// return nil

	var (
		err         error
		currentPage = 1
	)

	// 主页-搜索
	err = pc.WebBrowser().Get(pc.Opt().Domain)
	if err != nil {
		log.Error(err)
		return err
	}
	// 切换极速版
	{
		var switchItem selenium.WebElement
		switchItem, err = pc.WebBrowser().FindElement(selenium.ByClassName, `switch-item`)
		if err != nil {
			log.Error(err)
			return err
		}
		err = switchItem.Click()
		if err != nil {
			log.Error(err)
			return err
		}

		time.Sleep(time.Millisecond * pc.Opt().Delay)
	}

	var kwTextElem selenium.WebElement
	kwTextElem, err = pc.WebBrowser().FindElement(selenium.ByCSSSelector, `input[type="text"]`)
	if err != nil {
		log.Error(err)
		return err
	}

	err = kwTextElem.SendKeys(pc.Word().Keyword)
	if err != nil {
		log.Error(err)
		return err
	}

	time.Sleep(time.Millisecond * pc.Opt().Delay)

	btnElem, err := pc.WebBrowser().FindElement(selenium.ByID, "button")
	if err != nil {
		log.Error(err)
		return err
	}

	err = btnElem.Click()
	if err != nil {
		log.Error(err)
		return err
	}

	time.Sleep(time.Millisecond * pc.Opt().Delay)

	// 获取网页内容
	body, err := pc.WebBrowser().PageSource()
	if err != nil {
		log.Error(err)
		return err
	}

	// 解析网页
	// dom, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	// if err != nil {
	// 	log.Error(err)
	// 	return err
	// }

	// pc.parseCenterHTML(currentPage, dom)

	// 保存网页
	nethtp.SaveHTML(constant.ShenmaDir, pc.Word().Keyword, pc.Opt().Type, currentPage, body)

	// 点击下一页
	for {
		if currentPage == pc.Word().Page {
			break
		}
		currentPage++

		var nextPageElem selenium.WebElement
		nextPageElem, err = pc.WebBrowser().FindElement(selenium.ByClassName, "next")
		if err != nil {
			log.Error(err)
			continue
		}

		err = nextPageElem.Click()
		if err != nil {
			log.Error(err)
			continue
		}

		time.Sleep(time.Millisecond * pc.Opt().Delay)

		// 获取网页内容
		body, err = pc.WebBrowser().PageSource()
		if err != nil {
			log.Error(err)
			return err
		}

		// 解析网页
		// dom, err = goquery.NewDocumentFromReader(strings.NewReader(body))
		// if err != nil {
		// 	log.Error(err)
		// 	return err
		// }

		// pc.parseCenterHTML(currentPage, dom)

		// 保存网页
		nethtp.SaveHTML(constant.ShenmaDir, pc.Word().Keyword, pc.Opt().Type, currentPage, body)
	}

	return nil
}

func (pc *pcSearch) centerSearch() {
	param := &url.Values{}
	param.Add("q", pc.Word().Keyword)
	param.Add("page", fmt.Sprintf("%d", pc.Word().Page))
	param.Add("from", "smor")
	param.Add("safe", "1")
	param.Add("snum", "6")
	param.Add("tomode", "center") // 极速模式
	body, err := nethtp.FirstPageRequest(pc.Opt().Type, "http://yz.m.sm.cn", "/s", param, nil)
	if body == "" {
		log.Error(err)
		return
	}

	var currentPage = 1

	// 保存文件
	nethtp.SaveHTML(constant.ShenmaDir, pc.Word().Keyword, pc.Opt().Type, currentPage, body)

	// 解析网页
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		log.Error(err)
		return
	}

	nextURL, err := pc.getFastNextPageURL(dom)
	if err != nil {
		log.Error(err)
		return
	}

	title := dom.Find("title").Text()
	log.Info(title)

	// 解析第一页结果
	// pc.parseCenterHTML(currentPage, dom)

	for {
		if currentPage == pc.Word().Page {
			break
		}
		currentPage++

		body, err = nethtp.NextPageRequest(pc.Opt().Type, "http://yz.m.sm.cn", nextURL, nil)
		if err != nil {
			log.Error(err)
			continue
		}

		// 保存文件
		nethtp.SaveHTML(constant.ShenmaDir, pc.Word().Keyword, pc.Opt().Type, currentPage, body)

		dom, err = goquery.NewDocumentFromReader(strings.NewReader(body))
		if err != nil {
			log.Error(err)
			return
		}

		nextURL, err = pc.getFastNextPageURL(dom)
		if err != nil {
			log.Error(err)
			return
		}

		// pc.parseCenterHTML(currentPage, dom)

		time.Sleep(time.Millisecond * pc.Opt().Delay)
	}
}

func (pc *pcSearch) getFastNextPageURL(dom *goquery.Document) (*url.URL, error) {
	href, exists := dom.Find(`div[class="pager"]`).Find(`a[class="next"]`).Attr("href")
	if !exists {
		return nil, errors.New("href 不存在")
	}

	nextURL, err := url.Parse(href)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return nextURL, nil
}

func (pc *pcSearch) getFastPrevPageURL(dom *goquery.Document) (*url.URL, error) {
	href, exists := dom.Find(`div[class="pager"]`).
		Find(`div[class="prev"]`).Attr("href")
	if !exists {
		return nil, errors.New("href 不存在")
	}

	nextURL, err := url.Parse(href)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return nextURL, nil
}
