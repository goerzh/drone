// +build !oss

package config

import (
	"database/sql"
	"github.com/drone/drone/core"
	"github.com/drone/drone/store/shared/db"
)

// helper function converts the User structure to a set
// of named query parameters.
func toParams(config *core.Config) map[string]interface{} {
	return map[string]interface{}{
		"config_id":      config.ID,
		"config_repo_id": config.RepoID,
		"config_after":   config.After,
		"config_kind":    config.Kind,
		"config_data":    config.Data,
	}
}

// helper function scans the sql.Row and copies the column
// values to the destination object.
func scanRow(scanner db.Scanner, dst *core.Config) error {
	return scanner.Scan(
		&dst.ID,
		&dst.RepoID,
		&dst.After,
		&dst.Kind,
		&dst.Data,
	)
}

// helper function scans the sql.Row and copies the column
// values to the destination object.
func scanRows(rows *sql.Rows) ([]*core.Config, error) {
	defer rows.Close()

	configs := make([]*core.Config, 0)
	for rows.Next() {
		config := new(core.Config)
		err := scanRow(rows, config)
		if err != nil {
			return nil, err
		}
		configs = append(configs, config)
	}
	return configs, nil
}
