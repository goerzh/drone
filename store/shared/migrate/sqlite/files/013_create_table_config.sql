-- name: create-table-configs

CREATE TABLE IF NOT EXISTS configs (
 config_id          INTEGER PRIMARY KEY AUTOINCREMENT
,config_repo_id     INTEGER
,config_after       TEXT
,config_kind        TEXT
,config_data        TEXT
,UNIQUE(config_repo_id, config_after)
,FOREIGN KEY(config_repo_id) REFERENCES repos(repo_id) ON DELETE CASCADE
);

-- name: create-index-configs-repo

CREATE INDEX IF NOT EXISTS ix_config_repo ON configs (config_repo_id);

-- name: create-index-configs-repo-name

CREATE INDEX IF NOT EXISTS ix_config_repo_name ON configs (config_repo_id, config_after);