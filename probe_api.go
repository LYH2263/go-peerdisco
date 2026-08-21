package peerdisco

import "context"

// ProbeOnce 对目标执行一次 Ping，失败则标记怀疑。
func (c *Cluster) ProbeOnce(ctx context.Context, id string) error {
	if err := c.Ping(ctx, id); err != nil {
		_ = c.MarkSuspect(id, c.selfID)
		return err
	}
	return c.RefuteSuspect(id)
}

// AliveCount 存活成员数。
func (c *Cluster) AliveCount() int {
	n := 0
	for _, m := range c.Members() {
		if m.State == StateAlive {
			n++
		}
	}
	return n
}
