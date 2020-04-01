package dbsrv

import (
	"KeywordCollection/constant"
	"KeywordCollection/dao"
	"time"

	log "github.com/sirupsen/logrus"
)

/*
	更新关键字检索结果
*/

type UpdateResultTask struct {
	KeywordId    int    // 关键字ID
	Engine       string // 搜索引擎
	Device       string // 设备类型
	Keyword      string // 关键字
	Ranks        int    // 排行
	Content      string // 内容
	SerialNumber string // 序列号
}

func (u *UpdateResultTask) Run() {
	t := time.Now()
	err := dao.StoreKeywordResultToDbWithTransaction(u.KeywordId, u.SerialNumber, u.Engine, u.Device, u.Keyword, u.Ranks, u.Content, t, constant.KeywordStatusOk)
	if err != nil {
		log.Error(err)
		return
	}
	log.WithFields(log.Fields{"耗时": time.Since(t)}).Info("存储结果")
}
