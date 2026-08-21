package peerdisco_test

import (
	"path/filepath"
	"testing"

	"example.com/peerdisco"
	"example.com/peerdisco/internal/transport"
)

func TestBug10_CloseStopsLoopBeforeDrop(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "m.json")
	tr := &transport.Mock{}
	c, err := peerdisco.New(peerdisco.Options{
		SelfID: "s", SelfAddr: "127.0.0.1:1", Transport: tr, PersistPath: path,
	})
	if err != nil {
		t.Fatal(err)
	}
	c.StartLoop()
	if err := c.Join(peerdisco.JoinSpec{ID: "keep", Addr: "127.0.0.1:8", Tags: []string{"x"}}); err != nil {
		t.Fatal(err)
	}
	if err := c.Close(); err != nil {
		t.Fatal(err)
	}
	c2, err := peerdisco.New(peerdisco.Options{
		SelfID: "s", SelfAddr: "127.0.0.1:1", Transport: &transport.Mock{}, PersistPath: path,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer c2.Close()
	if err := c2.LoadPersist(); err != nil {
		t.Fatal(err)
	}
	if c2.AliveCount() < 1 {
		t.Fatal("Close flushed empty membership; Load lost members")
	}
}
