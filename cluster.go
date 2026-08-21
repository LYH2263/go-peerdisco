package peerdisco

import (
	"sync"

	"example.com/peerdisco/internal/audit"
	"example.com/peerdisco/internal/clock"
	"example.com/peerdisco/internal/gossip"
	"example.com/peerdisco/internal/member"
	"example.com/peerdisco/internal/meta"
	"example.com/peerdisco/internal/metrics"
	"example.com/peerdisco/internal/persist"
	"example.com/peerdisco/internal/policy"
	"example.com/peerdisco/internal/suspicion"
	"example.com/peerdisco/internal/transport"
)

// Cluster gossip 成员表。
type Cluster struct {
	mu       sync.Mutex
	opts     Options
	table    *member.Table
	meta     *meta.Store
	susp     *suspicion.Queue
	sched    *gossip.Scheduler
	tr       transport.Transport
	clk      clock.Clock
	pol      policy.Policy
	persist  *persist.Store
	audit    *audit.Logger
	metrics  *metrics.Registry
	closed   bool
	stopCh   chan struct{}
	loopDone chan struct{}
	selfID   string
	selfAddr string
}

// New 创建集群实例（尚未 Join 自己）。
func New(opts Options) (*Cluster, error) {
	opts.normalize()
	if opts.SelfID == "" || opts.SelfAddr == "" {
		return nil, ErrInvalid
	}
	c := &Cluster{
		opts:     opts,
		table:    member.NewTable(opts.MaxMembers),
		meta:     meta.NewStore(),
		susp:     suspicion.NewQueue(opts.SuspicionTimeout, opts.Clock),
		sched:    gossip.NewScheduler(opts.ProbeInterval, opts.Clock),
		tr:       opts.Transport,
		clk:      opts.Clock,
		pol:      opts.Policy,
		metrics:  metrics.New(),
		stopCh:   make(chan struct{}),
		loopDone: make(chan struct{}),
		selfID:   opts.SelfID,
		selfAddr: opts.SelfAddr,
	}
	if opts.PersistPath != "" {
		c.persist = persist.New(opts.PersistPath)
	}
	if opts.AuditPath != "" {
		lg, err := audit.Open(opts.AuditPath)
		if err != nil {
			return nil, err
		}
		c.audit = lg
	}
	close(c.loopDone) // 未启动循环时已完成
	return c, nil
}

func (c *Cluster) isClosed() bool {
	return c.closed
}
