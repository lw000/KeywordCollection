package qtasks

import (
	"KeywordCollection/dao/table"
	"KeywordCollection/global"
	"KeywordCollection/query"
	"github.com/Workiva/go-datastructures/queue"
	log "github.com/sirupsen/logrus"
	"sync"
)

type contentRetrievalServer struct {
	running bool
	quit    bool
	worker  int
	exit    chan struct{}
	wg      sync.WaitGroup
	tasks   *queue.PriorityQueue
}

const (
	defaultWorkers      int = 1
	defaultMaxTaskCount int = 2
)

var (
	contentRetrievalServerOnce     sync.Once
	contentRetrievalServerInstance *contentRetrievalServer
)

func ContentRetrievalServer() *contentRetrievalServer {
	contentRetrievalServerOnce.Do(func() {
		contentRetrievalServerInstance = &contentRetrievalServer{
			worker: defaultWorkers,
			tasks:  queue.NewPriorityQueue(defaultWorkers*defaultMaxTaskCount, false),
		}
		contentRetrievalServerInstance.init()
	})

	return contentRetrievalServerInstance
}

func (ts *contentRetrievalServer) init() {
	ts.exit = make(chan struct{}, 1)
}

func (ts *contentRetrievalServer) Start() *contentRetrievalServer {
	if !ts.running {
		ts.running = true

		{
			ts.wg.Add(1)
			go ts.run()
		}

		for i := 1; i <= ts.worker; i++ {
			ts.wg.Add(1)
			go ts.runTask(i)
		}
	}
	return ts
}

func (ts *contentRetrievalServer) Stop() {
	if ts == nil {
		return
	}
	close(ts.exit)
	ts.wg.Wait()
}

func (ts *contentRetrievalServer) run() {
	defer func() {
		ts.wg.Done()
		log.Error("taskServer，exit")
	}()

	for {
		select {
		case <-ts.exit:
			ts.quit = true
			return
		}
	}
}

func (ts *contentRetrievalServer) Put(v ...queue.Item) error {
	return ts.tasks.Put(v...)
}

func (ts *contentRetrievalServer) Len() int {
	return ts.tasks.Len()
}

func (ts *contentRetrievalServer) runTask(workerId int) {
	defer func() {
		ts.wg.Done()
		log.WithFields(log.Fields{"workerId": workerId}).Error("qtasks worker, exit")
	}()
	log.WithFields(log.Fields{"workerId": workerId}).Info("qtasks worker, running")

	for !ts.quit {
		items, err := ts.tasks.Get(defaultMaxTaskCount)
		if err != nil {
			log.Error(err)
			continue
		}

		var wg sync.WaitGroup
		for _, e := range items {
			switch e.(type) {
			case *table.QueryContext:
				queryCtx := e.(*table.QueryContext)
				queryCtx.SerialNumber = global.GetIdWorker().String()
				wg.Add(1)
				go query.KeywordQuery(&wg, queryCtx)
			default:
			}
		}

		wg.Wait()
	}
}
