package dao

import (
	"KeywordCollection/global"
	"database/sql"
	log "github.com/sirupsen/logrus"
	"time"
)

// 插入检索结果并同时更新关键字状态
func StoreKeywordResultToDbWithTransaction(keywordId int, serialNumber string, engine string, device string,
	keywords string, ranks int, content string, t time.Time, status int) error {
	// 开启事务
	var (
		err error
		tx  *sql.Tx
	)
	tx, err = global.DBReptiledata.DB().Begin()
	if err != nil {
		log.Error(err)
		return err
	}

	defer func() {
		switch {
		// 事务回滚
		case err != nil:
			if err = tx.Rollback(); err != nil {
				log.Error(err)
			}
			// 提交事务
		default:
			err = tx.Commit()
			if err != nil {
				log.Error(err)
			}
		}
	}()

	// 更新关键字-检索时间
	updateKeywordsFn := func(tx *sql.Tx, id int, searchtime time.Time) (int64, error) {
		var stmt *sql.Stmt
		query := `UPDATE keywords SET search_time =? WHERE id =? AND STATUS = 1;`
		stmt, err = tx.Prepare(query)
		if err != nil {
			log.Error(err)
			return -1, err
		}
		defer stmt.Close()

		var result sql.Result
		result, err = stmt.Exec(searchtime.Format("2006-01-02 15:04:05"), id)
		if err != nil {
			log.Error(err)
			return -1, err
		}

		return result.RowsAffected()
	}

	// 更新关键字-检索状态
	updateKeywordStatusFn := func(tx *sql.Tx, keywordid int, engine string, device string, status int) (int64, error) {
		var stmt *sql.Stmt
		query := `UPDATE keywords_status SET status =? WHERE keywords_id =? AND engines =? AND type =?;`
		stmt, err = tx.Prepare(query)
		if err != nil {
			log.Error(err)
			return -1, err
		}
		defer stmt.Close()

		var result sql.Result
		result, err = stmt.Exec(status, keywordid, engine, device)
		if err != nil {
			log.Error(err)
			return -1, err
		}
		return result.RowsAffected()
	}

	// 插入关键字-检索结果
	insertSearchResultFn := func(tx *sql.Tx, serialNumber string, engine string, device string, keywords string, ranks int, content string, t time.Time) (int64, error) {
		var stmt *sql.Stmt
		query := `INSERT INTO search_result (engines, type, keywords, ranks, content, create_time, serial_number) VALUES (?,?,?,?,?,?,?);`
		stmt, err = tx.Prepare(query)
		if err != nil {
			log.Error(err)
			return -1, err
		}
		defer stmt.Close()

		var result sql.Result
		result, err = stmt.Exec(engine, device, keywords, ranks, content, t.Format("2006-01-02 15:04:05"), serialNumber)
		if err != nil {
			log.Error(err)
			return -1, err
		}
		return result.LastInsertId()
	}

	// 1. 更新关键字-检索状态
	var n int64
	n, err = updateKeywordStatusFn(tx, keywordId, engine, device, status)
	if err != nil {
		log.Error(err)
		return err
	}
	if n > 0 {

	}

	// 2. 更新关键字-检索状态
	n, err = updateKeywordsFn(tx, keywordId, t)
	if err != nil {
		log.Error(err)
		return err
	}
	if n > 0 {

	}

	// 3. 关键字检索结果，写入数据库
	n, err = insertSearchResultFn(tx, serialNumber, engine, device, keywords, ranks, content, t)
	if err != nil {
		log.Error(err)
		return err
	}
	if n > 0 {

	}

	return nil
}
