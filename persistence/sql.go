package persistence

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
	db       *sql.DB
	inserter *sql.Stmt
	source   a.TradeSource
}

type SQLWriterConfig struct {
	Driver     string
	DataSource string
	TableName  string
}

func NewSQLWriter(config *SQLWriterConfig, source a.TradeSource) (*SQLWriter, error) {
	db, err := sql.Open(config.Driver, config.DataSource)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(fmt.Sprintf(tableCreateStmt, config.TableName))
	if err != nil {
		db.Close()
		return nil, err
	}

	inserter, err := db.Prepare(fmt.Sprintf(insertStmt, config.TableName))
	if err != nil {
		db.Close()
		return nil, err
	}

	return &SQLWriter{db: db, inserter: inserter, source: source}, nil
}

func (w *SQLWriter) Update(tid int) bool {
	_, t := w.source.Value()
	_, err := w.inserter.Exec(t.Id, t.Timestamp, int(t.Type), t.Price, t.Amount)
	if err != nil {
		panic(err)
	}
	return false
}

func (w *SQLWriter) Closed() bool {
	if w.source.Closed() {
		// Close methods are idempotent.
		w.db.Close()
		w.inserter.Close()
	}
	return w.source.Closed()
}

func init() {
	graphpipe.Regsiter("SQLWriter", NewSQLWriter)
}
