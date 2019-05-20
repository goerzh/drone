// +build !oss

package config

import (
	"github.com/drone/drone/core"
	"github.com/drone/drone/store/shared/db"
)

// helper function converts the User structure to a set
// of named query parameters.
func toParams(config *core.Config) map[string]interface{} {
	return map[string]interface{}{
		"config_id":    config.ID,
		"config_after": config.After,
		"config_kind":  config.Kind,
		"config_data":  config.Data,
	}
}

// helper function scans the sql.Row and copies the column
// values to the destination object.
func scanRow(scanner db.Scanner, dst *core.Config) error {
	return scanner.Scan(
		&dst.ID,
		&dst.After,
		&dst.Kind,
		&dst.Data,
	)
}
