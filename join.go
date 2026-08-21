package peerdisco

import (
	"context"

	"example.com/peerdisco/internal/member"
	"example.com/peerdisco/internal/meta"
)

// Join 将节点加入成员表并持久化。
func (c *Cluster) Join(spec JoinSpec) error {
	return c.JoinContext(context.Background(), spec)
}

// JoinContext 可取消的 Join。
func (c *Cluster) JoinContext(ctx context.Context, spec JoinSpec) error {
	if err := ctx.Err(); err != nil {
		return wrapCancel(err)
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.tr == nil {
		return ErrNoTransport
	}
	if spec.ID == "" || spec.Addr == "" {
		return ErrInvalid
	}
	if _, ok := c.table.Get(spec.ID); ok {
		return ErrExists
	}
	if err := c.sched.WaitJoin(ctx); err != nil {
		return wrapCancel(err)
	}
	m := member.Member{
		ID:          spec.ID,
		Addr:        spec.Addr,
		State:       member.Alive,
		Incarnation: 1,
		Tags:        append([]string(nil), spec.Tags...),
		JoinedAt:    c.clk.Now(),
		Updated:     c.clk.Now(),
	}
	metaCopy := meta.CloneMap(spec.Meta)
	if err := c.table.Add(m); err != nil {
		return err
	}
	c.meta.Set(spec.ID, metaCopy)
	if err := c.persistLocked(); err != nil {
		_ = c.table.Remove(spec.ID)
		c.meta.Delete(spec.ID)
		return err
	}
	c.metrics.IncJoins()
	if c.audit != nil {
		_ = c.audit.Write("join", spec.ID, spec.Addr)
	}
	_ = c.tr.Send(spec.Addr, []byte("join:"+spec.ID))
	return nil
}
