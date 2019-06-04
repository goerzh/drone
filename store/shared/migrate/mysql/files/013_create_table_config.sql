-- name: create-table-configs

CREATE TABLE IF NOT EXISTS configs (
 config_id          INTEGER PRIMARY KEY AUTO_INCREMENT
,config_repo_id     INTEGER
,config_after       VARCHAR(50)
,config_kind        VARCHAR(10)
,config_data        VARCHAR(2000)
,UNIQUE(config_repo_id, config_after)
,FOREIGN KEY(config_repo_id) REFERENCES repos(repo_id) ON DELETE CASCADE
);

-- name: create-index-configs-repo

CREATE INDEX ix_config_repo ON configs (config_repo_id);

-- name: create-index-configs-repo-after

CREATE INDEX ix_config_repo_after ON configs (config_repo_id, config_after);