// dnsilly - dns automation utility
// Copyright (C) 2025  bitrate16 (bitrate16@gmail.com)
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package config

import (
	"dnsilly/util"
	"errors"
	"flag"
	"os"

	"gopkg.in/yaml.v2"
)

func GetConfigPath() string {
	// Prepare config path
	var configPath string

	flag.StringVar(
		&configPath,
		"config",
		"dnsilly.yml",
		"path to config file",
	)

	flag.Parse()

	return configPath
}

func ParseConfig(configPath string) (*Config, error) {
	// Create if not exists
	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(configPath)
		if err != nil {
			return nil, err
		}

		err = yaml.NewEncoder(f).Encode(defaultConfig())
		f.Close()

		if err != nil {
			return nil, err
		}
	}

	// Prepare config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}

	// Make bank default
	util.SetDefaults(cfg)

	return cfg, nil
}
