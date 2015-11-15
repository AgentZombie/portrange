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
	"testing"
)

func TestNew(t *testing.T) {
	for _, tc := range []struct {
		id       string
		min, max uint16
		proto    uint8
		pRange   *PortRange
		err      error
	}{
		{
			id:  "Valid TCP, one port",
			min: 16, max: 16,
			proto:  ProtoTcp,
			pRange: &PortRange{minPort: 16, maxPort: 16, proto: ProtoTcp},
			err:    nil,
		},
		{
			id:  "Valid UDP, one port",
			min: 16, max: 16,
			proto:  ProtoUdp,
			pRange: &PortRange{minPort: 16, maxPort: 16, proto: ProtoUdp},
			err:    nil,
		},
		{
			id:  "Valid TCP, five ports",
			min: 16, max: 20,
			proto:  ProtoTcp,
			pRange: &PortRange{minPort: 16, maxPort: 20, proto: ProtoTcp},
			err:    nil,
		},
		{
			id:  "Valid UDP, five ports",
			min: 16, max: 20,
			proto:  ProtoUdp,
			pRange: &PortRange{minPort: 16, maxPort: 20, proto: ProtoUdp},
			err:    nil,
		},
		{
			id:  "Invalid TCP, port 0",
			min: 0, max: 16,
			proto:  ProtoTcp,
			pRange: nil,
			err:    ErrBadRange,
		},
		{
			id:  "Invalid UDP, port 0",
			min: 0, max: 16,
			proto:  ProtoUdp,
			pRange: nil,
			err:    ErrBadRange,
		},
		{
			id:  "Invalid proto, one port",
			min: 16, max: 16,
			proto:  ProtoInvalid,
			pRange: nil,
			err:    ErrBadProto,
		},
	} {
		t.Logf("Running New test case %q", tc.id)
		p, err := New(tc.min, tc.max, tc.proto)
		if err != tc.err {
			t.Errorf("Case %q, want error %q, got %q", tc.id, tc.err, err)
		}
		if p != nil && tc.pRange == nil {
			t.Errorf("Case %q, want nil range, got %q", tc.id, p)
		}
		if p == nil || tc.pRange == nil {
			continue
		}
		if *p != *tc.pRange {
			t.Errorf("Case %q, want range %q, got %q", tc.id, tc.pRange, p)
		}
	}
}

func TestOverlaps(t *testing.T) {
	for _, tc := range []struct {
		id       string
		p, o     *PortRange
		expected bool
	}{
		{
			id:       "Proto mismatch",
			p:        &PortRange{16, 16, ProtoTcp},
			o:        &PortRange{16, 16, ProtoUdp},
			expected: false,
		},
		{
			id:       "Adjacent non-overlap",
			p:        &PortRange{16, 20, ProtoTcp},
			o:        &PortRange{21, 26, ProtoTcp},
			expected: false,
		},
		{
			id:       "Identity overlap",
			p:        &PortRange{16, 16, ProtoTcp},
			o:        &PortRange{16, 16, ProtoTcp},
			expected: true,
		},
		{
			id:       "Same min overlap",
			p:        &PortRange{16, 25, ProtoTcp},
			o:        &PortRange{16, 21, ProtoTcp},
			expected: true,
		},
		{
			id:       "Same max overlap",
			p:        &PortRange{16, 25, ProtoTcp},
			o:        &PortRange{18, 25, ProtoTcp},
			expected: true,
		},
		{
			id:       "Nesting overlap",
			p:        &PortRange{16, 27, ProtoTcp},
			o:        &PortRange{18, 25, ProtoTcp},
			expected: true,
		},
		{
			id:       "Distinct overlap",
			p:        &PortRange{16, 25, ProtoTcp},
			o:        &PortRange{18, 27, ProtoTcp},
			expected: true,
		},
	} {
		t.Logf("Running Overlaps test case %q", tc.id)
		outcome := tc.p.Overlaps(tc.o)
		if outcome != tc.expected {
			t.Errorf("Case p(o) %q want %v, got %v", tc.id, tc.expected, outcome)
		}
		// Overlap is a symmetric property.
		outcome = tc.o.Overlaps(tc.p)
		if outcome != tc.expected {
			t.Errorf("Case o(p) %q want %v, got %v", tc.id, tc.expected, outcome)
		}
	}
}

func TestAdjacent(t *testing.T) {
	for _, tc := range []struct {
		id       string
		p, o     *PortRange
		expected bool
	}{
		{
			id:       "Proto mismatch",
			p:        &PortRange{16, 16, ProtoTcp},
			o:        &PortRange{16, 16, ProtoUdp},
			expected: false,
		},
		{
			id:       "Adjacent non-overlap",
			p:        &PortRange{16, 20, ProtoTcp},
			o:        &PortRange{21, 26, ProtoTcp},
			expected: true,
		},
		{
			id:       "Identity overlap",
			p:        &PortRange{16, 16, ProtoTcp},
			o:        &PortRange{16, 16, ProtoTcp},
			expected: false,
		},
		{
			id:       "Same min overlap",
			p:        &PortRange{16, 25, ProtoTcp},
			o:        &PortRange{16, 21, ProtoTcp},
			expected: false,
		},
		{
			id:       "Same max overlap",
			p:        &PortRange{16, 25, ProtoTcp},
			o:        &PortRange{18, 25, ProtoTcp},
			expected: false,
		},
		{
			id:       "Nesting overlap",
			p:        &PortRange{16, 27, ProtoTcp},
			o:        &PortRange{18, 25, ProtoTcp},
			expected: false,
		},
		{
			id:       "Distinct overlap",
			p:        &PortRange{16, 25, ProtoTcp},
			o:        &PortRange{18, 27, ProtoTcp},
			expected: false,
		},
		{
			id:       "Disjoint ranges",
			p:        &PortRange{16, 17, ProtoTcp},
			o:        &PortRange{20, 27, ProtoTcp},
			expected: false,
		},
		{
			id:       "Single port adjacent",
			p:        &PortRange{17, 17, ProtoTcp},
			o:        &PortRange{18, 18, ProtoTcp},
			expected: true,
		},
	} {
		t.Logf("Running Adjacent test case %q", tc.id)
		outcome := tc.p.Adjacent(tc.o)
		if outcome != tc.expected {
			t.Errorf("Case p(o) %q want %v, got %v", tc.id, tc.expected, outcome)
		}
		// Overlap is a symmetric property.
		outcome = tc.o.Adjacent(tc.p)
		if outcome != tc.expected {
			t.Errorf("Case o(p) %q want %v, got %v", tc.id, tc.expected, outcome)
		}
	}
}

func TestOverlap(t *testing.T) {
	for _, tc := range []struct {
		id      string
		s, d, e *PortRange
		err     error
	}{
		{
			id:  "Proto mismatch",
			s:   &PortRange{17, 20, ProtoTcp},
			d:   &PortRange{18, 21, ProtoUdp},
			e:   &PortRange{18, 21, ProtoUdp},
			err: ErrDisjointRanges,
		},
		{
			id:  "TCP Overlap",
			s:   &PortRange{17, 20, ProtoTcp},
			d:   &PortRange{18, 21, ProtoTcp},
			e:   &PortRange{18, 20, ProtoTcp},
			err: nil,
		},
		{
			id:  "UDP Overlap",
			s:   &PortRange{17, 20, ProtoUdp},
			d:   &PortRange{18, 21, ProtoUdp},
			e:   &PortRange{18, 20, ProtoUdp},
			err: nil,
		},
		{
			id:  "No Overlap",
			s:   &PortRange{17, 18, ProtoUdp},
			d:   &PortRange{19, 21, ProtoUdp},
			e:   &PortRange{19, 21, ProtoUdp},
			err: ErrDisjointRanges,
		},
	} {
		t.Logf("Running Overlap test case %q", tc.id)
		err := tc.s.Overlap(tc.d)
		if err != tc.err {
			t.Errorf("Case %q err want %v, got %v", tc.id, err, tc.err)
		}
		if *tc.d != *tc.e {
			t.Errorf("Case %q want %v, got %v", tc.id, tc.e, tc.d)
		}
	}
}

func TestMergeWith(t *testing.T) {
	for _, tc := range []struct {
		id      string
		s, d, e *PortRange
		err     error
	}{
		{
			id:  "Proto mismatch",
			s:   &PortRange{17, 20, ProtoTcp},
			d:   &PortRange{18, 21, ProtoUdp},
			e:   &PortRange{18, 21, ProtoUdp},
			err: ErrDisjointRanges,
		},
		{
			id:  "Disparate Ranges",
			s:   &PortRange{17, 18, ProtoUdp},
			d:   &PortRange{20, 21, ProtoUdp},
			e:   &PortRange{20, 21, ProtoUdp},
			err: ErrDisjointRanges,
		},
		{
			id:  "TCP Overlap 1",
			s:   &PortRange{17, 20, ProtoTcp},
			d:   &PortRange{18, 21, ProtoTcp},
			e:   &PortRange{17, 21, ProtoTcp},
			err: nil,
		},
		{
			id:  "TCP Overlap 2",
			s:   &PortRange{18, 21, ProtoTcp},
			d:   &PortRange{17, 20, ProtoTcp},
			e:   &PortRange{17, 21, ProtoTcp},
			err: nil,
		},
		{
			id:  "UDP Adjacent",
			s:   &PortRange{17, 18, ProtoUdp},
			d:   &PortRange{19, 21, ProtoUdp},
			e:   &PortRange{17, 21, ProtoUdp},
			err: nil,
		},
		{
			id:  "UDP Adjacent 2",
			s:   &PortRange{19, 21, ProtoUdp},
			d:   &PortRange{17, 18, ProtoUdp},
			e:   &PortRange{17, 21, ProtoUdp},
			err: nil,
		},
	} {
		t.Logf("Running MergeWith test case %q", tc.id)
		err := tc.s.MergeWith(tc.d)
		if err != tc.err {
			t.Errorf("Case %q err want %v, got %v", tc.id, err, tc.err)
		}
		if *tc.d != *tc.e {
			t.Errorf("Case %q want %v, got %v", tc.id, tc.e, tc.d)
		}
	}
}

func TestEntirelyLessThan(t *testing.T) {
	for _, tc := range []struct {
		id       string
		s, d     *PortRange
		expected bool
	}{
		{
			id:       "Proto Mismatch 1",
			s:        &PortRange{17, 20, ProtoTcp},
			d:        &PortRange{18, 21, ProtoUdp},
			expected: true,
		},
		{
			id:       "Proto Mismatch 2",
			s:        &PortRange{50, 55, ProtoTcp},
			d:        &PortRange{18, 21, ProtoUdp},
			expected: true,
		},
		{
			id:       "Range Overlap",
			s:        &PortRange{17, 20, ProtoTcp},
			d:        &PortRange{18, 21, ProtoTcp},
			expected: false,
		},
		{
			id:       "Range Superset",
			s:        &PortRange{17, 25, ProtoTcp},
			d:        &PortRange{18, 21, ProtoTcp},
			expected: false,
		},
		{
			id:       "Range Equal",
			s:        &PortRange{17, 21, ProtoTcp},
			d:        &PortRange{17, 21, ProtoTcp},
			expected: false,
		},
		{
			id:       "Range EntirelyLessThan",
			s:        &PortRange{17, 21, ProtoTcp},
			d:        &PortRange{24, 26, ProtoTcp},
			expected: true,
		},
	} {
		t.Logf("Running EntirelyLessThan test case %q", tc.id)
		got := tc.s.EntirelyLessThan(tc.d)
		if got != tc.expected {
			t.Errorf("Case %q want %v, got %v", tc.id, tc.expected, got)
		}
	}
}
