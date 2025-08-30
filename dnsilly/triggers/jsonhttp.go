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

package triggers

import (
	"bytes"
	"dnsilly/config"
	"dnsilly/rules"
	"encoding/json"
	"io"
	"net/http"
)

type TriggerEventPayload struct {
	Tag      string   `json:"tag"`
	Domain   string   `json:"domain"`
	Ipv4     []string `json:"ipv4"`
	Ipv6     []string `json:"ipv6"`
	ClientIP string   `json:"client_ip"`
}

type TriggerLifecyclePayload struct {
	State string `json:"state"`
}

func TriggerEventJSONHTTP(conf *config.Config, jhConf *config.ConfigTriggerJSONHTTP, rule *rules.Rule, domain string, ipv4 []string, ipv6 []string, client_ip string) error {
	if jhConf.EventEndpoint == "" {
		return nil
	}

	payload := TriggerEventPayload{
		Tag:      rule.Tag,
		Domain:   domain,
		Ipv4:     ipv4,
		Ipv6:     ipv6,
		ClientIP: client_ip,
	}
	payloadBytes, _ := json.Marshal(payload)

	resp, err := http.Post(
		jhConf.EventEndpoint,
		"application/json",
		bytes.NewBuffer(payloadBytes),
	)

	if err != nil {
		return err
	}

	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()

	return nil
}

func TriggerLifecycleJSONHTTP(conf *config.Config, jhConf *config.ConfigTriggerJSONHTTP, state string) error {
	if jhConf.LifecycleEndpoint == "" {
		return nil
	}

	payload := TriggerLifecyclePayload{
		State: state,
	}
	payloadBytes, _ := json.Marshal(payload)

	resp, err := http.Post(
		jhConf.LifecycleEndpoint,
		"application/json",
		bytes.NewBuffer(payloadBytes),
	)

	if err != nil {
		return err
	}

	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()

	return nil
}
