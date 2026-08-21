package peerdisco_test

import (
	"path/filepath"
	"testing"

	"example.com/peerdisco"
	"example.com/peerdisco/internal/transport"
)

func TestBug01_MembersMetaSliceAlias(t *testing.T) {
	tr := &transport.Mock{}
	c, err := peerdisco.New(peerdisco.Options{SelfID: "s", SelfAddr: "127.0.0.1:1", Transport: tr})
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	if err := c.Join(peerdisco.JoinSpec{ID: "a", Addr: "127.0.0.1:2", Tags: []string{"keep"}}); err != nil {
		t.Fatal(err)
	}
	ms := c.Members()
	if len(ms) != 1 {
		t.Fatalf("len=%d", len(ms))
	}
	ms[0].Tags[0] = "MUTATED"
	ms2 := c.Members()
	if ms2[0].Tags[0] != "keep" {
		t.Fatalf("tags aliased: %v", ms2[0].Tags)
	}
	_ = filepath.Separator
}
