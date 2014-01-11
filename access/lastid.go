package access

import (
	"github.com/thinxer/graphpipe"
)

type LastIdProviderConfig struct {
	TableName string
}

func newLastIdProvider(config *LastIdProviderConfig, sql SQLService) (int, error) {
	var id int
	row := sql.DB().QueryRow("SELECT MAX(id) FROM " + config.TableName)
	err := row.Scan(&id)
	return id, err
}

func init() {
	graphpipe.Regsiter("LastId", newLastIdProvider)
}
