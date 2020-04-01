package qserv

import (
	"KeywordCollection/constant"
	"KeywordCollection/dao/service"
	"KeywordCollection/dao/table"
	"KeywordCollection/server/dbserver"
	"KeywordCollection/server/taskQueue"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
	"sync"
	"sync/atomic"
	"time"
)

// 数据库关键字定时查询
// 查询的数据同步到TaskServer服务
type queryServer struct {
	running bool
	quit    bool
	offset  int32
	cr      *cron.Cron
	exit    chan struct{}
	wg      sync.WaitGroup
}

const (
	defaultQueryCount int = 100
)

var (
	queryServerInstance     *queryServer
	queryServerInstanceOnce sync.Once
)

func QueryServer() *queryServer {
	queryServerInstanceOnce.Do(func() {
		queryServerInstance = &queryServer{
			cr:     cron.New(),
			offset: 1,
		}
		queryServerInstance.init()
	})

	return queryServerInstance
}

func (qs *queryServer) init() {
	qs.exit = make(chan struct{}, 1)
}

func (qs *queryServer) Start() *queryServer {
	if !qs.running {
		qs.running = true
		// 每天凌晨1点执行一次：0 0 1 * * ?
		err := qs.cr.AddFunc("0 0 1 * * ?", qs.prepareTaskData)
		// err := qs.cr.AddFunc("0 */2 * * * ?", qs.prepareTaskData)
		// err := qs.cr.AddFunc("*/5 * * * * ?", qs.prepareTaskData)
		if err != nil {
			log.Panic(err)
		}

		qs.cr.Start()

		qs.wg.Add(1)
		go qs.run()

		qs.wg.Add(1)
		go qs.queryData()
	}
	return qs
}

func (qs *queryServer) Stop() {
	if qs == nil {
		return
	}
	close(qs.exit)
	qs.wg.Wait()
}

func (qs *queryServer) run() {
	defer func() {
		qs.wg.Done()
		log.Error("queryServer, exit")
	}()

	t := time.NewTicker(time.Second * 20)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			size := qtasks.ContentRetrievalServer().Len()
			if size < 5 {
				qs.wg.Add(1)
				go qs.queryData()
			}
		case <-qs.exit:
			qs.quit = true
			return
		}
	}
}

func (qs *queryServer) prepareTaskData() {
	log.Info("每日凌晨1点开始执行任务")

	// 重置关键字检索状态
	serv := service.KeywordsStatusDaoService{}
	err := serv.BatchUpdateStatus(constant.KeywordStatusDefault)
	if err != nil {
		log.Error(err)
		return
	}

	qs.wg.Add(1)
	go qs.queryData()
}

// 查询待检索关键字
func (qs *queryServer) queryData() {
	defer qs.wg.Done()

	// 查询关键字信息
	servKeywords := service.KeywordsDaoService{}
	keyswords, err := servKeywords.Query((qs.offset-1)*int32(defaultQueryCount), defaultQueryCount, constant.KeywordStatusDefault)
	if err != nil {
		log.Error(err)
		return
	}

	if len(keyswords) == 0 {
		atomic.StoreInt32(&qs.offset, 1)
		return
	}

	if len(keyswords) > 0 {
		atomic.AddInt32(&qs.offset, 1)
	}

	// 查询搜索引擎
	serv := service.SearchEnginesDaoService{}
	engines, err := serv.Query(constant.EnginesStatusEnable)
	if err != nil {
		log.Error(err)
		return
	}

	qcontexts := make([]*table.QueryContext, 0, defaultQueryCount*len(engines))

	var kwdsId []int
	for _, kwd := range keyswords {
		kwdsId = append(kwdsId, kwd.Id)
	}

	// 获取待检索关键字状态
	var keyswordsStatus []table.TKeyWordsStatus
	servKeywordsStatus := service.KeywordsStatusDaoService{}
	// 获取-等待检索状态的关键字,从新开始检索
	if keyswordsStatus, err = servKeywordsStatus.Query(kwdsId, constant.KeywordStatusDefault); err != nil {
		log.Error(err)
	}

	if len(keyswordsStatus) == 0 {
		// 获取-上次未检索完成的关键字,从新开始检索
		if keyswordsStatus, err = servKeywordsStatus.Query(kwdsId, constant.KeywordStatusChecking); err != nil {
			log.Error(err)
		}
	}

	if len(keyswordsStatus) == 0 {
		// 获取-检索失败的关键字,从新开始检索
		if keyswordsStatus, err = servKeywordsStatus.Query(kwdsId, constant.KeywordStatusFail); err != nil {
			log.Error(err)
		}
	}

	for _, engi := range engines {
		if engi.Page <= 0 {
			continue
		}

		for _, kwds := range keyswordsStatus {
			if (kwds.Engines == engi.Name) && (kwds.Type == engi.Type) {
				kwd, ok := keyswords[kwds.KeywordsId]
				if !ok {
					continue
				}
				ctx := &table.QueryContext{
					KeywordId: kwd.Id,
					Keyword:   kwd.Keywords,
					Page:      engi.Page,
					Engine:    kwds.Engines,
					Type:      kwds.Type,
				}
				qcontexts = append(qcontexts, ctx)
			}
		}
	}

	// for _, kwd := range keyswords {
	// 	// 获取待检索关键字状态
	// 	var keyswordsStatus []table.TKeyWordsStatus
	// 	servKeywordsStatus := service.KeywordsStatusDaoService{}
	// 	if keyswordsStatus, err = servKeywordsStatus.Query(kwd.Id, constant.KeywordStatusDefault); err != nil {
	// 		log.Error(err)
	// 		continue
	// 	}
	//
	// 	if len(keyswordsStatus) == 0 {
	// 		// 获取检索上次未检索完成的关键字从新开始检索
	// 		if keyswordsStatus, err = servKeywordsStatus.Query(kwd.Id, constant.KeywordStatusChecking); err != nil {
	// 			log.Error(err)
	// 			continue
	// 		}
	// 	}
	//
	// 	if len(keyswordsStatus) == 0 {
	// 		// 获取检索失败的关键字从新开始检索
	// 		if keyswordsStatus, err = servKeywordsStatus.Query(kwd.Id, constant.KeywordStatusFail); err != nil {
	// 			log.Error(err)
	// 			continue
	// 		}
	// 	}
	//
	// 	for _, kwds := range keyswordsStatus {
	// 		for _, engi := range engines {
	// 			if (kwds.Engines == engi.Name) && (kwds.Type == engi.Type) && (engi.Page > 0) {
	// 				ctx := &table.QueryContext{
	// 					KeywordId: kwd.Id,
	// 					Keyword:   kwd.Keywords,
	// 					Page:      engi.Page,
	// 					Engine:    kwds.Engines,
	// 					Type:      kwds.Type,
	// 				}
	// 				qcontexts = append(qcontexts, ctx)
	// 			}
	// 		}
	// 	}
	// }

	if len(qcontexts) == 0 {
		log.WithFields(log.Fields{"keyword": qcontexts}).Info("关键字")
	}

	for _, ctx := range qcontexts {
		// 更新关键字检索状态
		err = dbsrv.DBServer().AddJob(&dbsrv.UpdateKeywordsStatusTask{
			KeywordId: ctx.KeywordId,
			Engine:    ctx.Engine,
			Device:    ctx.Type,
			Status:    constant.KeywordStatusChecking,
		})
		if err != nil {
			log.Error(err)
			continue
		}

		// 添加关键字到检索队列服务中
		if err = qtasks.ContentRetrievalServer().Put(ctx); err != nil {
			log.Error(err)
			continue
		}
	}
}
