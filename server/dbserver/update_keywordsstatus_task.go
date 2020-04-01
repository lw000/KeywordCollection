package dbsrv

import (
	"KeywordCollection/global"
	"time"

	log "github.com/sirupsen/logrus"
)

type UpdateKeywordsStatusTask struct {
	KeywordId int
	Engine    string
	Device    string
	Status    int
}

func (u *UpdateKeywordsStatusTask) Run() {
	t := time.Now()

	result, err := global.DBReptiledata.DB().Exec("UPDATE keywords_status SET status=? WHERE keywords_id=? AND engines=? AND type=?;", u.Status, u.KeywordId, u.Engine, u.Device)
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

	log.WithFields(log.Fields{"status": u.Status, "ts": time.Since(t)}).Info("更新检索状态")
}
