// Copyright 2015 Jason Mansfield
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

package portrange

import (
	"errors"
)

const (
	// 0 is technically a valid protocol number, but this struct is only for
	// TCP and UDP.
	ProtoInvalid = 0
	ProtoTcp     = 6
	ProtoUdp     = 17
)

var (
	ErrBadRange       = errors.New("Bad port range")
	ErrBadProto       = errors.New("Bad protocol")
	ErrDisjointRanges = errors.New("Disjoint port ranges")
)

// A TCP or UDP port ranges.
type PortRange struct {
	minPort, maxPort uint16
	proto            uint8
}

// New creates a new PortRange. Returns ErrBadRange if either port 0 or
// if the max port is higher than the min port. Returns ErrBadProto if proto
// isn't one of ProtoTcp or ProtoUdp.
func New(minPort, maxPort uint16, proto uint8) (*PortRange, error) {
	p := &PortRange{minPort: minPort, maxPort: maxPort, proto: proto}
	err := p.Validate()
	if err != nil {
		return nil, err
	}
	return p, nil
}

// Validate returns an error indicating a semantic problem with the port range.
// nil indicates no problems.
func (p *PortRange) Validate() error {
	if p.minPort == 0 || p.maxPort == 0 {
		return ErrBadRange
	}
	if p.minPort > p.maxPort {
		return ErrBadRange
	}
	if p.proto != ProtoTcp && p.proto != ProtoUdp {
		return ErrBadProto
	}
	return nil
}

// Overlaps indicates if two ranges are the same protocol and share at least
// one port. This property is symmetric.
func (p *PortRange) Overlaps(o *PortRange) bool {
	if p.proto != o.proto {
		return false
	}
	if p.minPort > o.maxPort {
		return false
	}
	if o.minPort > p.maxPort {
		return false
	}
	return true
}

// Adjacent indicates if two ranges are the same protocol and the ranges are
// symmetric but do not overlap. This property is symmetric.
func (p *PortRange) Adjacent(o *PortRange) bool {
	if p.proto != o.proto {
		return false
	}
	if p.Overlaps(o) {
		return false
	}
	if p.maxPort+1 == o.minPort {
		return true
	}
	if o.maxPort+1 == p.minPort {
		return true
	}
	return false
}

// Overlap calculcates the overlap (intersection) between the two ranges and
// updates the argument to contain only that overlap. Returns ErrDisjointRanges
// if the ranges do not overlap.
func (p *PortRange) Overlap(o *PortRange) error {
	if !p.Overlaps(o) {
		return ErrDisjointRanges
	}
	minPort := p.minPort
	if o.minPort > minPort {
		minPort = o.minPort
	}
	maxPort := p.maxPort
	if o.maxPort < maxPort {
		maxPort = o.maxPort
	}
	o.minPort = minPort
	o.maxPort = maxPort
	return nil
}

// MergeWith updates the destination range to be the superset of both ranges if
// they overlap or are adjacent. Returns ErrDisjointRanges if they can't be
// merged.
func (p *PortRange) MergeWith(o *PortRange) error {
	if p.Overlaps(o) || p.Adjacent(o) {
		if o.minPort > p.minPort {
			o.minPort = p.minPort
		}
		if o.maxPort < p.maxPort {
			o.maxPort = p.maxPort
		}
		return nil
	}
	return ErrDisjointRanges
}

// EntirelyLessThan indicates that the entirety of one range is below the other.
// All TCP ports are lower than any UDP port. Returns false when ranges overlap.
func (p *PortRange) EntirelyLessThan(o *PortRange) bool {
	if p.proto < o.proto {
		return true
	}
	if p.Overlaps(o) {
		return false
	}
	if p.maxPort < o.minPort {
		return true
	}
	return false
}
