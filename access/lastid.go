package access

import (
	"database/sql"

	"github.com/thinxer/graphpipe"
)

type LastIdProviderConfig struct {
	TableName string
}

func newLastIdProvider(config *LastIdProviderConfig, db *sql.DB) (int, error) {
	var id int
	row := db.QueryRow("SELECT MAX(id) FROM " + config.TableName)
	err := row.Scan(&id)
	return id, err
}

func init() {
	graphpipe.Regsiter("LastId", newLastIdProvider)
}
