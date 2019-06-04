package config

import (
	"context"

	"github.com/drone/drone/core"
)

func Database(configs core.ConfigStore) core.ConfigService {
	return &database{configs}
}

type database struct {
	configs core.ConfigStore
}

func (d *database) Find(ctx context.Context, req *core.ConfigArgs) (*core.Config, error) {
	return d.configs.FindAfterOrExist(ctx, req.Repo.ID, req.Build.After)
}
