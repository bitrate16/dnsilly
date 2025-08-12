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

import "time"

func defaultConfig() *Config {
	return &Config{
		Verbose: true,
		Reload:  10 * time.Second,
		Rules:   "dnsilly.rules",
		Server: &ConfigServer{
			Host: "0.0.0.0",
			Port: 53,
		},
		Upstreams: []*ConfigUpstream{
			&ConfigUpstream{
				Host: "8.8.8.8",
				Port: 53,
			},
		},
		Trigger: &ConfigTrigger{
			Command: []*ConfigTriggerCommand{
				&ConfigTriggerCommand{
					Batch:             false,
					EventTemplate:     "echo 'tag={tag} domain={domain} type={type} ip={ip}'",
					LifecycleTemplate: "echo '{state}'",
					OnStart:           "echo on_start",
					OnStop:            "echo on_stop",
					OnPartialStart:    "echo on_partial_start",
					OnPartialStop:     "echo on_partial_stop",
				},
			},
			JSONHTTP: []*ConfigTriggerJSONHTTP{
				&ConfigTriggerJSONHTTP{
					EventEndpoint:     "https://api.example.com/v1/firewall/event",
					LifecycleEndpoint: "https://api.example.com/v1/firewall/lifecycle",
				},
			},
		},
	}
}
