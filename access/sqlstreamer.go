package access

import (
	"database/sql"
	"fmt"

	s "github.com/thinxer/coincross"
	"github.com/thinxer/graphpipe"
)

type SQLStreamer struct {
	tid     int
	value   *s.Trade
	pending chan *s.Trade

	query  *sql.Rows
	closed bool
	pair   s.Pair
}

type SQLStreamerConfig struct {
	TableName string
	Pair      string
}

func newSQLStreamer(config *SQLStreamerConfig, db *sql.DB, since, limit int) (*SQLStreamer, error) {
	q := fmt.Sprintf(`SELECT id, timestamp, type, price, amount FROM %s WHERE id >= ?`, config.TableName)
	params := []interface{}{since}
	if limit > 0 {
		q = q + " LIMIT ?"
		params = append(params, limit)
	}
	query, err := db.Query(q, params...)
	if err != nil {
		return nil, err
	}
	ms := &SQLStreamer{query: query, pending: make(chan *s.Trade, 4096)}
	(&ms.pair).Set(config.Pair)
	return ms, nil
}

func (m *SQLStreamer) Start(ch chan bool) {
	for m.query.Next() {
		trade := s.Trade{Pair: m.pair}
		m.query.Scan(&trade.Id, &trade.Timestamp, &trade.Type, &trade.Price, &trade.Amount)
		m.pending <- &trade
		ch <- true
	}
	m.query.Close()
	close(m.pending)
	ch <- true
	close(ch)
}

func (m *SQLStreamer) Update(tid int) graphpipe.Result {
	t, ok := <-m.pending
	if ok {
		m.tid, m.value = tid, t
	} else {
		m.closed = true
	}
	return graphpipe.Update
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
