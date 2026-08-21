package peerdisco

import (
	"context"

	"example.com/peerdisco/internal/member"
)

// MarkSuspect 将成员置为怀疑并启动超时。
func (c *Cluster) MarkSuspect(id, from string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return ErrClosed
	}
	m, ok := c.table.Get(id)
	if !ok {
		return ErrNotMember
	}
	m.State = member.Suspect
	m.Updated = c.clk.Now()
	_ = c.table.Update(m)
	c.susp.Add(id, from, c.clk.Now())
	c.metrics.IncSuspects()
	return nil
}

// PromoteDead 到期后晋升 Dead。
func (c *Cluster) PromoteDead(id string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return ErrClosed
	}
	m, ok := c.table.Get(id)
	if !ok {
		return ErrNotMember
	}
	m.State = member.Dead
	m.Updated = c.clk.Now()
	_ = c.table.Update(m)
	c.susp.Cancel(id)
	c.metrics.IncDead()
	return c.persistLocked()
}

// WaitSuspicion 等待怀疑超时（可取消）。
func (c *Cluster) WaitSuspicion(ctx context.Context, id string) error {
	c.mu.Lock()
	deadline, ok := c.susp.Deadline(id)
	c.mu.Unlock()
	if !ok {
		return ErrNotMember
	}
	return c.susp.WaitUntil(ctx, deadline)
}

// RefuteSuspect 收到 Ack 后恢复 Alive。
func (c *Cluster) RefuteSuspect(id string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	m, ok := c.table.Get(id)
	if !ok {
		return ErrNotMember
	}
	m.State = member.Alive
	m.Incarnation++
	m.Updated = c.clk.Now()
	_ = c.table.Update(m)
	c.susp.Cancel(id)
	return nil
}
