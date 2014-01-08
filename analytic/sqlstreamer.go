package analytic

import (
	"database/sql"
	"fmt"

	s "github.com/thinxer/gocoins"
	"github.com/thinxer/graphpipe"
)

var findStmt = `SELECT id, timestamp, type, price, amount FROM %s WHERE id >= ?;`

type SQLStreamer struct {
	tid   int
	value *s.Trade

	query  *sql.Rows
	closed bool
	pair   s.Pair
}

type SQLStreamerConfig struct {
	Driver     string
	DataSource string
	TableName  string
	Pair       string
	Since      int64
}

// TODO when to Close the DB?
func newSQLStreamer(config *SQLStreamerConfig) (*SQLStreamer, error) {
	db, err := sql.Open(config.Driver, config.DataSource)
	if err != nil {
		return nil, err
	}
	query, err := db.Query(fmt.Sprintf(findStmt, config.TableName), config.Since)
	if err != nil {
		return nil, err
	}
	ms := &SQLStreamer{query: query}
	(&ms.pair).Set(config.Pair)
	return ms, nil
}

func (m *SQLStreamer) Update(tid int) bool {
	if m.query.Next() {
		m.tid = tid
		trade := s.Trade{Pair: m.pair}
		m.query.Scan(&trade.Id, &trade.Timestamp, &trade.Type, &trade.Price, &trade.Amount)
		m.value = &trade
		return true
	} else {
		m.closed = true
		return false
	}
}

func (m *SQLStreamer) Value() (int, *s.Trade) {
	return m.tid, m.value
}

func (m *SQLStreamer) Closed() bool {
	return m.closed
}

func init() {
	graphpipe.Regsiter("SQLStreamer", newSQLStreamer)
}
