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

package server

import (
	"dnsilly/triggers"
	"dnsilly/util"
	"fmt"
	"net"
	"strconv"

	"github.com/miekg/dns"
)

func logHandler(w dns.ResponseWriter, r *dns.Msg) bool {
	for _, question := range r.Question {
		fmt.Printf("[%s] %s %s %s\n", util.Now(), dns.OpcodeToString[int(question.Qtype)], dns.ClassToString[question.Qclass], question.Name)

	}
	return true
}

type answerGroup struct {
	ipv4s []string
	ipv6s []string
}

// Group answer results per-domain
func makeAnswerGroups(answers []dns.RR) map[string]*answerGroup {
	answerGroups := make(map[string]*answerGroup)

	for _, answer := range answers {
		var ag *answerGroup
		var ok bool

		domain := answer.Header().Name
		if len(domain) == 0 {
			continue
		}

		if domain[len(domain)-1] == '.' {
			domain = domain[:len(domain)-1]
		}

		if ag, ok = answerGroups[domain]; !ok {
			ag = &answerGroup{
				ipv4s: make([]string, 0),
				ipv6s: make([]string, 0),
			}
			answerGroups[domain] = ag
		}

		if headerA, ok := answer.(*dns.A); ok {
			ag.ipv4s = append(ag.ipv4s, headerA.A.String())
		}

		if headerAAAA, ok := answer.(*dns.AAAA); ok {
			ag.ipv6s = append(ag.ipv6s, headerAAAA.AAAA.String())
		}
	}

	return answerGroups
}

func (s *Server) proxyHandler(w dns.ResponseWriter, request *dns.Msg) bool {
	// Seek for first available upstream
	for _, upstream := range s.config.Upstreams {
		addr := net.JoinHostPort(upstream.Host, strconv.Itoa(upstream.Port))

		// TODO: Optimize JoinHostPort
		upstreamResponse, _, err := s.client.Exchange(request, addr)
		if err != nil {
			fmt.Printf("[%s] Upstream %s not available: %v\n", util.Now(), addr, err)
			continue
		}

		answerGroups := makeAnswerGroups(upstreamResponse.Answer)

		// Trigger triggers for matching domains
		for domain, ag := range answerGroups {

			// Check rule match
			rule := s.rules.Match([]byte(domain))
			if rule != nil {
				triggers.TriggerEvent(s.config, rule, domain, ag.ipv4s, ag.ipv6s)
			}
		}

		// Pass response to client
		w.WriteMsg(upstreamResponse)

		return false
	}

	// No upstreams
	fmt.Printf("[%s] No upstream available\n", util.Now())
	response := &dns.Msg{}
	response.SetRcode(request, dns.RcodeServerFailure)
	w.WriteMsg(response)

	return false
}
