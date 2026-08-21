package metrics

import "sync"

type Registry struct {
	mu       sync.Mutex
	joins    int64
	leaves   int64
	pings    int64
	indirect int64
	suspects int64
	dead     int64
	ticks    int64
}

func New() *Registry { return &Registry{} }

func (r *Registry) IncJoins()    { r.add(&r.joins) }
func (r *Registry) IncLeaves()   { r.add(&r.leaves) }
func (r *Registry) IncPings()    { r.add(&r.pings) }
func (r *Registry) IncIndirect() { r.add(&r.indirect) }
func (r *Registry) IncSuspects() { r.add(&r.suspects) }
func (r *Registry) IncDead()     { r.add(&r.dead) }
func (r *Registry) IncTicks()    { r.add(&r.ticks) }

func (r *Registry) add(p *int64) {
	r.mu.Lock()
	*p++
	r.mu.Unlock()
}

func (r *Registry) Snapshot() map[string]int64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	return map[string]int64{
		"joins":    r.joins,
		"leaves":   r.leaves,
		"pings":    r.pings,
		"indirect": r.indirect,
		"suspects": r.suspects,
		"dead":     r.dead,
		"ticks":    r.ticks,
	}
}
