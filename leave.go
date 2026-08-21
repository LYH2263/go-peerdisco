package peerdisco

import "example.com/peerdisco/internal/member"

// Leave 将成员标记为 Left 并移出活跃表。
func (c *Cluster) Leave(id string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return ErrClosed
	}
	m, ok := c.table.Get(id)
	if !ok {
		return wrapNotMember(id)
	}
	m.State = member.Left
	m.Updated = c.clk.Now()
	_ = c.table.Update(m)
	c.susp.Cancel(id)
	c.meta.Delete(id)
	if err := c.table.Remove(id); err != nil {
		return wrapNotMember(id)
	}
	c.metrics.IncLeaves()
	if c.audit != nil {
		_ = c.audit.Write("leave", id, m.Addr)
	}
	return c.persistLocked()
}
