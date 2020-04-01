package searchengine

import (
	"encoding/json"
	"errors"
	"fmt"
)

// StoreItem ...
type StoreItem struct {
	Title     string `json:"title"`     // 标题
	Hyperlink string `json:"hyperlink"` // 链接地址
}

// SearchWord 搜索关键字
type SearchWord struct {
	Engine       string // 搜索引擎
	Device       string // 搜索设备
	KeywordId    int    // 关键字ID
	Keyword      string // 关键字
	Domain       string // 关键字域名
	Page         int    // 页数
	ClientId     string // 客户端ID
	SerialNumber string
}

// NewSearchWordsStore 创建关键词存储对象
func NewSearchWordsStore(wd *SearchWord) *SearchWordsStore {
	swm := &SearchWordsStore{
		wd:    wd,
		input: make(chan *StoreItem, 4096),
		done:  make(chan struct{}, 1),
	}
	return swm
}

// KeywordID ...
func (s *SearchWordsStore) KeywordID() string {
	return fmt.Sprintf("%d", s.wd.KeywordId)
}

// SetKeywordID ...
func (s *SearchWordsStore) SetKeywordID(keywordId int) {
	s.wd.KeywordId = keywordId
}

// Start ...
func (s *SearchWordsStore) Start() *SearchWordsStore {
	go s.run()
	return s
}

// Store ...
func (s *SearchWordsStore) Store(values map[string]interface{}) error {
	title := values["title"].(string)
	hyperlink := values["href"].(string)

	// if title == "" && hyperlink == "" {
	// 	return errors.New(fmt.Sprintf("title and href is empty"))
	// }

	data := &StoreItem{Title: title, Hyperlink: hyperlink}
	select {
	case s.input <- data:
	default:
		return errors.New("缓冲区已满")
	}

	return nil
}

// Print ...
func (s *SearchWordsStore) Foreach(fn func(index int, data *StoreItem) bool) {
	s.store.All(func(index int, value interface{}) bool {
		gogo := fn(index, value.(*StoreItem))
		return gogo
	})
}

func (s *SearchWordsStore) Top(n int) (string, error) {
	type Item struct {
		Title     string `json:"title"`     // 标题
		Hyperlink string `json:"hyperlink"` // 链接地址
	}
	var ars []Item
	for i, v := range s.store.Values() {
		if i >= n {
			break
		}
		sitem := v.(*StoreItem)
		ars = append(ars, Item{Title: sitem.Title, Hyperlink: sitem.Hyperlink})
	}

	if len(ars) > 0 {
		d, err := json.Marshal(ars)
		if err != nil {
			return "", err
		}
		return string(d), nil
	}

	return "", nil
}

func (s *SearchWordsStore) Length() int {
	return s.store.Size()
}

func (s *SearchWordsStore) run() {
	for {
		select {
		case <-s.done:
			return
		case data := <-s.input:
			s.store.Add(data)
		}
	}
}

// Close ...
func (s *SearchWordsStore) Close() {
	s.done <- struct{}{}
}
