package table

import (
	"github.com/Workiva/go-datastructures/queue"
	"time"
)

// TSearchResult 搜索结果
type TSearchResult struct {
	SerialNumber string    `db:"serial_number"`
	Domain       string    `db:"domain_id"`
	Engines      string    `db:"engines_id"`
	Keywords     string    `db:"keywords"`
	Ranks        int       `db:"ranks"`
	TopThree     string    `db:"top_three"`
	CreateTime   time.Time `db:"create_time"`
}

// QueryKeyWordContext ...
type QueryContext struct {
	Priority     int
	SerialNumber string
	KeywordId    int
	Engine       string
	Type         string
	Keyword      string
	Page         int
	ClientId     string
	QueryId      string
}

func (q *QueryContext) Compare(other queue.Item) int {
	d := other.(*QueryContext)
	if q.Priority == d.Priority {
		return 0
	}

	if q.Priority > d.Priority {
		return 1
	}

	return -1
}
