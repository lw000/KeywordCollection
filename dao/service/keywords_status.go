package service

import (
	"KeywordCollection/dao/table"
	"KeywordCollection/global"
	"database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
)

// 关键字状态数据库服务
type KeywordsStatusDaoService struct {
}

func (kws *KeywordsStatusDaoService) Query(keywordsId []int, status int) ([]table.TKeyWordsStatus, error) {
	return kws.selectFromDb(keywordsId, status)
}

func (kws *KeywordsStatusDaoService) QueryWithKeywordsId(keywordsId int, status int) ([]table.TKeyWordsStatus, error) {
	return kws.selectWithKeywordsIdFromDb(keywordsId, status)
}

func (kws *KeywordsStatusDaoService) selectFromDb(keywordsId []int, status int) ([]table.TKeyWordsStatus, error) {
	query := `SELECT keywords_id, engines, type, status FROM keywords_status WHERE keywords_id in(?) AND status=?;`
	var s string

	for _, v := range keywordsId {
		s += fmt.Sprintf("%d,", v)
	}
	s = strings.TrimRight(s, ",")
	rows, err := global.DBReptiledata.DB().Query(query, s, status)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer rows.Close()

	var words []table.TKeyWordsStatus
	for rows.Next() {
		var word table.TKeyWordsStatus
		err = rows.Scan(&word.KeywordsId, &word.Engines, &word.Type, &word.Status)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		words = append(words, word)
	}

	return words, nil
}

func (kws *KeywordsStatusDaoService) selectWithKeywordsIdFromDb(keywordsId int, status int) ([]table.TKeyWordsStatus, error) {
	query := `SELECT keywords_id, engines, type, status FROM keywords_status WHERE keywords_id in(?) AND status=?;`
	// rows, err := global.DBReptiledata.DB().Query(query, fmt.Sprintf("%d, %d", keywordsId, keywordsId), status)
	rows, err := global.DBReptiledata.DB().Query(query, keywordsId, status)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer rows.Close()

	var words []table.TKeyWordsStatus
	for rows.Next() {
		var word table.TKeyWordsStatus
		err = rows.Scan(&word.KeywordsId, &word.Engines, &word.Type, &word.Status)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		words = append(words, word)
	}

	return words, nil
}

// 批量更新关键字检索状态
func (kws *KeywordsStatusDaoService) BatchUpdateStatus(status int) error {
	query := "UPDATE keywords_status SET status=? WHERE status in (2,3,4);"
	result, err := global.DBReptiledata.DB().Exec(query, status)
	if err != nil {
		log.Error(err)
		return err
	}

	var n int64
	n, err = result.RowsAffected()
	if err != nil {
		log.Error(err)
		return err
	}
	if n > 0 {

	}

	return nil
}

// 更新关键字检索状态
func (kws *KeywordsStatusDaoService) UpdateStatus(keywordId int, engine string, device string, status int, tx *sql.Tx) error {
	var (
		err  error
		stmt *sql.Stmt
	)
	query := `UPDATE keywords_status SET status =? WHERE keywords_id =? AND engines =? AND type =?;`
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
	result, err = stmt.Exec(status, keywordId, engine, device)
	if err != nil {
		log.Error(err)
		return err
	}

	var n int64
	n, err = result.RowsAffected()
	if err != nil {
		log.Error(err)
		return err
	}
	if n > 0 {

	}
	return nil
}
