package table

type TDomain struct {
	Id     int    `db:"id"`
	Name   string `db:"name"`
	Domain string `db:"domain"`
	Status int    `db:"status"`
}

func (t *TDomain) Decode(data map[string]interface{}) error {
	t.Id = int(data["id"].(int64))
	t.Name = string(data["name"].([]uint8))
	t.Domain = string(data["domain"].([]uint8))
	t.Status = int(data["status"].(int64))
	return nil
}
