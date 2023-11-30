// Copyright 2023 Laurynas Četyrkinas
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bunny

import (
	"context"
	"net/http"
	"net/netip"
	"sync"
	"time"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"

	"github.com/digilolnet/bunnynetedgeips/pkg/bunnynetedgeips"
)

func init() {
	caddy.RegisterModule(BunnyIPRange{})
}

// BunnyIPRange provides a range of IP address prefixes (CIDRs) retrieved from https://api.bunny.net/system/edgeserverlist and https://api.bunny.net/system/edgeserverlist/ipv6.
type BunnyIPRange struct {
	// refresh Interval
	Interval caddy.Duration `json:"interval,omitempty"`
	// request Timeout
	Timeout caddy.Duration `json:"timeout,omitempty"`

	// Holds the parsed CIDR ranges from Ranges.
	ranges []netip.Prefix

	ctx  caddy.Context
	lock *sync.RWMutex
}

// CaddyModule returns the Caddy module information.
func (BunnyIPRange) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.ip_sources.bunny",
		New: func() caddy.Module { return new(BunnyIPRange) },
	}
}

// getContext returns a cancelable context, with a timeout if configured.
func (s *BunnyIPRange) getContext() (context.Context, context.CancelFunc) {
	if s.Timeout > 0 {
		return context.WithTimeout(s.ctx, time.Duration(s.Timeout))
	}
	return context.WithCancel(s.ctx)
}

func (s *BunnyIPRange) getPrefixes() ([]netip.Prefix, error) {
	ctx, cancel := s.getContext()
	defer cancel()
	ips, err := bunnynetedgeips.BunnynetEdgeIPs(ctx)
	if err != nil {
		return nil, err
	}
	var prefixes []netip.Prefix
	for _, ip := range ips {
		prefix, err := caddyhttp.CIDRExpressionToPrefix(ip)
		if err != nil {
			return nil, err
		}
		prefixes = append(prefixes, prefix)
	}
	return prefixes, nil
}

func (s *BunnyIPRange) refreshLoop() {
	if s.Interval == 0 {
		s.Interval = caddy.Duration(time.Hour)
	}
	ticker := time.NewTicker(time.Duration(s.Interval))
	s.lock.Lock()
	s.ranges, _ = s.getPrefixes()
	s.lock.Unlock()
	for {
		select {
		case <-ticker.C:
			prefixes, err := s.getPrefixes()
			if err != nil {
				break
			}
			s.lock.Lock()
			s.ranges = prefixes
			s.lock.Unlock()
		case <-s.ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func (s *BunnyIPRange) Provision(ctx caddy.Context) error {
	s.ctx = ctx
	s.lock = new(sync.RWMutex)
	go s.refreshLoop()
	return nil
}

func (s *BunnyIPRange) GetIPRanges(_ *http.Request) []netip.Prefix {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.ranges
}

// UnmarshalCaddyfile implements caddyfile.Unmarshaler.
//
//	bunny {
//	   interval val
//	   timeout val
//	}
func (m *BunnyIPRange) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	d.Next() // Skip module name.

	// No same-line options are supported
	if d.NextArg() {
		return d.ArgErr()
	}

	for nesting := d.Nesting(); d.NextBlock(nesting); {
		switch d.Val() {
		case "interval":
			if !d.NextArg() {
				return d.ArgErr()
			}
			val, err := caddy.ParseDuration(d.Val())
			if err != nil {
				return err
			}
			m.Interval = caddy.Duration(val)
		case "timeout":
			if !d.NextArg() {
				return d.ArgErr()
			}
			val, err := caddy.ParseDuration(d.Val())
			if err != nil {
				return err
			}
			m.Timeout = caddy.Duration(val)
		default:
			return d.ArgErr()
		}
	}

	return nil
}

// interface guards
var (
	_ caddy.Module            = (*BunnyIPRange)(nil)
	_ caddy.Provisioner       = (*BunnyIPRange)(nil)
	_ caddyfile.Unmarshaler   = (*BunnyIPRange)(nil)
	_ caddyhttp.IPRangeSource = (*BunnyIPRange)(nil)
)
