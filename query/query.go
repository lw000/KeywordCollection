package query

import (
	"KeywordCollection/constant"
	"KeywordCollection/dao/service"
	"KeywordCollection/dao/table"
	"KeywordCollection/errors"
	"KeywordCollection/global"
	"KeywordCollection/searchengine"
	"KeywordCollection/searchengine/360"
	"KeywordCollection/searchengine/baidu"
	"KeywordCollection/searchengine/shenma"
	"KeywordCollection/searchengine/sogou"
	log "github.com/sirupsen/logrus"
	"sync"
)

// KeywordQuery 检索关键字
func KeywordQuery(w *sync.WaitGroup, qctx *table.QueryContext) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}

		if w != nil {
			w.Done()
		}
	}()

	engines, err := initEngines()
	if err != nil {
		log.Error(err)
		return
	}

	wg := &sync.WaitGroup{}
	for _, enin := range engines {
		if enin.Opt().Engine == qctx.Engine && enin.Opt().Type == qctx.Type {
			enin1 := enin
			wg.Add(1)
			go query(wg, enin1, qctx)
		}
	}
	wg.Wait()
}

func query(wg *sync.WaitGroup, engi searchengine.WebSearchEngine, qctx *table.QueryContext) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}

		engi.Close()

		if wg != nil {
			wg.Done()
		}
	}()

	word := searchengine.SearchWord{
		Engine:       qctx.Engine,
		Device:       qctx.Type,
		KeywordId:    qctx.KeywordId,
		Keyword:      qctx.Keyword,
		Page:         qctx.Page,
		ClientId:     qctx.ClientId,
		SerialNumber: qctx.SerialNumber,
	}
	var err error
	if err = engi.Start(&word); err != nil {
		log.Error(err)
		return
	}

	if err = engi.SearchChrome(); err != nil {
		log.Error(err)
		return
	}
}

func initEngines() (map[searchengine.SearchKey]searchengine.WebSearchEngine, error) {
	serv := service.SearchEnginesDaoService{}
	enginesConfig, err := serv.Query(constant.EnginesStatusEnable)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if len(enginesConfig) == 0 {
		err = errors.New(0, "未配置搜索引擎", "")
		log.Error(err)
		return nil, err
	}

	engines := make(map[searchengine.SearchKey]searchengine.WebSearchEngine)

	for _, cfg := range enginesConfig {
		if cfg.Name == "" || cfg.Type == "" || cfg.Url == "" {
			log.WithFields(log.Fields{"engine": cfg}).Error("搜索引擎配置错误")
			continue
		}

		key := searchengine.SearchKey{Engine: cfg.Name, Device: cfg.Type}
		opt := &searchengine.SearchOption{Domain: cfg.Url, Engine: cfg.Name, Type: cfg.Type}
		switch cfg.Name {
		case constant.EngineBaidu:
			switch cfg.Type {
			case constant.DevicePc:
				engines[key] = baidu.NewPcSearch(global.Chromes[constant.DevicePc], opt)
			case constant.DeviceMobile:
				engines[key] = baidu.NewMobileSearch(global.Chromes[constant.DeviceMobile], opt)
			}
		case constant.EngineSogou:
			switch cfg.Type {
			case constant.DevicePc:
				engines[key] = sogou.NewPcSearch(global.Chromes[constant.DevicePc], opt)
			case constant.DeviceMobile:
				engines[key] = sogou.NewMobileSearch(global.Chromes[constant.DeviceMobile], opt)
			}
		case constant.EngineSo360:
			switch cfg.Type {
			case constant.DevicePc:
				engines[key] = so360.NewPcSearch(global.Chromes[constant.DevicePc], opt)
			case constant.DeviceMobile:
				engines[key] = so360.NewMobileSearch(global.Chromes[constant.DeviceMobile], opt)
			}
		case constant.EngineShenma:
			switch cfg.Type {
			case constant.DevicePc:
			case constant.DeviceMobile:
				engines[key] = shenma.NewMobileSearch(global.Chromes[constant.DeviceMobile], opt)
			}
		default:
		}
	}

	return engines, nil
}
