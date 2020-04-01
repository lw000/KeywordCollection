package dbsrv

import (
	"KeywordCollection/errors"
	"sync"

	"github.com/Workiva/go-datastructures/queue"
	log "github.com/sirupsen/logrus"
)

type DBFuncJob func()

func (f DBFuncJob) Run() {
	f()
}

type DBJob interface {
	Run()
}

type dbServer struct {
	quit    bool
	running bool
	worker  int
	input   chan interface{}
	tasks   queue.Queue
	exit    chan struct{}
	wg      sync.WaitGroup
}

const (
	defaultWorker             = 4
	defaultMaxTaskCount int64 = 2
)

var (
	instance     *dbServer
	instanceOnce sync.Once
)

// DBServer ...
func DBServer() *dbServer {
	instanceOnce.Do(func() {
		instance = &dbServer{
			worker: defaultWorker,
			input:  make(chan interface{}, 4096),
		}
		instance.init()
	})
	return instance
}

func (d *dbServer) init() {
	d.exit = make(chan struct{}, 1)
}

func (d *dbServer) Start() *dbServer {
	if !d.running {
		d.running = true

		d.wg.Add(1)
		go d.runControl()

		for i := 1; i <= d.worker; i++ {
			d.wg.Add(1)
			go d.runTask(i)
		}
	}
	return d
}

func (d *dbServer) AddJob(jobs ...DBJob) error {
	for _, job := range jobs {
		select {
		case d.input <- job:
		default:
			return errors.New(0, "缓冲区已满", "")
		}
	}

	return nil
}

// Stop ...
func (d *dbServer) Stop() {
	if d == nil {
		return
	}
	close(d.exit)
	d.wg.Wait()
}

func (d *dbServer) runControl() {
	defer func() {
		d.wg.Done()
		log.Error("dbServer, exit")
	}()

	for {
		select {
		case data := <-d.input:
			if err := d.tasks.Put(data); err != nil {
				log.Error(err)
				break
			}
		case <-d.exit:
			d.quit = true
			return
		}
	}
}

func (d *dbServer) runTask(workerId int) {
	defer func() {
		d.wg.Done()
		log.WithFields(log.Fields{"workerId": workerId}).Error("dbServer worker, exit")
	}()

	log.WithFields(log.Fields{"workerId": workerId}).Info("dbServer worker, running")

	for !d.quit {
		items, err := d.tasks.Get(defaultMaxTaskCount)
		if err != nil {
			log.Error(err)
			continue
		}

		for _, e := range items {
			eCopy := e
			go d.run(eCopy)
		}
	}
}

func (d *dbServer) run(e interface{}) {
	defer func() {
		if x := recover(); x != nil {
			log.Error(x)
		}
	}()

	switch e.(type) {
	case DBJob:
		job := e.(DBJob)
		job.Run()
	case DBFuncJob:
		fn := e.(DBFuncJob)
		fn()
	default:

	}
}
