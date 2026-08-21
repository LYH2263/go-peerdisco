package peerdisco

import "example.com/peerdisco/internal/meta"

// SetMeta 更新节点元数据并审计。
func (c *Cluster) SetMeta(id string, kv map[string]string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return ErrClosed
	}
	if _, ok := c.table.Get(id); !ok {
		return ErrNotMember
	}
	cloned := meta.CloneMap(kv)
	c.meta.Set(id, cloned)
	if c.audit != nil {
		_ = c.audit.Write("meta", id, meta.Format(cloned))
		// 轮转依赖 Close 真正释放句柄
	}
	return c.persistLocked()
}

// GetMeta 返回元数据拷贝。
func (c *Cluster) GetMeta(id string) (map[string]string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.table.Get(id); !ok {
		return nil, false
	}
	return c.meta.GetClone(id), true
}
