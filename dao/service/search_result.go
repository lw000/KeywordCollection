package service

import (
	"KeywordCollection/global"
	"database/sql"
	log "github.com/sirupsen/logrus"
	"time"
)

// 搜索结果数据库服务
type SearchResultDaoService struct {
}

// 插入检索结果
func (sr *SearchResultDaoService) Insert(serialNumber string, engine string, device string, keywords string, ranks int, content string, t time.Time, tx *sql.Tx) error {
	var (
		err  error
		stmt *sql.Stmt
	)
	query := `INSERT INTO search_result (engines, type, keywords, ranks, content, create_time, serial_number) VALUES (?,?,?,?,?,?,?);`
	if tx == nil {
		stmt, err = global.DBReptiledata.DB().Prepare(query)
	} else {
		stmt, err = tx.Prepare(query)
	}
	if err != nil {
		log.Error(err)
		return err
	}
	defer stmt.Close()

	var result sql.Result
	result, err = stmt.Exec(engine, device, keywords, ranks, content, t.Format("2006-01-02 15:04:05"), serialNumber)
	if err != nil {
		log.Error(err)
		return err
	}
	var n int64
	n, err = result.LastInsertId()
	if err != nil {
		log.Error(err)
		return err
	}
	if n > 0 {

	}
	return nil
}
