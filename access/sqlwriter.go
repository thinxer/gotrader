package access

import (
	"fmt"
	a "github.com/thinxer/gotrader/analytic"
	"github.com/thinxer/graphpipe"

	"database/sql"
)

var (
	tableCreateStmt = `
CREATE TABLE IF NOT EXISTS %s (
	id BIGINT NOT NULL PRIMARY KEY,
	timestamp BIGINT NOT NULL,
	type TINYINT NOT NULL,
	price REAL NOT NULL,
	amount REAL NOT NULL
);`
	insertStmt = `REPLACE INTO %s (id, timestamp, type, price, amount) VALUES (?, ?, ?, ?, ?)`
)

type SQLWriter struct {
	inserter *sql.Stmt
	source   a.TradeSource
}

type SQLWriterConfig struct {
	TableName string
}

func newSQLWriter(config *SQLWriterConfig, db *sql.DB) (*SQLWriter, error) {
	var err error
	_, err = db.Exec(fmt.Sprintf(tableCreateStmt, config.TableName))
	if err != nil {
		return nil, err
	}

	inserter, err := db.Prepare(fmt.Sprintf(insertStmt, config.TableName))
	if err != nil {
		inserter.Close()
		return nil, err
	}

	return &SQLWriter{inserter: inserter}, nil
}

func (s *SQLWriter) SetInput(source a.TradeSource) {
	s.source = source
}

func (w *SQLWriter) Update(tid int) graphpipe.Result {
	_, t := w.source.Value()
	_, err := w.inserter.Exec(t.Id, t.Timestamp, int(t.Type), t.Price, t.Amount)
	if err != nil {
		panic(err)
	}
	return graphpipe.Skip
}

func (w *SQLWriter) Closed() bool {
	if w.source.Closed() {
		// Close methods are idempotent.
		w.inserter.Close()
	}
	return w.source.Closed()
}

func init() {
	graphpipe.Regsiter("SQLWriter", newSQLWriter)
}
