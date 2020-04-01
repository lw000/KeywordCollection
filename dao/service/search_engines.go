package service

import (
	"KeywordCollection/dao/table"
	"KeywordCollection/global"
	log "github.com/sirupsen/logrus"
)

// 搜索引擎数据库服务
type SearchEnginesDaoService struct {
}

func (sess *SearchEnginesDaoService) Query(status int) ([]table.TSearchEngines, error) {
	return sess.selectFromDb(status)
}

func (sess *SearchEnginesDaoService) selectFromDb(status int) ([]table.TSearchEngines, error) {
	query := "SELECT title, name, url, type, status, page FROM search_engines WHERE status=?;"
	rows, err := global.DBReptiledata.DB().Query(query, status)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var engines []table.TSearchEngines
	for rows.Next() {
		var engine table.TSearchEngines
		err = rows.Scan(&engine.Title, &engine.Name, &engine.Url, &engine.Type, &engine.Status, &engine.Page)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		engines = append(engines, engine)
	}

	return engines, nil
}
