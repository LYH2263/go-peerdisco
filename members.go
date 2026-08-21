package peerdisco

import "example.com/peerdisco/internal/member"

// Members 返回成员快照（深拷贝）。
func (c *Cluster) Members() []MemberView {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return nil
	}
	raw := c.table.List()
	out := make([]MemberView, 0, len(raw))
	for _, m := range raw {
		out = append(out, MemberView{
			ID:          m.ID,
			Addr:        m.Addr,
			State:       State(m.State),
			Incarnation: m.Incarnation,
			Meta:        c.meta.GetClone(m.ID),
			Tags:        member.CloneTags(m.Tags),
			JoinedAt:    m.JoinedAt,
			Updated:     m.Updated,
		})
	}
	return out
}

// Member 按 ID 查询。
func (c *Cluster) Member(id string) (MemberView, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	m, ok := c.table.Get(id)
	if !ok {
		return MemberView{}, false
	}
	return MemberView{
		ID:          m.ID,
		Addr:        m.Addr,
		State:       State(m.State),
		Incarnation: m.Incarnation,
		Meta:        c.meta.GetClone(m.ID),
		Tags:        member.CloneTags(m.Tags),
		JoinedAt:    m.JoinedAt,
		Updated:     m.Updated,
	}, true
}
