// Copyright 2019 Drone IO, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package core

import "context"

type (
	// Config represents a pipeline config file.
	Config struct {
		ID     int64  `json:"id,omitempty"`
		RepoID int64  `json:"repo_id,omitempty"`
		After  string `json:"after,omitempty"`
		Data   string `json:"data"`
		Kind   string `json:"kind"`
	}

	// ConfigArgs represents a request for the pipeline
	// configuration file (e.g. .drone.yml)
	ConfigArgs struct {
		User   *User       `json:"-"`
		Repo   *Repository `json:"repo,omitempty"`
		Build  *Build      `json:"build,omitempty"`
		Config *Config     `json:"config,omitempty"`
	}

	// ConfigStore defines operations for working with configs.
	ConfigStore interface {
		// Find returns a build from the datastore.
		Find(context.Context, int64) (*Config, error)

		// List returns a config list from the datastore.
		List(context.Context, int64) ([]*Config, error)

		// FindAfter returns a config from the configstore.
		FindAfter(context.Context, int64, string) (*Config, error)

		// FindAfterOrExist returns a config from the configstore, or nil, nil if not exist
		FindAfterOrExist(context.Context, int64, string) (*Config, error)

		// UpdateOrCreate updates a build in the datastore, or create a new entry if not exist
		UpdateOrCreate(context.Context, *Config) error

		// Create persists a build to the datastore.
		Create(context.Context, *Config) error

		// Update updates a build in the datastore.
		Update(context.Context, *Config) error

		// Delete deletes a build from the datastore.
		Delete(context.Context, *Config) error
	}

	// ConfigService provides pipeline configuration from an
	// external service.
	ConfigService interface {
		Find(context.Context, *ConfigArgs) (*Config, error)
	}
)
