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
	"dnsilly/config"
	"dnsilly/rules"
	"dnsilly/util"
	"fmt"
)

func TriggerEvent(conf *config.Config, rule *rules.Rule, domain string, ipv4 []string, ipv6 []string) {
	if conf.Verbose {
		fmt.Printf("[%s] Trigger event: domain=%s, rule=%s\n", util.Now(), domain, rule.Tag)
	}

	if conf.Trigger == nil {
		return
	}

	for _, cmdConf := range conf.Trigger.Command {
		if cmdConf.Async {
			go func() {
				err := TriggerEventCommand(conf, cmdConf, rule, domain, ipv4, ipv6)
				if err != nil {
					fmt.Printf("[%s] Trigger event command error: %v\n", util.Now(), err)
				}
			}()
		} else {
			err := TriggerEventCommand(conf, cmdConf, rule, domain, ipv4, ipv6)
			if err != nil {
				fmt.Printf("[%s] Trigger event command error: %v\n", util.Now(), err)
			}
		}
	}

	for _, jhConf := range conf.Trigger.JSONHTTP {
		if jhConf.Async {
			go func() {
				err := TriggerEventJSONHTTP(conf, jhConf, rule, domain, ipv4, ipv6)
				if err != nil {
					fmt.Printf("[%s] Trigger event json http error: %v\n", util.Now(), err)
				}
			}()
		} else {
			err := TriggerEventJSONHTTP(conf, jhConf, rule, domain, ipv4, ipv6)
			if err != nil {
				fmt.Printf("[%s] Trigger event json http error: %v\n", util.Now(), err)
			}
		}
	}
}

func TriggerLifecycle(conf *config.Config, state string) {
	if conf.Verbose {
		fmt.Printf("[%s] Trigger lifecycle state: %s\n", util.Now(), state)
	}

	if conf.Trigger == nil {
		return
	}

	for _, cmdConf := range conf.Trigger.Command {
		if cmdConf.Async {
			go func() {
				err := TriggerLifecycleCommand(conf, cmdConf, state)
				if err != nil {
					fmt.Printf("[%s] Trigger lifecycle command error: %v\n", util.Now(), err)
				}
			}()
		} else {
			err := TriggerLifecycleCommand(conf, cmdConf, state)
			if err != nil {
				fmt.Printf("[%s] Trigger lifecycle command error: %v\n", util.Now(), err)
			}
		}
	}

	for _, jhConf := range conf.Trigger.JSONHTTP {
		if jhConf.Async {
			go func() {
				err := TriggerLifecycleJSONHTTP(conf, jhConf, state)
				if err != nil {
					fmt.Printf("[%s] Trigger lifecycle json http error: %v\n", util.Now(), err)
				}
			}()
		} else {
			err := TriggerLifecycleJSONHTTP(conf, jhConf, state)
			if err != nil {
				fmt.Printf("[%s] Trigger lifecycle json http error: %v\n", util.Now(), err)
			}
		}
	}
}
