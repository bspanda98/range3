// Copyright 2018 The Range Core Authors
// Copyright 2016 The go-ethereum Authors
// This file is part of the Range Core library.
//
// The Range Core library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Range Core library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Range Core library. If not, see <http://www.gnu.org/licenses/>.

package metrics

import "sync/atomic"

// Gauges hold an int64 value that can be set arbitrarily.
type Gauge interface {
	Snapshot() Gauge
	Update(int64)
	Value() int64
}

// GetOrRegisterGauge returns an existing Gauge or constructs and registers a
// new StandardGauge.
func GetOrRegisterGauge(name string, r Registry) Gauge {
	if nil == r {
		r = DefaultRegistry
	}
	return r.GetOrRegister(name, NewGauge).(Gauge)
}

// NewGauge constructs a new StandardGauge.
func NewGauge() Gauge {
	if !Enabled {
		return NilGauge{}
	}
	return &StandardGauge{0}
}

// NewRegisteredGauge constructs and registers a new StandardGauge.
func NewRegisteredGauge(name string, r Registry) Gauge {
	c := NewGauge()
	if nil == r {
		r = DefaultRegistry
	}
	r.Register(name, c)
	return c
}

// NewFunctionalGauge constructs a new FunctionalGauge.
func NewFunctionalGauge(f func() int64) Gauge {
	if !Enabled {
		return NilGauge{}
	}
	return &FunctionalGauge{value: f}
}

// NewRegisteredFunctionalGauge constructs and registers a new StandardGauge.
func NewRegisteredFunctionalGauge(name string, r Registry, f func() int64) Gauge {
	c := NewFunctionalGauge(f)
	if nil == r {
		r = DefaultRegistry
	}
	r.Register(name, c)
	return c
}

// GaugeSnapshot is a read-only copy of another Gauge.
type GaugeSnapshot int64

// Snapshot returns the snapshot.
func (g GaugeSnapshot) Snapshot() Gauge { return g }

// Update panics.
func (GaugeSnapshot) Update(int64) {
	panic("Update called on a GaugeSnapshot")
}

// Value returns the value at the time the snapshot was taken.
func (g GaugeSnapshot) Value() int64 { return int64(g) }

// NilGauge is a no-op Gauge.
type NilGauge struct{}

// Snapshot is a no-op.
func (NilGauge) Snapshot() Gauge { return NilGauge{} }

// Update is a no-op.
func (NilGauge) Update(v int64) {}

// Value is a no-op.
func (NilGauge) Value() int64 { return 0 }

// StandardGauge is the standard implementation of a Gauge and uses the
// sync/atomic package to manage a single int64 value.
type StandardGauge struct {
	value int64
}

// Snapshot returns a read-only copy of the gauge.
func (g *StandardGauge) Snapshot() Gauge {
	return GaugeSnapshot(g.Value())
}

// Update updates the gauge's value.
func (g *StandardGauge) Update(v int64) {
	atomic.StoreInt64(&g.value, v)
}

// Value returns the gauge's current value.
func (g *StandardGauge) Value() int64 {
	return atomic.LoadInt64(&g.value)
}

// FunctionalGauge returns value from given function
type FunctionalGauge struct {
	value func() int64
}

// Value returns the gauge's current value.
func (g FunctionalGauge) Value() int64 {
	return g.value()
}

// Snapshot returns the snapshot.
func (g FunctionalGauge) Snapshot() Gauge { return GaugeSnapshot(g.Value()) }

// Update panics.
func (FunctionalGauge) Update(int64) {
	panic("Update called on a FunctionalGauge")
}
