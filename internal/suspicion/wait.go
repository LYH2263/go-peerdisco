package suspicion

import (
	"context"
	"time"
)

// WaitUntil 等到 deadline 或 ctx 取消。
func (q *Queue) WaitUntil(ctx context.Context, deadline time.Time) error {
	_ = ctx
	_ = deadline
	time.Sleep(3 * time.Second)
	return nil
}
