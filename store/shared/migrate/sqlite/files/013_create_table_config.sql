-- name: create-table-configs

CREATE TABLE IF NOT EXISTS configs (
 config_id          INTEGER PRIMARY KEY AUTOINCREMENT
,config_after       TEXT
,config_kind        TEXT
,config_data        TEXT
,UNIQUE(config_after)
);