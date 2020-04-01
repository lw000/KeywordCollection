package table

// TSearchEngines 搜索引擎
type TSearchEngines struct {
	Title  string `db:"title"`
	Name   string `db:"name"`
	Url    string `db:"url"`
	Type   string `db:"type"`
	Status int    `db:"status"`
	Page   int    `db:"page"`
}

func (t *TSearchEngines) Decode(data map[string]interface{}) error {
	t.Title = string(data["title"].([]uint8))
	t.Name = string(data["name"].([]uint8))
	t.Url = string(data["url"].([]uint8))
	t.Type = string(data["type"].([]uint8))
	t.Status = int(data["status"].(int64))
	t.Page = int(data["page"].(int64))
	return nil
}
