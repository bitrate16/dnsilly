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
	"dnsilly/util"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var shell string
var hasShell bool

func init() {
	// Use env shell
	shell = os.Getenv("SHELL")
	if shell != "" {
		hasShell = true
		return
	}

	// Use hardcoded shell
	commonShells := []string{
		"/bin/sh",
		"/bin/bash",
		"/bin/zsh",
		"/bin/fish",
	}
	for _, shell = range commonShells {
		if _, err := exec.LookPath(shell); err == nil {
			hasShell = true
			return
		}
	}

	// No shell
	hasShell = false
}

func executeForError(command string, verbose bool) error {
	if verbose {
		fmt.Printf("[%s] exec: `%s`\n", util.Now(), command)
	}

	proc := exec.Command(shell, "-c", command)
	output, err := proc.CombinedOutput()
	if err != nil {
		return fmt.Errorf("exec failed (%v): %v", err, string(output))
	}

	if verbose {
		fmt.Printf("[%s] stdout: `%s`\n", util.Now(), bytes.TrimSpace(output))
	}

	return nil
}

func partialTriggerEventCommand(conf *config.Config, cmdConf *config.ConfigTriggerCommand, command string, ips []string, proto string) error {
	if len(ips) != 0 {
		subCommand := strings.ReplaceAll(command, "{type}", proto)

		if cmdConf.Batch {
			subSubCommand := strings.ReplaceAll(subCommand, "{ips}", strings.Join(ips, ","))

			// Execute
			err := executeForError(subSubCommand, conf.Verbose)
			if err != nil {
				return err
			}
		} else {
			for _, ip := range ips {
				subSubCommand := strings.ReplaceAll(subCommand, "{ip}", ip)

				// Execute
				err := executeForError(subSubCommand, conf.Verbose)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func TriggerEventCommand(conf *config.Config, cmdConf *config.ConfigTriggerCommand, rule *rules.Rule, domain string, ipv4 []string, ipv6 []string) error {
	if !hasShell {
		return errors.New("shell not found")
	}

	if cmdConf.EventTemplate == "" {
		return nil
	}

	command := cmdConf.EventTemplate
	command = strings.ReplaceAll(command, "{tag}", rule.Tag)
	command = strings.ReplaceAll(command, "{domain}", domain)

	err := partialTriggerEventCommand(conf, cmdConf, command, ipv4, "A")
	if err != nil {
		return err
	}

	return partialTriggerEventCommand(conf, cmdConf, command, ipv4, "AAAA")
}

func TriggerLifecycleCommand(conf *config.Config, cmdConf *config.ConfigTriggerCommand, state string) error {
	if !hasShell {
		return errors.New("shell not found")
	}

	// Execute distinct triggers
	if (state == OnStart) && (cmdConf.OnStart != "") {
		err := executeForError(cmdConf.OnStart, conf.Verbose)
		if err != nil {
			return err
		}
	}

	if (state == OnStop) && (cmdConf.OnStop != "") {
		err := executeForError(cmdConf.OnStop, conf.Verbose)
		if err != nil {
			return err
		}
	}

	if (state == OnPartialStart) && (cmdConf.OnPartialStart != "") {
		err := executeForError(cmdConf.OnPartialStart, conf.Verbose)
		if err != nil {
			return err
		}
	}

	if (state == OnPartialStop) && (cmdConf.OnPartialStop != "") {
		err := executeForError(cmdConf.OnPartialStop, conf.Verbose)
		if err != nil {
			return err
		}
	}

	// Execute handler script
	if cmdConf.LifecycleTemplate == "" {
		return nil
	}

	command := cmdConf.LifecycleTemplate
	command = strings.ReplaceAll(command, "{state}", state)

	return executeForError(command, conf.Verbose)
}
