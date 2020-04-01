package table

// TKeyWordsStatus 关键字状态
type TKeyWordsStatus struct {
	KeywordsId int    `db:"keywords_id"` // 关键字ID
	Engines    string `db:"engines"`     // 搜索引擎
	Type       string `db:"type"`        // 搜索引擎类型
	Status     int    `db:"status"`      // 检索状态
	QueyId     int    `db:"quey_id"`     // 检索ID
}

func (t *TKeyWordsStatus) Decode(data map[string]interface{}) error {
	t.KeywordsId = int(data["keywords_id"].(int64))
	t.Engines = string(data["engines"].([]uint8))
	t.Type = string(data["type"].([]uint8))
	t.Status = int(data["status"].(int64))
	return nil
}
