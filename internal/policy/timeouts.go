package policy

import "time"

const (
	DefaultProbeInterval    = 200 * time.Millisecond
	DefaultSuspicionTimeout = 500 * time.Millisecond
	DefaultIndirectFanout   = 3
	DefaultMaxMembers       = 256
)

type Policy struct {
	ProbeInterval    time.Duration
	SuspicionTimeout time.Duration
	IndirectFanout   int
	MaxMembers       int
}

func Default() Policy {
	return Policy{
		ProbeInterval:    DefaultProbeInterval,
		SuspicionTimeout: DefaultSuspicionTimeout,
		IndirectFanout:   DefaultIndirectFanout,
		MaxMembers:       DefaultMaxMembers,
	}
}

func (p Policy) ClampFanout(n int) int {
	if n <= 0 {
		return p.IndirectFanout
	}
	if n > p.MaxMembers {
		return p.MaxMembers
	}
	return n
}
