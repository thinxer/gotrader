package access

import (
	"database/sql"

	"github.com/thinxer/graphpipe"
)

type SQLService interface {
	DB() *sql.DB
}

type SQLProvider struct {
	db *sql.DB
}

type SQLProviderConfig struct {
	Driver     string
	DataSource string
}

func newSQLProvider(config *SQLProviderConfig) (*SQLProvider, error) {
	db, err := sql.Open(config.Driver, config.DataSource)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	return &SQLProvider{db}, nil
}

func (s *SQLProvider) DB() *sql.DB {
	return s.db
}

func init() {
	graphpipe.Regsiter("SQL", newSQLProvider)
}
