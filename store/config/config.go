// +build !oss

package config

import (
	"context"

	"github.com/drone/drone/core"
	"github.com/drone/drone/store/shared/db"
)

func New(db *db.DB) core.ConfigStore {
	return &configStore{db}
}

type configStore struct {
	db *db.DB
}

// Find returns a build from the datastore.
func (s *configStore) Find(ctx context.Context, id int64) (*core.Config, error) {
	out := &core.Config{ID: id}
	err := s.db.View(func(queryer db.Queryer, binder db.Binder) error {
		params := toParams(out)
		query, args, err := binder.BindNamed(queryKey, params)
		if err != nil {
			return err
		}
		row := queryer.QueryRow(query, args...)
		return scanRow(row, out)
	})
	if err != nil && err.Error() == "sql: no rows in result set" {
		return nil, nil
	}
	return out, err
}

// FindAfter returns a config from the configstore.
func (s *configStore) FindAfter(ctx context.Context, after string) (*core.Config, error) {
	out := &core.Config{After: after}
	err := s.db.View(func(queryer db.Queryer, binder db.Binder) error {
		params := toParams(out)
		query, args, err := binder.BindNamed(queryAfter, params)
		if err != nil {
			return err
		}
		row := queryer.QueryRow(query, args...)
		return scanRow(row, out)
	})
	if err != nil && err.Error() == "sql: no rows in result set" {
		return nil, nil
	}
	return out, err
}

// Create persists a build to the datastore.
func (s *configStore) Create(ctx context.Context, config *core.Config) error {
	return s.db.Lock(func(execer db.Execer, binder db.Binder) error {
		params := toParams(config)
		stmt, args, err := binder.BindNamed(stmtInsert, params)
		if err != nil {
			return err
		}
		res, err := execer.Exec(stmt, args...)
		if err != nil {
			return err
		}
		config.ID, err = res.LastInsertId()
		return err
	})
}

// Update updates a build in the datastore.
func (s *configStore) Update(ctx context.Context, config *core.Config) error {
	return s.db.Lock(func(execer db.Execer, binder db.Binder) error {
		params := toParams(config)
		stmt, args, err := binder.BindNamed(stmtUpdate, params)
		if err != nil {
			return err
		}
		_, err = execer.Exec(stmt, args...)
		return err
	})
}

// Delete deletes a build from the datastore.
func (s *configStore) Delete(ctx context.Context, config *core.Config) error {
	return s.db.Lock(func(execer db.Execer, binder db.Binder) error {
		params := toParams(config)
		stmt, args, err := binder.BindNamed(stmtDelete, params)
		if err != nil {
			return err
		}
		_, err = execer.Exec(stmt, args...)
		return err
	})
}

const queryBase = `
SELECT
 config_id
,config_after
,config_kind
,config_data
`

const queryKey = queryBase + `
FROM configs
WHERE config_id = :config_id
LIMIT 1
`

const queryAfter = queryBase + `
FROM configs
WHERE config_after = :config_after
LIMIT 1
`

const stmtInsert = `
INSERT INTO configs (
 config_after
,config_kind
,config_data
) VALUES (
 :config_after
,:config_kind
,:config_data
)
`

const stmtUpdate = `
UPDATE configs SET
 config_id = :config_id
,config_after = :config_after
,config_kind = :config_kind
,config_data = :config_data
WHERE config_id = :config_id
`

const stmtDelete = `
DELETE FROM configs
WHERE config_id = :config_id
`
