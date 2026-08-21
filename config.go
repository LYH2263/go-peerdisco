package peerdisco

import "time"

// ConfigView 对外配置只读视图。
type ConfigView struct {
	SelfID           string
	SelfAddr         string
	ProbeInterval    time.Duration
	SuspicionTimeout time.Duration
	IndirectFanout   int
	MaxMembers       int
}

// Config 返回当前配置快照。
func (c *Cluster) Config() ConfigView {
	c.mu.Lock()
	defer c.mu.Unlock()
	return ConfigView{
		SelfID:           c.selfID,
		SelfAddr:         c.selfAddr,
		ProbeInterval:    c.opts.ProbeInterval,
		SuspicionTimeout: c.opts.SuspicionTimeout,
		IndirectFanout:   c.opts.IndirectFanout,
		MaxMembers:       c.opts.MaxMembers,
	}
}

// HasMember 快速存在性检查。
func (c *Cluster) HasMember(id string) bool {
	_, ok := c.Member(id)
	return ok
}
