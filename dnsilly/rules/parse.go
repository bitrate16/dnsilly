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

package rules

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Supported rules:
// Wildcard:
// - "*" - match anything
// - "*.foo" - match any ".foo" subdomain
// - "?" - match single character
//
// TODO:
//
// Optional dot:
// - "+.foo" - match "*.foo" and "foo".
//   - "mar.foo"
//   - "foo" - only if not starting with
//
// - "bar.+.foo" - match anything ending (prefix) or surrounded or with dot or nothing "bar.*.foo" and "bar.foo"
//   - "bar.foo"
//   - "bar.tar.foo"
//
// - "^" - optional dot (example: "*^bar.com" - match "foo.bar" and "foobar" and "bar")
//
// Specific letter classes
// - "?%" and "*%" - match numeric character
// - "?$" and "*$" - match letter character
// - "?@" and "*@" - match non-numeric and non-letter character (example: ".-_")
func makeRegexp(matcher string) (*regexp.Regexp, error) {
	matcher = strings.ReplaceAll(matcher, ".", "\\.")
	matcher = strings.ReplaceAll(matcher, "*", ".*")
	matcher = strings.ReplaceAll(matcher, "?", ".")

	return regexp.Compile("^" + matcher + "$")
}

func ParseRules(configPath string) (*Rules, error) {
	// Create if not exists
	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(configPath)
		if err != nil {
			return nil, err
		}

		f.WriteString("# Example:\n")
		f.WriteString("# block example.com\n")
		f.WriteString("# allow analytics.example.com\n")
		f.WriteString("# block *.example.com\n")
		f.Close()
	}

	// Prepare config
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)
	lineno := 0

	rules := &Rules{
		Rules: make([]*Rule, 0),
	}

	for scanner.Scan() {
		line := scanner.Bytes()
		lineno += 1

		// Remove comments '#'
		index := bytes.Index(line, []byte("#"))
		if index != -1 {
			line = line[:index]
		}

		// Clean string
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		// Parse rule
		fields := bytes.Fields(line)

		if len(fields) != 2 {
			return nil, fmt.Errorf("invalid rule at line %d", lineno)
		}

		tag := string(fields[0])
		regex, err := makeRegexp(string(fields[1]))
		if err != nil {
			return nil, err
		}

		rules.Rules = append(rules.Rules, &Rule{
			Tag:     tag,
			Regexp:  regex,
			Pattern: string(fields[1]),
		})
	}

	return rules, nil
}
