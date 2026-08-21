package peerdisco_test

import (
	"errors"
	"testing"

	"example.com/peerdisco"
	"example.com/peerdisco/internal/transport"
)

func TestBug03_JoinAfterCloseNoPanic(t *testing.T) {
	tr := &transport.Mock{}
	c, err := peerdisco.New(peerdisco.Options{SelfID: "s", SelfAddr: "127.0.0.1:1", Transport: tr})
	if err != nil {
		t.Fatal(err)
	}
	if err := c.Close(); err != nil {
		t.Fatal(err)
	}
	err = c.Join(peerdisco.JoinSpec{ID: "z", Addr: "127.0.0.1:9"})
	if !errors.Is(err, peerdisco.ErrClosed) {
		t.Fatalf("got %v", err)
	}
}
