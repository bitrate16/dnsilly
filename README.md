# DNSilly

Transparent DNS proxy for automation needs.

# What?

Have you ever wanted ~~to become a better version of yourself?~~ to trigger certain actions based on DNS queries?

If the answer is yes, then read below. If the answer is no, still read.

DNSilly allows you to trigger custom rules based on DNS queries from clients.

Example use-cases include:
- Route based on DNS responses
- Build firewall based on DNS queries (block all except allowed)
- Monitor your DNS queries
- Anything else

# How to use

Everything starts with single yaml file:

```yml
# Enable verbose logging
verbose: true

# Reload interval, 0 - to disable
reload: 10s

# Server configuration
server:
  # Server host
  host: 0.0.0.0

  # Server port
  port: 53

# Path to file with rules
rules: dnsilly.rules

# List of upstreams to forward queries
upstreams:
- host: 8.8.8.8
  port: 53

# Trigger rules, optional
trigger:

  # Trigger to execute shell script
  command:
    -
      # Async mode, don't wait for execution
      async: false

      # Batch mode: concatenate ips in comma-separated string instead of calling script for each ip separately
      batch: true

      # Domain hit template:
      # {tag} - your rule tag
      # {domain} - matched domain name
      # {type} - type of query: A or AAAA
      # {ips} - comma-separated list in batch mode
      # {ip} - single ip in non-batch mode
      event_template: echo 'tag={tag} domain={domain} type={type} ips={ips} ip={ip}'

      # Lifecycle trigger template:
      # {state} - one of [start, stop, partial_start, partial_stop]
      lifecycle_template: echo '{state}'

      # Separate triggers
      on_start: echo on_start
      on_stop: echo on_stop
      on_partial_start: echo on_partial_start
      on_partial_stop: echo on_partial_stop

  # HTTP JSON request trigger
  json_http:
    -
      # Async mode, don't wait for execution
      async: true

      # Event trigger endpoint
      # POST Payload:
      # {
      #     "tag": "<rule tag>",
      #     "domain": "<domain name>",
      #     "ipv4": [
      #         "comma-separated list of ipv4 in response",
      #     ],
      #     "ipv6": [
      #         "comma-separated list of ipv6 in response",
      #     ]
      # }
      event_endpoint: https://api.example.com/v1/firewall/event

      # Lifecycle trigger endpoint
      # POST Payload:
      # {
      #     "state": "<lifecycle state>",
      # }
      lifecycle_endpoint: https://api.example.com/v1/firewall/lifecycle
```

And rules file:
```
# Example:
block example.com
allow analytics.example.com
block *.example.com
```

Rules are evaluated sequently and first matching rule activates trigger.

# Rules

Rules support patterns:
- `*` - match any character of any count
- `?` - match any single character

# Examples

## Route based on tag

```yml
verbose: true
reload: 10s
server:
  host: 0.0.0.0
  port: 53

rules: dnsilly.rules

upstreams:
- host: 8.8.8.8
  port: 53

trigger:
  json_http:
    -
      # Async mode, don't wait for execution
      async: true
      event_endpoint: https://api.example.com/v1/firewall/event
```

## Log DNS queries to remote server

```yml
verbose: true
reload: 10s
server:
  host: 0.0.0.0
  port: 53

rules: dnsilly.rules

upstreams:
- host: 8.8.8.8
  port: 53

trigger:
  command:
    -
      async: true
      batch: false
      event_template: echo 'ip route add {ip} dev {tag}'
```


# LICENSE

```
dnsilly - dns automation utility
Copyright (C) 2025  bitrate16 (bitrate16@gmail.com)

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
```
