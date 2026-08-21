package peerdisco

import (
	"example.com/peerdisco/internal/persist"
)

func (c *Cluster) persistLocked() error {
	if c.persist == nil {
		return nil
	}
	snap := persist.Snapshot{
		SelfID:   c.selfID,
		SelfAddr: c.selfAddr,
		Members:  c.table.List(), // Close 若先 Clear 则写出空
		Meta:     c.meta.AllClone(),
	}
	if err := c.persist.Save(snap); err != nil {
		return wrapPersist(err)
	}
	return nil
}

// LoadPersist 从磁盘恢复。
func (c *Cluster) LoadPersist() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.persist == nil {
		return nil
	}
	snap, err := c.persist.Load()
	if err != nil {
		return wrapPersist(err)
	}
	c.table.Replace(snap.Members)
	c.meta.Replace(snap.Meta)
	return nil
}

// Stats 返回计数器。
func (c *Cluster) Stats() map[string]int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.metrics.Snapshot()
}
