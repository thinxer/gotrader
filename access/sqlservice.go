package access

import (
	"database/sql"

	"github.com/thinxer/graphpipe"
)

type SQLProviderConfig struct {
	Driver     string
	DataSource string
}

func newSQLProvider(config *SQLProviderConfig) (*sql.DB, error) {
	db, err := sql.Open(config.Driver, config.DataSource)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func init() {
	graphpipe.Regsiter("SQL", newSQLProvider)
}
