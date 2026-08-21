package peerdisco

import "example.com/peerdisco/internal/meta"

// Snapshot 导出完整成员+元数据视图。
func (c *Cluster) Snapshot() []MemberView {
	c.mu.Lock()
	defer c.mu.Unlock()
	raw := c.table.List()
	out := make([]MemberView, 0, len(raw))
	for _, m := range raw {
		tags := append([]string(nil), m.Tags...)
		out = append(out, MemberView{
			ID:          m.ID,
			Addr:        m.Addr,
			State:       State(m.State),
			Incarnation: m.Incarnation,
			Meta:        meta.CloneMap(c.meta.Get(m.ID)),
			Tags:        tags,
			JoinedAt:    m.JoinedAt,
			Updated:     m.Updated,
		})
	}
	return out
}

// SuspicionList 怀疑队列快照。
func (c *Cluster) SuspicionList() []SuspicionEvent {
	c.mu.Lock()
	defer c.mu.Unlock()
	items := c.susp.List()
	out := make([]SuspicionEvent, 0, len(items))
	for _, it := range items {
		out = append(out, SuspicionEvent{
			ID:       it.ID,
			Since:    it.Since,
			Deadline: it.Deadline,
			From:     it.From,
		})
	}
	return out
}
