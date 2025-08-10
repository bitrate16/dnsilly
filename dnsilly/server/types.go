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
	"dnsilly/config"
	"dnsilly/rules"
	"dnsilly/util"
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/miekg/dns"
)

type Server struct {
	config  *config.Config
	rules   *rules.Rules
	server  *dns.Server
	running bool
	lock    sync.Mutex
	client  *dns.Client

	// Channel to be used for server exit
	onExited chan struct{}
}

func NewServer(
	config *config.Config,
) *Server {
	return &Server{
		config: config,
	}
}

func (s *Server) Start() error {
	s.lock.Lock()

	if s.running {
		s.lock.Unlock()
		return errors.New("server is running")
	}
	s.running = true

	// Register handlers
	chain := NewHandlerChain()

	if s.config.Verbose {
		chain.Add(logHandler)
	}

	chain.Add(s.proxyHandler)

	dns.Handle(".", chain)

	// Create server
	addr := s.config.Server.Host + ":" + strconv.Itoa(s.config.Server.Port)
	s.client = new(dns.Client)
	s.server = &dns.Server{
		Addr: addr,
		Net:  "udp",
	}

	// Make exit callback channel
	s.onExited = make(chan struct{}, 1)
	s.lock.Unlock()

	// Start server
	fmt.Printf("[%s] Listening on %s\n", util.Now(), addr)
	err := s.server.ListenAndServe()
	s.onExited <- struct{}{}

	return err
}

func (s *Server) Stop() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if !s.running {
		return errors.New("server is not running")
	}
	s.running = false

	s.server.Shutdown()
	<-s.onExited

	fmt.Printf("[%s] %s\n", util.Now(), "Server stopped")

	return nil
}

func (s *Server) SetRules(r *rules.Rules) {
	s.rules = r
}
