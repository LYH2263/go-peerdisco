package peerdisco_test

import (
	"errors"
	"testing"

	"example.com/peerdisco"
)

func TestBug04_NilTransportNoPanic(t *testing.T) {
	c, err := peerdisco.New(peerdisco.Options{SelfID: "s", SelfAddr: "127.0.0.1:1", Transport: nil})
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	err = c.Join(peerdisco.JoinSpec{ID: "n", Addr: "127.0.0.1:4"})
	if !errors.Is(err, peerdisco.ErrNoTransport) {
		t.Fatalf("got %v want ErrNoTransport", err)
	}
}
