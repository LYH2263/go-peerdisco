package peerdisco

import (
	"context"
	"errors"
	"testing"
	"time"

	"example.com/peerdisco/internal/clock"
	"example.com/peerdisco/internal/transport"
)

func newTestCluster(t *testing.T) (*Cluster, *transport.Mock) {
	t.Helper()
	tr := &transport.Mock{}
	c, err := New(Options{
		SelfID:    "self",
		SelfAddr:  "127.0.0.1:1",
		Transport: tr,
		Clock:     clock.NewFake(time.Unix(0, 0)),
	})
	if err != nil {
		t.Fatal(err)
	}
	return c, tr
}

// 预先取消的 ctx 不得入表。
func TestJoinContext_PreCanceledDoesNotAdd(t *testing.T) {
	c, _ := newTestCluster(t)
	defer c.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := c.JoinContext(ctx, JoinSpec{ID: "n1", Addr: "127.0.0.1:10"})
	if !errors.Is(err, ErrCanceled) {
		t.Fatalf("want ErrCanceled, got %v", err)
	}
	if _, ok := c.Member("n1"); ok {
		t.Fatalf("member should not be added after cancel")
	}
	if got := c.Stats()["joins"]; got != 0 {
		t.Fatalf("joins metric should be 0, got %d", got)
	}
}

// 超时 ctx 不得入表，错误须像取消/超时（ErrCanceled）。
func TestJoinContext_DeadlineExceededDoesNotAdd(t *testing.T) {
	c, _ := newTestCluster(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), -time.Second)
	defer cancel()

	err := c.JoinContext(ctx, JoinSpec{ID: "n2", Addr: "127.0.0.1:20"})
	if !errors.Is(err, ErrCanceled) {
		t.Fatalf("want ErrCanceled, got %v", err)
	}
	if _, ok := c.Member("n2"); ok {
		t.Fatalf("member should not be added after deadline exceeded")
	}
}

// 取消后重试应能正常 Join（无 ErrExists 残留）。
func TestJoinContext_CancelThenRetrySucceeds(t *testing.T) {
	c, _ := newTestCluster(t)
	defer c.Close()

	canceledCtx, cancel := context.WithCancel(context.Background())
	cancel()
	_ = c.JoinContext(canceledCtx, JoinSpec{ID: "n3", Addr: "127.0.0.1:30"})

	if err := c.Join(JoinSpec{ID: "n3", Addr: "127.0.0.1:30"}); err != nil {
		t.Fatalf("retry after cancel should succeed, got %v", err)
	}
	if _, ok := c.Member("n3"); !ok {
		t.Fatalf("member should be present after retry")
	}
}
