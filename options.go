package peerdisco

import (
	"time"

	"example.com/peerdisco/internal/clock"
	"example.com/peerdisco/internal/policy"
	"example.com/peerdisco/internal/transport"
)

// Options 集群配置。
type Options struct {
	SelfID           string
	SelfAddr         string
	ProbeInterval    time.Duration
	SuspicionTimeout time.Duration
	IndirectFanout   int
	MaxMembers       int
	PersistPath      string
	Transport        transport.Transport
	Clock            clock.Clock
	Policy           policy.Policy
	AuditPath        string
}

func (o *Options) normalize() {
	if o.ProbeInterval <= 0 {
		o.ProbeInterval = policy.DefaultProbeInterval
	}
	if o.SuspicionTimeout <= 0 {
		o.SuspicionTimeout = policy.DefaultSuspicionTimeout
	}
	if o.IndirectFanout <= 0 {
		o.IndirectFanout = policy.DefaultIndirectFanout
	}
	if o.MaxMembers <= 0 {
		o.MaxMembers = policy.DefaultMaxMembers
	}
	if o.Clock == nil {
		o.Clock = clock.Real{}
	}
	if o.Policy.ProbeInterval == 0 {
		o.Policy = policy.Default()
	}
}
