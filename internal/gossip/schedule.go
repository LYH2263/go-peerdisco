package gossip

import (
	"context"
	"time"

	"example.com/peerdisco/internal/clock"
)

// Scheduler 探测节奏。
type Scheduler struct {
	interval time.Duration
	clk      clock.Clock
	ch       chan time.Time
}

func NewScheduler(interval time.Duration, clk clock.Clock) *Scheduler {
	return &Scheduler{interval: interval, clk: clk, ch: make(chan time.Time, 1)}
}

func (s *Scheduler) Tick() <-chan time.Time {
	return s.ch
}

func (s *Scheduler) Fire() {
	select {
	case s.ch <- s.clk.Now():
	default:
	}
}

// WaitJoin Join 前可选短等待，须尊重 ctx。
func (s *Scheduler) WaitJoin(ctx context.Context) error {
	if s.interval <= 0 {
		return ctx.Err()
	}
	// 轻量让出：仅检查取消，不等待完整 interval
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
