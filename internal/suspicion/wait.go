package suspicion

import (
	"context"
	"time"
)

// WaitUntil 等到 deadline 或 ctx 取消。
// deadline 已到则返回 nil；ctx 取消则返回 ctx.Err()（尽快返回，不再阻塞）。
func (q *Queue) WaitUntil(ctx context.Context, deadline time.Time) error {
	// 剩余等待时长。非正表示已到期，直接返回，不进入 select，避免
	// 构造已过期的 timer 造成无谓阻塞。
	remaining := time.Until(deadline)
	if remaining <= 0 {
		return nil
	}

	t := time.NewTimer(remaining)
	defer t.Stop()

	select {
	case <-t.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
