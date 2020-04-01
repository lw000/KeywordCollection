package parseserv

import (
	"github.com/Workiva/go-datastructures/queue"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

type parseServer struct {
	sync.Mutex
	running bool
	quit    bool
	tasks   *queue.Queue
	exit    chan struct{}
	wg      sync.WaitGroup
	worker  int
	handler []func(*htmlParse) error
}

const (
	defaultWorkers      int   = 4
	defaultMaxTaskCount int64 = 4
)

var (
	parseServerOnce     sync.Once
	parseServerInstance *parseServer
)

func ParseServer() *parseServer {
	parseServerOnce.Do(func() {
		parseServerInstance = &parseServer{
			worker: defaultWorkers,
			tasks:  queue.New(512),
		}
		parseServerInstance.init()
	})

	return parseServerInstance
}

func (p *parseServer) init() {
	p.exit = make(chan struct{}, 1)
	p.AddParseHandler(saveKeywordResultHandler, saveDomainRankResultHandler)
}

func (p *parseServer) Start() *parseServer {
	if !p.running {
		p.running = true

		p.wg.Add(1)
		go p.run()

		for i := 1; i <= p.worker; i++ {
			p.wg.Add(1)
			go p.runTask(i)
		}
	}
	return p
}

func (p *parseServer) Stop() {
	if p == nil {
		return
	}
	close(p.exit)
	p.wg.Wait()
}

func (p *parseServer) AddParseHandler(handler ...func(*htmlParse) error) {
	p.Lock()
	defer p.Unlock()
	p.handler = append(p.handler, handler...)
}

func (p *parseServer) AddTask(tasks ...HTMLData) error {
	return p.tasks.Put(tasks)
}

func (p *parseServer) run() {
	defer func() {
		p.wg.Done()
		log.Error("parseServer, exit")
	}()

	t := time.NewTicker(time.Second * 5)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			if p.tasks.Len() < 5 {
			}
		case <-p.exit:
			p.quit = true
			return
		}
	}
}

func (p *parseServer) Len() int64 {
	return p.tasks.Len()
}

func (p *parseServer) runTask(workerId int) {
	defer func() {
		p.wg.Done()
		log.WithFields(log.Fields{"workerId": workerId}).Error("parseServer worker, exit")
	}()

	log.WithFields(log.Fields{"workerId": workerId}).Info("parseServer worker, running")

	for !p.quit {
		items, err := p.tasks.Get(defaultMaxTaskCount)
		if err != nil {
			log.Error(err)
			continue
		}

		for _, e := range items {
			eCopy := e
			switch eCopy.(type) {
			case HTMLData:
				data := eCopy.(HTMLData)
				go p.parse(data)
			default:
			}
		}
	}
}

// 解析网页数据服务
func (p *parseServer) parse(data HTMLData) {
	defer func() {
		if x := recover(); x != nil {
			log.Error(x)
		}
	}()

	parseHtml := NewHTMLParse(data, true)
	defer parseHtml.Close()

	var err error
	if err = parseHtml.Start(); err != nil {
		log.Error(err)
		return
	}

	if err = parseHtml.Do(); err != nil {
		log.Error(err)
		return
	}

	for _, h := range p.handler {
		if err = h(parseHtml); err != nil {
			log.Error(err)
		}
	}
}
