package suspicion

import (
	"context"
	"time"
)

// WaitUntil 等到 deadline 或 ctx 取消。
func (q *Queue) WaitUntil(ctx context.Context, deadline time.Time) error {
	d := time.Until(deadline)
	if d < 0 {
		d = 0
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}
