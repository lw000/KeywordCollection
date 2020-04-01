package models

import (
	"errors"
	"fmt"
	tyutils "github.com/lw000/gocommon/utils"

	"github.com/Workiva/go-datastructures/queue"
	log "github.com/sirupsen/logrus"
)

// SearchRequestItem ...
type SearchRequestItem struct {
	RequestID string // 请求ID
	Engine    string // 搜索引擎
	Wd        string // 关键字
	Page      int
}

// SearchRequestManager ...
type SearchRequestManager struct {
	isStart bool
	input   chan *SearchRequestItem
	done    chan struct{}
	queue   *queue.Queue
}

// NewSearchRequestManager 创建搜索请求队列管理
func NewSearchRequestManager() *SearchRequestManager {
	return &SearchRequestManager{
		input: make(chan *SearchRequestItem, 4096),
		done:  make(chan struct{}, 1),
		queue: queue.New(4096),
	}
}

// Add ...
func (s *SearchRequestManager) Add(req *SearchRequestItem) (queryID string, err error) {
	select {
	case s.input <- req:
	default:
		return "", errors.New("缓冲区已满")
	}
	return tyutils.UUID(), nil
}

// Start ...
func (s *SearchRequestManager) Start() *SearchRequestManager {
	if !s.isStart {
		s.isStart = true
		go s.run()
	}
	return s
}

func (s *SearchRequestManager) run() {
	for {
		select {
		case <-s.done:
			return
		case req := <-s.input:
			if err := s.queue.Put(req); err != nil {
				log.Error(err)
			}
		}
	}
}

// Close ...
func (s *SearchRequestManager) Close() {
	s.done <- struct{}{}
}

// Get ...
func (s *SearchRequestManager) Get() (*SearchRequestItem, error) {
	v, err := s.queue.Peek()
	if err != nil {
		return nil, err
	}
	return v.(*SearchRequestItem), nil

	// return nil, nil
}

func (s SearchRequestItem) String() string {
	return fmt.Sprintf("Engine:%s Wd:%s Page:%d", s.Engine, s.Wd, s.Page)
}
