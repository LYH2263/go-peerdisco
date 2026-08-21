package peerdisco_test

import (
	"context"
	"testing"
	"time"

	"example.com/peerdisco"
	"example.com/peerdisco/internal/transport"
)

func TestBug08_SuspicionWaitHonorsContext(t *testing.T) {
	tr := &transport.Mock{}
	c, err := peerdisco.New(peerdisco.Options{
		SelfID: "s", SelfAddr: "127.0.0.1:1", Transport: tr,
		SuspicionTimeout: time.Hour,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	if err := c.Join(peerdisco.JoinSpec{ID: "q", Addr: "127.0.0.1:7"}); err != nil {
		t.Fatal(err)
	}
	if err := c.MarkSuspect("q", "s"); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()
	start := time.Now()
	err = c.WaitSuspicion(ctx, "q")
	if err == nil {
		t.Fatal("expected cancel/timeout")
	}
	if time.Since(start) > 2*time.Second {
		t.Fatalf("WaitSuspicion ignored ctx, took %v", time.Since(start))
	}
}
