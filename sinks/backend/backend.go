/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package backend

import (
	"github.com/orb-community/orb/pkg/types"
)

type Backend interface {
	Metadata() interface{}
	CreateFeatureConfig() []ConfigFeature
	ValidateConfiguration(config types.Metadata) error
	ParseConfig(format string, config string) (types.Metadata, error)
	ConfigToFormat(format string, metadata types.Metadata) (string, error)
}

const ConfigFeatureTypePassword = "password"
const ConfigFeatureTypeText = "text"

type ConfigFeature struct {
	Type     string `json:"type"`
	Input    string `json:"input"`
	Title    string `json:"title"`
	Name     string `json:"name"`
	Required bool   `json:"required"`
}

type SinkFeature struct {
	Backend     string          `json:"backend"`
	Description string          `json:"description"`
	Config      []ConfigFeature `json:"config"`
}

var registry = make(map[string]Backend)

func Register(name string, b Backend) {
	registry[name] = b
}

func GetList() []string {
	keys := make([]string, 0, len(registry))
	for k := range registry {
		keys = append(keys, k)
	}
	return keys
}

func HaveBackend(name string) bool {
	_, prs := registry[name]
	return prs
}

func GetBackend(name string) Backend {
	if name == "" {
		return nil
	}
	return registry[name]
}
