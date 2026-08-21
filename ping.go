package peerdisco

import (
	"context"

	"example.com/peerdisco/internal/gossip"
)

// Ping 对成员发直接探测。
func (c *Cluster) Ping(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return wrapCancel(err)
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return ErrClosed
	}
	if c.tr == nil {
		return ErrNoTransport
	}
	m, ok := c.table.Get(id)
	if !ok {
		return ErrNotMember
	}
	payload := gossip.EncodePing(c.selfID, id)
	if err := c.tr.Send(m.Addr, payload); err != nil {
		return err
	}
	c.metrics.IncPings()
	return nil
}

// IndirectPing 选 fanout 个节点间接探测。
func (c *Cluster) IndirectPing(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return wrapCancel(err)
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return ErrClosed
	}
	if c.tr == nil {
		return ErrNoTransport
	}
	helpers := c.table.PickHelpers(id, c.opts.IndirectFanout)
	payload := gossip.EncodeIndirect(c.selfID, id)
	for _, h := range helpers {
		_ = c.tr.Send(h.Addr, payload)
	}
	c.metrics.IncIndirect()
	return nil
}
