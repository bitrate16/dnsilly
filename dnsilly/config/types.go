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

// Server configuration
type ConfigServer struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// DNS Upstream config
type ConfigUpstream struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type ConfigTriggerCommand struct {
	// Asynchronous trigger
	Async bool `yaml:"async"`

	// Run commands in batch instead of per-ip
	Batch bool `yaml:"batch"`

	// Command template
	// Accepts parameters:
	// - {tag} - rule tag
	// - {domain} - domain name
	// - {client_ip} - client ip
	// - {type} - DNS response type (A or AAAA)
	// - {ips} - comma-separated list of ips from response if `batch=true`
	// - {ip} - ip from response if `batch=false`
	EventTemplate string `yaml:"event_template"`

	// Lifecycle template
	// Accepts parameters:
	// - {state} - lifecycle state
	LifecycleTemplate string `yaml:"lifecycle_template"`

	// Execute on server start
	OnStart string `yaml:"on_start"`

	// Execute on server stop
	OnStop string `yaml:"on_stop"`

	// Execute on server start during config reload
	OnPartialStart string `yaml:"on_partial_start"`

	// Execute on server stop during config reload
	OnPartialStop string `yaml:"on_partial_stop"`
}

type ConfigTriggerJSONHTTP struct {
	// Asynchronous trigger
	Async bool `yaml:"async"`

	// JSON HTTP request
	// Payload:
	// {
	//     "tag": "<rule tag>",
	//     "domain": "<domain name>",
	//     "mask": "<dns mask from rule>",
	//     "ipv4": [
	//         "comma-separated list of ipv4 in response",
	//     ],
	//     "ipv6": [
	//         "comma-separated list of ipv6 in response",
	//     ],
	//     "client_ip": "Client IP Address"
	// }
	EventEndpoint string `yaml:"event_endpoint"`

	// JSON HTTP request
	// Payload:
	// {
	//     "state": "<lifecycle state>",
	// }
	LifecycleEndpoint string `yaml:"lifecycle_endpoint"`
}

// Trigger config
type ConfigTrigger struct {
	Command  []*ConfigTriggerCommand  `yaml:"command"`
	JSONHTTP []*ConfigTriggerJSONHTTP `yaml:"json_http"`
}

type Config struct {
	// Verbose log
	Verbose bool `yaml:"verbose"`

	// Reload interval
	Reload time.Duration `yaml:"reload"`

	// Server configuration
	Server *ConfigServer `yaml:"server"`

	// Rules file path
	Rules string `default:"dnsilly.rules" yaml:"rules"`

	// Upstreams for DNS resolution
	Upstreams []*ConfigUpstream `yaml:"upstreams"`
	Trigger   *ConfigTrigger    `yaml:"trigger"`
}
