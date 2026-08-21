package peerdisco_test

import (
	"errors"
	"testing"

	"example.com/peerdisco"
	"example.com/peerdisco/internal/transport"
)

func TestBug05_LeaveErrorWrapsSentinel(t *testing.T) {
	tr := &transport.Mock{}
	c, err := peerdisco.New(peerdisco.Options{SelfID: "s", SelfAddr: "127.0.0.1:1", Transport: tr})
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	err = c.Leave("missing")
	if !errors.Is(err, peerdisco.ErrNotMember) {
		t.Fatalf("got %v", err)
	}
}
