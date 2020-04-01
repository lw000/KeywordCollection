package table

import "time"

// TKeyWords 搜索关键字
type TKeyWords struct {
	Id         int       `db:"id"`          // 关键字ID
	Keywords   string    `db:"keywords"`    // 关键字
	LevelId    int       `db:"level_id"`    // 关键字等级
	Status     int       `db:"status"`      // 关键字状态
	CreateTime time.Time `db:"create_time"` // 创建时间
	SearchTime time.Time `db:"starch_time"` // 检索时间
}

func (t *TKeyWords) Decode(data map[string]interface{}) error {
	t.Id = int(data["id"].(int64))
	t.Keywords = string(data["keywords"].([]uint8))
	t.LevelId = int(data["level_id"].(int64))
	t.Status = int(data["status"].(int64))
	return nil
}
