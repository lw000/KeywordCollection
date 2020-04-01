package service

import (
	"KeywordCollection/dao/table"
	"KeywordCollection/global"
	"database/sql"
	log "github.com/sirupsen/logrus"
	"time"
)

// 关键字数据库服务
type KeywordsDaoService struct {
}

func (kw *KeywordsDaoService) Query(page int32, count int, status int) (map[int]table.TKeyWords, error) {
	if page < 0 {
		page = 0
	}

	if count < 0 {
		count = 0
	}

	log.WithFields(log.Fields{"page": page, "count": count}).Info("查询关键字")

	query := `SELECT id, keywords, level_id, status FROM keywords WHERE status=? LIMIT ?,?;`
	rows, err := global.DBReptiledata.DB().Query(query, status, page, count)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer rows.Close()

	var words map[int]table.TKeyWords
	words = make(map[int]table.TKeyWords, count)

	for rows.Next() {
		var word table.TKeyWords
		err = rows.Scan(&word.Id, &word.Keywords, &word.LevelId, &word.Status)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		// words = append(words, word)
		words[word.Id] = word
	}

	return words, nil
}

func (kw *KeywordsDaoService) UpdateSearchTime(id int, searchtime time.Time, tx *sql.Tx) error {
	var (
		err  error
		stmt *sql.Stmt
	)
	query := `UPDATE keywords SET search_time =? WHERE id =? AND status = 1;`
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
	result, err = stmt.Exec(searchtime.Format("2006-01-02 15:04:05"), id)
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
