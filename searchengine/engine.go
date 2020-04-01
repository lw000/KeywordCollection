package searchengine

import (
	"KeywordCollection/chrome"
	"encoding/json"
	"fmt"
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/tebeka/selenium"
	"time"
)

// SearchKey 关键字
type SearchKey struct {
	Engine string // 搜索引擎
	Device string // 搜索设备
}

// SearchOption 搜索配置
type SearchOption struct {
	Domain string        // 网页地址
	Engine string        // 搜索引擎
	Type   string        // 搜索设备类型
	Delay  time.Duration // 操作延迟时间。单位毫秒(ms)
}

// SearchEngine 搜索引擎
type SearchEngine struct {
	opt        *SearchOption        // 配置
	word       *SearchWord          // 关键字
	webBrowser selenium.WebDriver   // 浏览器
	chrome     *chrome.ChromeDriver // Chrome服务
}

// SearchWordsStore ...
type SearchWordsStore struct {
	wd    *SearchWord // 搜索关键字
	store arraylist.List
	input chan *StoreItem
	done  chan struct{}
}

func (s *SearchEngine) Word() *SearchWord {
	return s.word
}

func (s *SearchEngine) SetWord(word *SearchWord) {
	s.word = word
	if s.word.Page < 1 {
		s.word.Page = 1
	}
}

// Opt ...
func (s *SearchEngine) Opt() *SearchOption {
	return s.opt
}

// SetOpt ...
func (s *SearchEngine) SetOpt(opt *SearchOption) {
	s.opt = opt
	if s.opt.Delay < 1000 {
		s.opt.Delay = 1000
	}
}

// SetChrome ...
func (s *SearchEngine) SetChrome(chrome *chrome.ChromeDriver) {
	s.chrome = chrome
}

// WebBrowser ...
func (s *SearchEngine) WebBrowser() selenium.WebDriver {
	return s.webBrowser
}

// Start 开始
func (s *SearchEngine) Start(word *SearchWord) error {
	var err error
	s.webBrowser, err = s.chrome.OpenWebBrowser()
	if err != nil {
		return err
	}

	s.SetWord(word)

	return nil
}

// Close 关闭
func (s *SearchEngine) Close() {
	if s.webBrowser != nil {
		_ = s.webBrowser.Quit()
	}
}

func (s StoreItem) String() string {
	return fmt.Sprintf("Title:%s, Hyperlink:%s", s.Title, s.Hyperlink)
}

// JSON ...
func (s StoreItem) JSON() string {
	data, err := json.Marshal(s)
	if err != nil {
		return "{}"
	}
	return string(data)
}
