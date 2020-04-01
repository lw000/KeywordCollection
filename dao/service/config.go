package service

import (
	"KeywordCollection/dao/table"
	"KeywordCollection/global"
	log "github.com/sirupsen/logrus"
)

// 数据库系统配置服务
type ConfigDaoService struct {
}

func (conf *ConfigDaoService) Query() (map[string]interface{}, error) {
	return conf.selectFromDb()
}

func (conf *ConfigDaoService) selectFromDb() (map[string]interface{}, error) {
	query := `SELECT id, name, title, value FROM config;`
	rows, err := global.DBReptiledata.DB().Query(query)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer rows.Close()

	cfgs := make(map[string]interface{})
	for rows.Next() {
		var c table.TConfig
		err = rows.Scan(&c.Id, &c.Name, &c.Title, &c.Value)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		cfgs[c.Name] = c.Value
	}
	return cfgs, nil
}
