package table

import "time"

type TDomainRanks struct {
	DomainId     int       `db:"domain_id"`
	KeywordsId   int       `db:"keywords_id"`
	Engines      string    `db:"engines"`
	Type         string    `db:"type"`
	Ranks        int       `db:"ranks"`
	Content      string    `db:"content"`
	SerialNumber string    `db:"serial_number"`
	CreateTime   time.Time `db:"create_time"`
}

func (t *TDomainRanks) Decode(data map[string]interface{}) error {
	t.DomainId = int(data["domain_id"].(int64))
	t.KeywordsId = int(data["keywords_id"].(int64))
	t.Engines = string(data["engines"].([]uint8))
	t.Type = string(data["type"].([]uint8))
	t.Ranks = int(data["ranks"].(int64))
	t.Content = string(data["content"].([]uint8))
	t.SerialNumber = string(data["serial_number"].([]uint8))
	t.CreateTime = data["create_time"].(time.Time)
	return nil
}
