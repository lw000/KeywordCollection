package dbsrv

import (
	"KeywordCollection/global"
	log "github.com/sirupsen/logrus"
	"time"
)

/*
	更新关键字域名检索结果
*/

type Result struct {
	Rank    int
	Domain  string
	Content string
}

type UpdateDomainRankTask struct {
	DomainId     int      // 域名ID
	KeywordsId   int      // 关键字ID
	Engine       string   // 搜索引擎
	Device       string   // 搜索引擎类型
	Results      []Result // 结果
	SerialNumber string   // 序列号
	Rank         int      // 排行
	Content      string   // 内容
}

func (u *UpdateDomainRankTask) AddResult(domain string, rank int) {
	u.Results = append(u.Results, Result{Domain: domain, Rank: rank})
}

func (u *UpdateDomainRankTask) Run() {
	t := time.Now()
	query := `INSERT INTO domain_ranks (
					domain_id,
					keywords_id,
					engines,
					type,
					ranks,
					content,
					serial_number,
					create_time
				)
				VALUES
					(?,?,?,?,?,?,?,?);`
	result, err := global.DBReptiledata.Exec(query, u.DomainId, u.KeywordsId, u.Engine, u.Device, u.Rank, u.Content, u.SerialNumber, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Error(err)
		return
	}

	n, err := result.LastInsertId()
	if err != nil {
		log.Error(err)
		return
	}

	if n > 0 {

	}

	log.WithFields(log.Fields{"耗时": time.Since(t)}).Info("存储结果")
}
