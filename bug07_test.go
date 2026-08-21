package peerdisco_test

import (
	"context"
	"errors"
	"testing"

	"example.com/peerdisco"
	"example.com/peerdisco/internal/transport"
)

func TestBug07_JoinContextHonorsCancel(t *testing.T) {
	tr := &transport.Mock{}
	c, err := peerdisco.New(peerdisco.Options{SelfID: "s", SelfAddr: "127.0.0.1:1", Transport: tr})
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err = c.JoinContext(ctx, peerdisco.JoinSpec{ID: "c", Addr: "127.0.0.1:6"})
	if err == nil || (!errors.Is(err, peerdisco.ErrCanceled) && !errors.Is(err, context.Canceled)) {
		t.Fatalf("got %v", err)
	}
}
